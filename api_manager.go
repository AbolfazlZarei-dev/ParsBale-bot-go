package ParsBale

import (
	"context"
	"log"
	"regexp"
	"strings"
)

type HandlerFunc func(bot *Bot, update Update)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Dispatcher struct {
	Bot         *Bot
	State       StateStorage
	handlers    []handlerEntry
	middlewares []MiddlewareFunc
}

type handlerEntry struct {
	filter func(Update) bool
	handle HandlerFunc
}

func NewDispatcher(bot *Bot) *Dispatcher {
	state := NewMemoryState()
	bot.State = state

	return &Dispatcher{
		Bot:         bot,
		State:       state,
		handlers:    make([]handlerEntry, 0),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (d *Dispatcher) Use(m MiddlewareFunc) {
	d.middlewares = append(d.middlewares, m)
}

func (d *Dispatcher) Handle(filter func(Update) bool, handler HandlerFunc) {
	h := handler
	// Reverse order middleware execution
	for i := len(d.middlewares) - 1; i >= 0; i-- {
		h = d.middlewares[i](h)
	}
	d.handlers = append(d.handlers, handlerEntry{filter: filter, handle: h})
}

// --- Routing Helpers ---

func (d *Dispatcher) OnCommand(cmd string, handler HandlerFunc) {
	d.Handle(func(u Update) bool {
		if u.Message == nil || u.Message.Text == "" {
			return false
		}
		text := u.Message.Text
		if !strings.HasPrefix(text, "/") {
			return false
		}
		parts := strings.Fields(text)
		command := strings.ToLower(parts[0])
		if strings.Contains(command, "@") {
			command = strings.Split(command, "@")[0]
		}
		return command == "/"+strings.ToLower(cmd)
	}, handler)
}

func (d *Dispatcher) OnText(regex string, handler HandlerFunc) {
	re := regexp.MustCompile(regex)
	d.Handle(func(u Update) bool {
		if u.Message == nil || u.Message.Text == "" {
			return false
		}
		return re.MatchString(u.Message.Text)
	}, handler)
}

func (d *Dispatcher) OnCallback(prefix string, handler HandlerFunc) {
	d.Handle(func(u Update) bool {
		return u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, prefix)
	}, handler)
}

func (d *Dispatcher) OnState(state string, handler HandlerFunc) {
	d.Handle(func(u Update) bool {
		var userID int64
		if u.Message != nil && u.Message.From != nil {
			userID = u.Message.From.ID
		} else if u.CallbackQuery != nil && u.CallbackQuery.From != nil {
			userID = u.CallbackQuery.From.ID
		} else {
			return false
		}

		currentState, _ := d.State.Get(userID)
		return currentState == state
	}, handler)
}

// --- Execution ---

func (d *Dispatcher) StartPolling(ctx context.Context) error {
	poller := NewPoller(d.Bot)
	updates := make(chan Update, 100)

	go func() {
		if err := poller.Start(ctx, updates); err != nil {
			log.Fatalf("Polling failed: %v", err)
		}
	}()

	log.Println("ParsBale Bot v1.0.0 Started polling...")
	for {
		select {
		case <-ctx.Done():
			return nil
		case u := <-updates:
			go d.processUpdate(u)
		}
	}
}

func (d *Dispatcher) StartWebhook(ctx context.Context, addr, path string) error {
	server := NewWebhookServer(d.Bot, addr)

	go func() {
		for {
			select {
			case u := <-server.Updates():
				go d.processUpdate(u)
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("ParsBale Bot v1.0.0 Started webhook on %s", addr)
	return server.Start(path)
}

func (d *Dispatcher) processUpdate(u Update) {
	for _, h := range d.handlers {
		if h.filter(u) {
			h.handle(d.Bot, u)
			return
		}
	}
}
