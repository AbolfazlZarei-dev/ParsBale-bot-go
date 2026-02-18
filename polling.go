package ParsBale

import (
	"context"
	"log"
	"time"
)

type Poller struct {
	bot    *Bot
	limit  int
	offset int
}

func NewPoller(bot *Bot) *Poller {
	return &Poller{
		bot:   bot,
		limit: 100,
	}
}

func (p *Poller) Start(ctx context.Context, updates chan<- Update) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			updatesResp, err := p.bot.GetUpdates(p.offset, p.limit, 30)
			if err != nil {
				log.Printf("Polling error: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for _, u := range updatesResp {
				p.offset = u.UpdateID + 1
				updates <- u
			}
		}
	}
}
