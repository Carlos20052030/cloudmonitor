package notifier

import "github.com/Carlos20052030/cloudmonitor/internal/domain"

// Notifier is the contract consumed by main. Any implementation
// (Telegram, Slack, test fake) satisfies it implicitly.
type Notifier interface {
	Notify(r domain.Result)
}