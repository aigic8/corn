package notif

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/discord"
	"github.com/nikoksr/notify/service/telegram"
)

type Notifier = notify.Notifier

type (
	Notif struct {
		Timeout  time.Duration
		services map[string]*notify.Notify
	}

	DiscordNotifierOpts struct {
		OAuth2Token string
		BotToken    string
		ChanelIDs   []string
	}
)

func NewNotif(timeout time.Duration, services map[string][]Notifier, disabled bool) *Notif {
	resServices := make(map[string]*notify.Notify, len(services))
	for serviceName, notifiers := range services {
		option := notify.Enable
		if disabled {
			option = notify.Disable
		}
		n := notify.NewWithOptions(option)
		n.UseServices(notifiers...)
		resServices[serviceName] = n
	}
	return &Notif{services: resServices, Timeout: timeout}
}

func (n *Notif) UseService(service string) func(subject, text string) error {
	if service != "" {
		return func(subject, text string) error {
			return n.Send(service, subject, text)
		}
	} else {
		return func(subject, text string) error { return nil }
	}
}

//	generalized function to send a message
//
// it probably should not be used directly
// but only internally in `UseService`
func (n *Notif) Send(service, subject, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), n.Timeout)
	defer cancel()
	return n.services[service].Send(ctx, subject, text)
}

func NewTelegramNotifier(token string, receivers ...int64) (Notifier, error) {
	notifier, err := telegram.New(token)
	if err != nil {
		return notifier, fmt.Errorf("creating telegram notifier: %w", err)
	}

	notifier.AddReceivers(receivers...)

	return notifier, nil
}

func NewDiscordNotifier(o *DiscordNotifierOpts) (Notifier, error) {
	notifier := discord.New()
	if o.BotToken != "" {
		if err := notifier.AuthenticateWithBotToken(o.BotToken); err != nil {
			return notifier, fmt.Errorf("authenticating discord with bot token: %w", err)
		}
	} else if o.OAuth2Token != "" {
		if err := notifier.AuthenticateWithOAuth2Token(o.BotToken); err != nil {
			return notifier, fmt.Errorf("authenticating discord with oAuth2 token: %w", err)
		}
	} else {
		return notifier, errors.New("no authentication method (bot token or oAuth2 token) is provided for discord")
	}

	return notifier, nil
}
