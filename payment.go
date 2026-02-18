package ParsBale

import (
	"strconv"
)

func (b *Bot) SendInvoice(chatID int64, title, desc, payload, providerToken string, prices []LabeledPrice, photoURL string) (*Message, error) {
	body := struct {
		ChatID        string         `json:"chat_id"`
		Title         string         `json:"title"`
		Description   string         `json:"description"`
		Payload       string         `json:"payload"`
		ProviderToken string         `json:"provider_token"`
		Prices        []LabeledPrice `json:"prices"`
		PhotoURL      string         `json:"photo_url,omitempty"`
	}{
		ChatID:        strconv.FormatInt(chatID, 10),
		Title:         title,
		Description:   desc,
		Payload:       payload,
		ProviderToken: providerToken,
		Prices:        prices,
		PhotoURL:      photoURL,
	}

	data, err := b.Client.Request("sendInvoice", nil, body)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) CreateInvoiceLink(title, desc, payload, providerToken string, prices []LabeledPrice) (string, error) {
	body := struct {
		Title         string         `json:"title"`
		Description   string         `json:"description"`
		Payload       string         `json:"payload"`
		ProviderToken string         `json:"provider_token"`
		Prices        []LabeledPrice `json:"prices"`
	}{
		Title:         title,
		Description:   desc,
		Payload:       payload,
		ProviderToken: providerToken,
		Prices:        prices,
	}

	data, err := b.Client.Request("createInvoiceLink", nil, body)
	if err != nil {
		return "", err
	}
	var link string
	if err := parseResult(data, &link); err != nil {
		return "", err
	}
	return link, nil
}

func (b *Bot) AnswerPreCheckoutQuery(queryID string, ok bool, errorMsg string) (bool, error) {
	params := map[string]string{
		"pre_checkout_query_id": queryID,
		"ok":                    strconv.FormatBool(ok),
	}
	if !ok && errorMsg != "" {
		params["error_message"] = errorMsg
	}

	data, err := b.Client.Request("answerPreCheckoutQuery", params, nil)
	if err != nil {
		return false, err
	}
	var res bool
	if err := parseResult(data, &res); err != nil {
		return false, err
	}
	return res, nil
}

// InquireTransaction - اصلاح شده برای برگرداندن ساختار Transaction
func (b *Bot) InquireTransaction(transactionID string) (*Transaction, error) {
	params := map[string]string{"transaction_id": transactionID}
	data, err := b.Client.Request("inquireTransaction", params, nil)
	if err != nil {
		return nil, err
	}
	var res Transaction
	if err := parseResult(data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
