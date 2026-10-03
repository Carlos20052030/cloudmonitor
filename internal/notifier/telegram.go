package notifier

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

// antiSpam keeps the last known status per URL and decides
// whether a new notification should be sent.
type antiSpam struct {
	last map[string]domain.Status
}

func newAntiSpam() *antiSpam {
	return &antiSpam{last: make(map[string]domain.Status)}
}

// shouldNotify returns true if the status differs from the last one
// recorded for that URL. On the first sight of a URL, it always returns true.
func (a *antiSpam) shouldNotify(url string, status domain.Status) bool {
	previous, seen := a.last[url]
	if seen && previous == status {
		return false
	}
	a.last[url] = status
	return true
}

// Telegram sends messages to a chat when a target changes state.
type Telegram struct {
	bot      *tgbotapi.BotAPI
	chatID   int64
	antiSpam *antiSpam
}

func NewTelegram(token string, chatID int64) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot:      bot,
		chatID:   chatID,
		antiSpam: newAntiSpam(),
	}, nil
}

func (t *Telegram) Notify(r domain.Result) {
	if !t.antiSpam.shouldNotify(r.URL, r.Status) {
		return
	}

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
