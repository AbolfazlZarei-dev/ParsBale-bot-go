package ParsBale

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Bot struct {
	Client *HTTPClient
	Self   *User
	State  StateStorage
}

func NewBot(token string) (*Bot, error) {
	client := NewHTTPClient(token)
	bot := &Bot{Client: client}
	bot.State = NewMemoryState()

	user, err := bot.GetMe()
	if err != nil {
		return nil, err
	}
	bot.Self = user
	return bot, nil
}

func parseResult(data []byte, dest interface{}) error {
	var resp struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return ParseError{Raw: data, Err: err}
	}
	return json.Unmarshal(resp.Result, dest)
}

// --- Core Methods ---

func (b *Bot) GetMe() (*User, error) {
	data, err := b.Client.Request("getMe", nil, nil)
	if err != nil {
		return nil, err
	}
	var user User
	if err := parseResult(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (b *Bot) GetUpdates(offset, limit, timeout int) ([]Update, error) {
	params := map[string]string{
		"offset":  strconv.Itoa(offset),
		"limit":   strconv.Itoa(limit),
		"timeout": strconv.Itoa(timeout),
	}
	data, err := b.Client.Request("getUpdates", params, nil)
	if err != nil {
		return nil, err
	}
	var updates []Update
	if err := parseResult(data, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}

func (b *Bot) SetWebhook(webhookURL string) (bool, error) {
	params := map[string]string{"url": webhookURL}
	data, err := b.Client.Request("setWebhook", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) DeleteWebhook() (bool, error) {
	data, err := b.Client.Request("deleteWebhook", nil, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) GetWebhookInfo() (*WebhookInfo, error) {
	data, err := b.Client.Request("getWebhookInfo", nil, nil)
	if err != nil {
		return nil, err
	}
	var info WebhookInfo
	if err := parseResult(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// --- Sending Messages ---

type SendMessageOptions struct {
	ReplyToMessageID int64
	ParseMode        ParseMode
	ReplyMarkup      interface{}
}

func (b *Bot) SendMessage(chatID int64, text string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"text":    text,
	}

	if opts != nil {
		if opts.ReplyToMessageID != 0 {
			params["reply_to_message_id"] = strconv.FormatInt(opts.ReplyToMessageID, 10)
		}
		if opts.ParseMode != "" {
			params["parse_mode"] = string(opts.ParseMode)
		}
		if opts.ReplyMarkup != nil {
			params["reply_markup"] = mustMarshal(opts.ReplyMarkup)
		}
	}

	data, err := b.Client.Request("sendMessage", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) ForwardMessage(chatID, fromChatID int64, messageID int64) (*Message, error) {
	params := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"from_chat_id": strconv.FormatInt(fromChatID, 10),
		"message_id":   strconv.FormatInt(messageID, 10),
	}
	data, err := b.Client.Request("forwardMessage", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) CopyMessage(chatID, fromChatID int64, messageID int64) (*MessageId, error) {
	params := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"from_chat_id": strconv.FormatInt(fromChatID, 10),
		"message_id":   strconv.FormatInt(messageID, 10),
	}
	data, err := b.Client.Request("copyMessage", params, nil)
	if err != nil {
		return nil, err
	}
	var msgId MessageId
	if err := parseResult(data, &msgId); err != nil {
		return nil, err
	}
	return &msgId, nil
}

// --- File Upload Support (Corrected) ---

// SendPhoto accepts either string (file_id/url) or FileUpload for uploading files
func (b *Bot) SendPhoto(chatID int64, photo interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}
	if opts != nil {
		if opts.ParseMode != "" {
			params["parse_mode"] = string(opts.ParseMode)
		}
		if opts.ReplyMarkup != nil {
			params["reply_markup"] = mustMarshal(opts.ReplyMarkup)
		}
	}

	var data []byte
	var err error

	switch v := photo.(type) {
	case string:
		params["photo"] = v
		data, err = b.Client.Request("sendPhoto", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendPhoto", params, map[string]FileUpload{"photo": v})
	default:
		return nil, fmt.Errorf("invalid photo type: expected string or FileUpload")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendAudio(chatID int64, audio interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	var data []byte
	var err error

	switch v := audio.(type) {
	case string:
		params["audio"] = v
		data, err = b.Client.Request("sendAudio", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendAudio", params, map[string]FileUpload{"audio": v})
	default:
		return nil, fmt.Errorf("invalid audio type")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendDocument(chatID int64, doc interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	var data []byte
	var err error

	switch v := doc.(type) {
	case string:
		params["document"] = v
		data, err = b.Client.Request("sendDocument", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendDocument", params, map[string]FileUpload{"document": v})
	default:
		return nil, fmt.Errorf("invalid document type")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendVideo(chatID int64, video interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	var data []byte
	var err error

	switch v := video.(type) {
	case string:
		params["video"] = v
		data, err = b.Client.Request("sendVideo", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendVideo", params, map[string]FileUpload{"video": v})
	default:
		return nil, fmt.Errorf("invalid video type")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendAnimation(chatID int64, animation interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	var data []byte
	var err error

	switch v := animation.(type) {
	case string:
		params["animation"] = v
		data, err = b.Client.Request("sendAnimation", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendAnimation", params, map[string]FileUpload{"animation": v})
	default:
		return nil, fmt.Errorf("invalid animation type")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendVoice(chatID int64, voice interface{}, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if caption != "" {
		params["caption"] = caption
	}

	var data []byte
	var err error

	switch v := voice.(type) {
	case string:
		params["voice"] = v
		data, err = b.Client.Request("sendVoice", params, nil)
	case FileUpload:
		data, err = b.Client.RequestMultipart("sendVoice", params, map[string]FileUpload{"voice": v})
	default:
		return nil, fmt.Errorf("invalid voice type")
	}

	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// --- Other Methods ---

func (b *Bot) SendLocation(chatID int64, latitude, longitude float64, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":   strconv.FormatInt(chatID, 10),
		"latitude":  fmt.Sprintf("%f", latitude),
		"longitude": fmt.Sprintf("%f", longitude),
	}
	data, err := b.Client.Request("sendLocation", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendContact(chatID int64, phoneNumber, firstName string, lastName string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":      strconv.FormatInt(chatID, 10),
		"phone_number": phoneNumber,
		"first_name":   firstName,
	}
	if lastName != "" {
		params["last_name"] = lastName
	}
	data, err := b.Client.Request("sendContact", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) SendChatAction(chatID int64, action ChatAction) (bool, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"action":  string(action),
	}
	data, err := b.Client.Request("sendChatAction", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) SendMediaGroup(chatID int64, media []InputMedia, opts *SendMessageOptions) ([]Message, error) {
	mediaJson, _ := json.Marshal(media)
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"media":   string(mediaJson),
	}

	data, err := b.Client.Request("sendMediaGroup", params, nil)
	if err != nil {
		return nil, err
	}
	var msgs []Message
	if err := parseResult(data, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// --- Callbacks ---

func (b *Bot) AnswerCallbackQuery(queryID, text string, showAlert bool) (bool, error) {
	params := map[string]string{
		"callback_query_id": queryID,
		"text":              text,
		"show_alert":        strconv.FormatBool(showAlert),
	}

	data, err := b.Client.Request("answerCallbackQuery", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

// --- Editing & Admin methods remain the same ---

func (b *Bot) EditMessageText(chatID int64, messageID int64, text string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.FormatInt(messageID, 10),
		"text":       text,
	}

	if opts != nil && opts.ReplyMarkup != nil {
		params["reply_markup"] = mustMarshal(opts.ReplyMarkup)
	}

	data, err := b.Client.Request("editMessageText", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) EditMessageCaption(chatID int64, messageID int64, caption string, opts *SendMessageOptions) (*Message, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.FormatInt(messageID, 10),
		"caption":    caption,
	}

	if opts != nil && opts.ReplyMarkup != nil {
		params["reply_markup"] = mustMarshal(opts.ReplyMarkup)
	}

	data, err := b.Client.Request("editMessageCaption", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) EditMessageReplyMarkup(chatID int64, messageID int64, markup *InlineKeyboardMarkup) (*Message, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.FormatInt(messageID, 10),
	}
	if markup != nil {
		params["reply_markup"] = mustMarshal(markup)
	}

	data, err := b.Client.Request("editMessageReplyMarkup", params, nil)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := parseResult(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (b *Bot) DeleteMessage(chatID int64, messageID int64) (bool, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.FormatInt(messageID, 10),
	}
	data, err := b.Client.Request("deleteMessage", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

// --- Chat Admin (Methods unchanged, just part of the file) ---

func (b *Bot) GetChat(chatID int64) (*ChatFullInfo, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("getChat", params, nil)
	if err != nil {
		return nil, err
	}
	var chat ChatFullInfo
	if err := parseResult(data, &chat); err != nil {
		return nil, err
	}
	return &chat, nil
}

func (b *Bot) GetChatAdministrators(chatID int64) ([]ChatMember, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("getChatAdministrators", params, nil)
	if err != nil {
		return nil, err
	}
	var members []ChatMember
	if err := parseResult(data, &members); err != nil {
		return nil, err
	}
	return members, nil
}

func (b *Bot) GetChatMembersCount(chatID int64) (int, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("getChatMembersCount", params, nil)
	if err != nil {
		return 0, err
	}
	var count int
	if err := parseResult(data, &count); err != nil {
		return 0, err
	}
	return count, nil
}

func (b *Bot) GetChatMember(chatID int64, userID int64) (*ChatMember, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"user_id": strconv.FormatInt(userID, 10),
	}

	data, err := b.Client.Request("getChatMember", params, nil)
	if err != nil {
		return nil, err
	}
	var member ChatMember
	if err := parseResult(data, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

func (b *Bot) BanChatMember(chatID int64, userID int64) (bool, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"user_id": strconv.FormatInt(userID, 10),
	}
	data, err := b.Client.Request("banChatMember", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) UnbanChatMember(chatID int64, userID int64, onlyIfBanned bool) (bool, error) {
	params := map[string]string{
		"chat_id":        strconv.FormatInt(chatID, 10),
		"user_id":        strconv.FormatInt(userID, 10),
		"only_if_banned": strconv.FormatBool(onlyIfBanned),
	}
	data, err := b.Client.Request("unbanChatMember", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) PromoteChatMember(chatID int64, userID int64, perms map[string]bool) (bool, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"user_id": strconv.FormatInt(userID, 10),
	}
	for k, v := range perms {
		params[k] = strconv.FormatBool(v)
	}

	data, err := b.Client.Request("promoteChatMember", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) RestrictChatMember(chatID int64, userID int64, permissions ChatPermissions, untilDate int64) (bool, error) {
	params := map[string]string{
		"chat_id":     strconv.FormatInt(chatID, 10),
		"user_id":     strconv.FormatInt(userID, 10),
		"permissions": mustMarshal(permissions),
	}
	if untilDate != 0 {
		params["until_date"] = strconv.FormatInt(untilDate, 10)
	}

	data, err := b.Client.Request("restrictChatMember", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) SetChatPermissions(chatID int64, permissions ChatPermissions) (bool, error) {
	params := map[string]string{
		"chat_id":     strconv.FormatInt(chatID, 10),
		"permissions": mustMarshal(permissions),
	}

	data, err := b.Client.Request("setChatPermissions", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) PinChatMessage(chatID, messageID int64, disableNotification bool) (bool, error) {
	params := map[string]string{
		"chat_id":              strconv.FormatInt(chatID, 10),
		"message_id":           strconv.FormatInt(messageID, 10),
		"disable_notification": strconv.FormatBool(disableNotification),
	}
	data, err := b.Client.Request("pinChatMessage", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) UnpinChatMessage(chatID, messageID int64) (bool, error) {
	params := map[string]string{
		"chat_id":    strconv.FormatInt(chatID, 10),
		"message_id": strconv.FormatInt(messageID, 10),
	}
	data, err := b.Client.Request("unpinChatMessage", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) UnpinAllChatMessages(chatID int64) (bool, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("unpinAllChatMessages", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) LeaveChat(chatID int64) (bool, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("leaveChat", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) SetChatTitle(chatID int64, title string) (bool, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
		"title":   title,
	}
	data, err := b.Client.Request("setChatTitle", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) SetChatDescription(chatID int64, description string) (bool, error) {
	params := map[string]string{
		"chat_id":     strconv.FormatInt(chatID, 10),
		"description": description,
	}
	data, err := b.Client.Request("setChatDescription", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) DeleteChatPhoto(chatID int64) (bool, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("deleteChatPhoto", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) SetChatPhoto(chatID int64, photo FileUpload) (bool, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	data, err := b.Client.RequestMultipart("setChatPhoto", params, map[string]FileUpload{"photo": photo})
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) CreateChatInviteLink(chatID int64, name string, expireDate int, memberLimit int) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id": strconv.FormatInt(chatID, 10),
	}
	if name != "" {
		params["name"] = name
	}
	if expireDate != 0 {
		params["expire_date"] = strconv.Itoa(expireDate)
	}
	if memberLimit != 0 {
		params["member_limit"] = strconv.Itoa(memberLimit)
	}

	data, err := b.Client.Request("createChatInviteLink", params, nil)
	if err != nil {
		return nil, err
	}
	var link ChatInviteLink
	if err := parseResult(data, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (b *Bot) RevokeChatInviteLink(chatID int64, inviteLink string) (*ChatInviteLink, error) {
	params := map[string]string{
		"chat_id":     strconv.FormatInt(chatID, 10),
		"invite_link": inviteLink,
	}
	data, err := b.Client.Request("revokeChatInviteLink", params, nil)
	if err != nil {
		return nil, err
	}
	var link ChatInviteLink
	if err := parseResult(data, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (b *Bot) ExportChatInviteLink(chatID int64) (string, error) {
	params := map[string]string{"chat_id": strconv.FormatInt(chatID, 10)}
	data, err := b.Client.Request("exportChatInviteLink", params, nil)
	if err != nil {
		return "", err
	}
	var link string
	if err := parseResult(data, &link); err != nil {
		return "", err
	}
	return link, nil
}

// --- File & Sticker ---

func (b *Bot) GetFile(fileID string) (*File, error) {
	params := map[string]string{"file_id": fileID}
	data, err := b.Client.Request("getFile", params, nil)
	if err != nil {
		return nil, err
	}
	var file File
	if err := parseResult(data, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (b *Bot) FileLink(file *File) string {
	return fmt.Sprintf(FileEndpoint, b.Client.Token, file.FilePath)
}

func (b *Bot) UploadStickerFile(userID int64, sticker FileUpload, stickerFormat string) (*File, error) {
	params := map[string]string{
		"user_id":        strconv.FormatInt(userID, 10),
		"sticker_format": stickerFormat,
	}
	data, err := b.Client.RequestMultipart("uploadStickerFile", params, map[string]FileUpload{"sticker": sticker})
	if err != nil {
		return nil, err
	}
	var f File
	if err := parseResult(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (b *Bot) CreateNewStickerSet(userID int64, name, title string, stickers []InputSticker) (bool, error) {
	body := struct {
		UserID   int64          `json:"user_id"`
		Name     string         `json:"name"`
		Title    string         `json:"title"`
		Stickers []InputSticker `json:"stickers"`
	}{
		UserID:   userID,
		Name:     name,
		Title:    title,
		Stickers: stickers,
	}
	data, err := b.Client.Request("createNewStickerSet", nil, body)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) AddStickerToSet(userID int64, name string, sticker InputSticker) (bool, error) {
	body := struct {
		UserID  int64        `json:"user_id"`
		Name    string       `json:"name"`
		Sticker InputSticker `json:"sticker"`
	}{
		UserID:  userID,
		Name:    name,
		Sticker: sticker,
	}
	data, err := b.Client.Request("addStickerToSet", nil, body)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (b *Bot) AskReview(userID int64, delaySeconds int) (bool, error) {
	params := map[string]string{
		"user_id":       strconv.FormatInt(userID, 10),
		"delay_seconds": strconv.Itoa(delaySeconds),
	}
	data, err := b.Client.Request("askReview", params, nil)
	if err != nil {
		return false, err
	}
	var ok bool
	if err := parseResult(data, &ok); err != nil {
		return false, err
	}
	return ok, nil
}
