# Deployment

Step-by-step guide to deploy CloudMonitor on a Linux VPS.

## Requirements

- A VPS with Ubuntu 22.04 or later (tested on 26.04 LTS).
- A public IP address. No domain is required.
- A Telegram bot token and chat ID (see the README).
- An SSH key pair on the local machine.

## 1. Provision the VPS

Any provider works. This guide was validated on a Hetzner Cloud CX23
instance (2 vCPU, 4 GB RAM, 40 GB SSD, x86) in Helsinki.

Ubuntu 26.04 LTS is the recommended image.

When creating the instance, add your SSH public key. Never enable
password authentication.

## 2. Initial hardening

Connect as `root` with the SSH key:

```
ssh root@YOUR_SERVER_IP
```

Update the system and reboot to apply the new kernel:

```
apt update && apt upgrade -y
reboot
```

Reconnect after the reboot. Create a non-root user:

```
adduser carlos
usermod -aG sudo carlos
```

Copy the SSH authorized keys from `root` to the new user:

```
mkdir -p /home/carlos/.ssh
cp /root/.ssh/authorized_keys /home/carlos/.ssh/
chown -R carlos:carlos /home/carlos/.ssh
chmod 700 /home/carlos/.ssh
chmod 600 /home/carlos/.ssh/authorized_keys
```

Test the new user in a separate terminal before continuing:

```
ssh carlos@YOUR_SERVER_IP
```

## 3. Disable root login over SSH

Create a drop-in configuration file:

```
sudo bash -c "printf 'PermitRootLogin no\nPasswordAuthentication no\n' > /etc/ssh/sshd_config.d/99-cloudmonitor.conf"
sudo sshd -t
sudo systemctl restart ssh
```

Verify in a new terminal that `root` is now rejected:

```
ssh root@YOUR_SERVER_IP
# Permission denied (publickey).
```

## 4. Firewall

Open only the ports the service needs:

```
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

## 5. Build the binary locally

From the project root on your development machine:

```
GOOS=linux GOARCH=amd64 go build -o bin/cloudmonitor ./cmd/monitor
```

This produces a static binary with no runtime dependencies.

## 6. Upload to the server

```
scp bin/cloudmonitor carlos@YOUR_SERVER_IP:~/
```

If the transfer is slow, add `-C` to enable compression:

```
scp -C bin/cloudmonitor carlos@YOUR_SERVER_IP:~/
```

## 7. Prepare the working directory

On the server, as `carlos`:

```
mkdir -p ~/cloudmonitor-app
mv ~/cloudmonitor ~/cloudmonitor-app/cloudmonitor
chmod +x ~/cloudmonitor-app/cloudmonitor
cd ~/cloudmonitor-app
```

Create `config.yaml` in this directory. It must contain the Telegram
token, chat ID, and the list of targets. See `config.example.yaml`
for the format.

## 8. Test manually

Before installing as a service, run the binary once to confirm it works:

```
./cloudmonitor
```

Expected output:

```
level=INFO msg="telegram enabled" chat_id=...
level=INFO msg="monitor starting" targets=3 interval=30s
level=INFO msg=check target=Example status=UP status_code=200 latency_ms=...
```

Stop with `Ctrl+C`.

## 9. Create the systemd unit

```
sudo nano /etc/systemd/system/cloudmonitor.service
```

Paste:

```ini
[Unit]
Description=CloudMonitor - HTTP uptime monitor
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=carlos
WorkingDirectory=/home/carlos/cloudmonitor-app
ExecStart=/home/carlos/cloudmonitor-app/cloudmonitor
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Save and reload systemd:

```
sudo systemctl daemon-reload
sudo systemctl enable cloudmonitor
sudo systemctl start cloudmonitor
```

## 10. Verify

```
sudo systemctl status cloudmonitor
```

Expected:

```
Active: active (running) since ...
Main PID: ... (cloudmonitor)
```

Watch live logs:

```
sudo journalctl -u cloudmonitor -f
```

Press `Ctrl+C` to exit.

## Operations

| Task | Command |
| --- | --- |
| Status | `sudo systemctl status cloudmonitor` |
| Live logs | `sudo journalctl -u cloudmonitor -f` |
| Last 50 log lines | `sudo journalctl -u cloudmonitor -n 50` |
| Restart | `sudo systemctl restart cloudmonitor` |
| Stop | `sudo systemctl stop cloudmonitor` |
| Disable at boot | `sudo systemctl disable cloudmonitor` |

## Updating the service

From the development machine:

```
GOOS=linux GOARCH=amd64 go build -o bin/cloudmonitor ./cmd/monitor
scp -C bin/cloudmonitor carlos@YOUR_SERVER_IP:~/cloudmonitor-app/cloudmonitor
```

On the server:

```
sudo systemctl restart cloudmonitor
```

## Troubleshooting

**Service fails to start:** check `sudo journalctl -u cloudmonitor -n 50`.
The most common cause is a missing or malformed `config.yaml`.

**Telegram messages not arriving:** confirm the token and chat ID in
`config.yaml`. Send `/start` to the bot in Telegram if you never did.

**High memory usage:** the service uses under 20 MB under normal load.
If it grows beyond that, check the number of targets and the check
interval.