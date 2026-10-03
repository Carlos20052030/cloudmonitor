package notifier

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

// Telegram sends messages to a chat when a target changes state.
// It keeps the last known state per URL in memory to avoid spam.
type Telegram struct {
	bot       *tgbotapi.BotAPI
	chatID    int64
	lastState map[string]domain.Status
}

func NewTelegram(token string, chatID int64) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot:       bot,
		chatID:    chatID,
		lastState: make(map[string]domain.Status),
	}, nil
}

// Notify sends a message only when the status differs from the last
// known one for that URL. On the first check of each URL, it always sends.
func (t *Telegram) Notify(r domain.Result) {
	last, seen := t.lastState[r.URL]
	if seen && last == r.Status {
		return
	}
	t.lastState[r.URL] = r.Status

	var msg string
	if r.Status == domain.StatusDOWN {
		msg = fmt.Sprintf("🔴 DOWN: %s\nURL: %s\nError: %s", r.TargetName, r.URL, r.Error)
	} else {
		msg = fmt.Sprintf("🟢 UP: %s\nURL: %s\nLatency: %dms", r.TargetName, r.URL, r.LatencyMS)
	}

	m := tgbotapi.NewMessage(t.chatID, msg)
	if _, err := t.bot.Send(m); err != nil {
		slog.Error("telegram send", "err", err, "target", r.TargetName)
	}
}