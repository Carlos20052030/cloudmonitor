package notifier

import (
	"testing"

	"github.com/Carlos20052030/cloudmonitor/internal/domain"
)

func TestAntiSpam_FirstTimeSends(t *testing.T) {
	a := newAntiSpam()
	if !a.shouldNotify("https://x.com", domain.StatusUP) {
		t.Fatal("first notification should be sent")
	}
}

func TestAntiSpam_SameStatusDoesNotSend(t *testing.T) {
	a := newAntiSpam()
	a.shouldNotify("https://x.com", domain.StatusUP)
	if a.shouldNotify("https://x.com", domain.StatusUP) {
		t.Fatal("repeated UP should not send")
	}
}

func TestAntiSpam_StatusChangeSends(t *testing.T) {
	a := newAntiSpam()
	a.shouldNotify("https://x.com", domain.StatusUP)
	if !a.shouldNotify("https://x.com", domain.StatusDOWN) {
		t.Fatal("UP -> DOWN should send")
	}
	if !a.shouldNotify("https://x.com", domain.StatusUP) {
		t.Fatal("DOWN -> UP should send")
	}
}

func TestAntiSpam_IndependentPerURL(t *testing.T) {
	a := newAntiSpam()
	a.shouldNotify("https://a.com", domain.StatusUP)
	a.shouldNotify("https://b.com", domain.StatusUP)

	if a.shouldNotify("https://a.com", domain.StatusUP) {
		t.Fatal("a.com same status should not send")
	}
	if a.shouldNotify("https://b.com", domain.StatusUP) {
		t.Fatal("b.com same status should not send")
	}
}
