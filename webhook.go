package ParsBale

import (
	"encoding/json"
	"net/http"
)

type WebhookServer struct {
	Bot     *Bot
	Addr    string
	updates chan Update
}

func NewWebhookServer(bot *Bot, addr string) *WebhookServer {
	return &WebhookServer{
		Bot:     bot,
		Addr:    addr,
		updates: make(chan Update, 100),
	}
}

func (ws *WebhookServer) Updates() <-chan Update {
	return ws.updates
}

func (ws *WebhookServer) Start(path string) error {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var u Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		ws.updates <- u
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(ws.Addr, nil)
}
