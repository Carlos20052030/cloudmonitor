#!/bin/bash

# Auto-Healing Watchdog Script for Docker Services
# Checks health of critical containers and restarts them if needed

set -euo pipefail

# Configuration
LOG_FILE="/var/log/nuntix-auto-heal.log"
TELEGRAM_BOT_TOKEN="${TELEGRAM_BOT_TOKEN:-}"
TELEGRAM_CHAT_ID="${TELEGRAM_CHAT_ID:-}"

# Container configurations: name -> working directory
declare -A CONTAINERS=(
    ["bookmarks-api-api-1"]="/home/carlos/bookmarks-api"
    ["webhook-dispatcher-dispatcher-1"]="/home/carlos/webhook-dispatcher"
    ["reverse-proxy"]="/home/carlos/reverse-proxy"
)

# Ensure log file exists
mkdir -p "$(dirname "$LOG_FILE")"
touch "$LOG_FILE"

# Function to log messages
log_message() {
    local level="$1"
    local message="$2"
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

# Function to send Telegram notification
send_telegram_notification() {
    local message="$1"
    if [[ -z "$TELEGRAM_BOT_TOKEN" || -z "$TELEGRAM_CHAT_ID" ]]; then
        log_message "WARN" "Telegram credentials not configured, skipping notification"
        return 0
    fi

    local url="https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage"
    local payload
    payload=$(jq -n --arg chat_id "$TELEGRAM_CHAT_ID" --arg text "$message" '{chat_id: $chat_id, text: $text, parse_mode: "HTML"}')

    curl -s -X POST "$url" \
        -H "Content-Type: application/json" \
        -d "$payload" > /dev/null 2>&1 || log_message "ERROR" "Failed to send Telegram notification"
}

# Function to check if container is running
is_container_running() {
    local container_name="$1"
    local status
    status=$(docker inspect --format="{{.State.Running}}" "$container_name" 2>/dev/null || echo "false")
    [[ "$status" == "true" ]]
}

# Function to restart container via docker compose
restart_container() {
    local container_name="$1"
    local work_dir="$2"

    log_message "INFO" "Attempting to restart $container_name in $work_dir"

    if [[ ! -d "$work_dir" ]]; then
        log_message "ERROR" "Directory $work_dir does not exist for container $container_name"
        return 1
    fi

    cd "$work_dir"
    if docker compose up -d; then
        log_message "INFO" "Successfully restarted $container_name"
        return 0
    else
        log_message "ERROR" "Failed to restart $container_name"
        return 1
    fi
}

# Main health check loop
main() {
    log_message "INFO" "Starting auto-heal health check"

    local any_restarted=false
    local restart_messages=()

    for container_name in "${!CONTAINERS[@]}"; do
        local work_dir="${CONTAINERS[$container_name]}"

        if is_container_running "$container_name"; then
            log_message "DEBUG" "Container $container_name is healthy"
        else
            log_message "WARN" "Container $container_name is not running!"

            if restart_container "$container_name" "$work_dir"; then
                any_restarted=true
                restart_messages+=("✅ <b>$container_name</b> restarted successfully in <code>$work_dir</code>")
            else
                restart_messages+=("❌ <b>$container_name</b> FAILED to restart in <code>$work_dir</code>")
            fi
        fi
    done

    # Send Telegram notification if any container was restarted
    if [[ "$any_restarted" == "true" ]]; then
        local message="🏥 <b>Auto-Heal Alert</b>%0A%0A"
        message+="The following containers were restarted:%0A"
        for msg in "${restart_messages[@]}"; do
            message+="$msg%0A"
        done
        message+="%0A🕐 <i>$(date '+%Y-%m-%d %H:%M:%S')</i>"
        send_telegram_notification "$message"
    fi

    log_message "INFO" "Auto-heal health check completed"
}

main "$@"