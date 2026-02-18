package ParsBale

import "encoding/json"

// --- Updates ---

type Update struct {
	UpdateID         int               `json:"update_id"`
	Message          *Message          `json:"message,omitempty"`
	EditedMessage    *Message          `json:"edited_message,omitempty"`
	CallbackQuery    *CallbackQuery    `json:"callback_query,omitempty"`
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
}

// --- Base Types ---

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID        int64    `json:"id"`
	Type      ChatType `json:"type"`
	Title     string   `json:"title,omitempty"`
	Username  string   `json:"username,omitempty"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
}

type ChatFullInfo struct {
	Chat
	Bio           string           `json:"bio,omitempty"`
	Photo         *ChatPhoto       `json:"photo,omitempty"`
	Description   string           `json:"description,omitempty"`
	InviteLink    string           `json:"invite_link,omitempty"`
	PinnedMessage *Message         `json:"pinned_message,omitempty"`
	Permissions   *ChatPermissions `json:"permissions,omitempty"`
}

// --- Message ---

type Message struct {
	MessageID            int64                 `json:"message_id"`
	From                 *User                 `json:"from,omitempty"`
	SenderChat           *Chat                 `json:"sender_chat,omitempty"`
	Date                 int64                 `json:"date"`
	Chat                 *Chat                 `json:"chat"`
	ForwardFrom          *User                 `json:"forward_from,omitempty"`
	ForwardFromChat      *Chat                 `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID int64                 `json:"forward_from_message_id,omitempty"`
	ForwardDate          int64                 `json:"forward_date,omitempty"`
	EditDate             int64                 `json:"edit_date,omitempty"`      // اضافه شده
	MediaGroupID         string                `json:"media_group_id,omitempty"` // اضافه شده
	Text                 string                `json:"text,omitempty"`
	Entities             []MessageEntity       `json:"entities,omitempty"`
	Photo                []PhotoSize           `json:"photo,omitempty"`
	Video                *Video                `json:"video,omitempty"`
	Audio                *Audio                `json:"audio,omitempty"`
	Document             *Document             `json:"document,omitempty"`
	Voice                *Voice                `json:"voice,omitempty"`
	Animation            *Animation            `json:"animation,omitempty"`
	Sticker              *Sticker              `json:"sticker,omitempty"`
	Contact              *Contact              `json:"contact,omitempty"`
	Location             *Location             `json:"location,omitempty"`
	Invoice              *Invoice              `json:"invoice,omitempty"`
	SuccessfulPayment    *SuccessfulPayment    `json:"successful_payment,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	ReplyToMessage       *Message              `json:"reply_to_message,omitempty"`
	Caption              string                `json:"caption,omitempty"`
	NewChatMembers       []User                `json:"new_chat_members,omitempty"`
	LeftChatMember       *User                 `json:"left_chat_member,omitempty"`
	WebAppData           *WebAppData           `json:"web_app_data,omitempty"`
}

type MessageEntity struct {
	Type   MessageEntityType `json:"type"`
	Offset int               `json:"offset"`
	Length int               `json:"length"`
	URL    string            `json:"url,omitempty"`
	User   *User             `json:"user,omitempty"`
}

// --- Media Types ---

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Title        string `json:"title,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	Performer    string `json:"performer,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type Animation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type Sticker struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Type         string     `json:"type"` // "regular", "mask"
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	IsAnimated   bool       `json:"is_animated"`
	IsVideo      bool       `json:"is_video"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	Emoji        string     `json:"emoji,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

type StickerSet struct {
	Name      string     `json:"name"`
	Title     string     `json:"title"`
	Stickers  []Sticker  `json:"stickers"`
	Thumbnail *PhotoSize `json:"thumbnail,omitempty"`
}

type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// --- Input Media ---

type InputMedia interface {
	MediaType() string
}

type InputMediaPhoto struct {
	Type    string `json:"type"` // "photo"
	Media   string `json:"media"`
	Caption string `json:"caption,omitempty"`
}

func (InputMediaPhoto) MediaType() string { return "photo" }

type InputMediaVideo struct {
	Type      string `json:"type"` // "video"
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (InputMediaVideo) MediaType() string { return "video" }

type InputMediaAnimation struct {
	Type      string `json:"type"` // "animation"
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (InputMediaAnimation) MediaType() string { return "animation" }

type InputMediaAudio struct {
	Type      string `json:"type"` // "audio"
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Title     string `json:"title,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (InputMediaAudio) MediaType() string { return "audio" }

type InputMediaDocument struct {
	Type      string `json:"type"` // "document"
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

func (InputMediaDocument) MediaType() string { return "document" }

// --- Sticker Input ---

type InputSticker struct {
	Sticker  string   `json:"sticker"` // file_id or upload URL
	Format   string   `json:"format"`  // "static", "animated", "video"
	Keywords []string `json:"keywords"`
	Emoji    string   `json:"emoji"`
}

// --- Other Types ---

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter,omitempty"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

type SuccessfulPayment struct {
	Currency                string `json:"currency"`
	TotalAmount             int    `json:"total_amount"`
	InvoicePayload          string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
}

type WebAppData struct {
	Data string `json:"data"`
}

// --- Keyboards ---

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string          `json:"text"`
	URL          string          `json:"url,omitempty"`
	CallbackData string          `json:"callback_data,omitempty"`
	WebApp       *WebAppInfo     `json:"web_app,omitempty"`
	CopyText     *CopyTextButton `json:"copy_text,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	Selective       bool               `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text            string      `json:"text"`
	RequestContact  bool        `json:"request_contact,omitempty"`
	RequestLocation bool        `json:"request_location,omitempty"`
	WebApp          *WebAppInfo `json:"web_app,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type CopyTextButton struct {
	Text string `json:"text"`
}

// --- Payment & Queries ---

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message,omitempty"`
	Data    string   `json:"data,omitempty"`
}

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	From           *User  `json:"from"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}

type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"` // ریال
}

// --- Chat Admin ---

type ChatMember struct {
	User                *User  `json:"user"`
	Status              string `json:"status"` // creator, administrator, member, restricted, left, kicked
	CanChangeInfo       bool   `json:"can_change_info,omitempty"`
	CanPostMessages     bool   `json:"can_post_messages,omitempty"`
	CanEditMessages     bool   `json:"can_edit_messages,omitempty"`
	CanDeleteMessages   bool   `json:"can_delete_messages,omitempty"`
	CanInviteUsers      bool   `json:"can_invite_users,omitempty"`
	CanRestrictMembers  bool   `json:"can_restrict_members,omitempty"`
	CanPinMessages      bool   `json:"can_pin_messages,omitempty"`
	CanPromoteMembers   bool   `json:"can_promote_members,omitempty"`
	CanManageVideoChats bool   `json:"can_manage_video_chats,omitempty"`
	CanManageChat       bool   `json:"can_manage_chat,omitempty"`
	CanPostStories      bool   `json:"can_post_stories,omitempty"` // اضافه شده
	IsMember            bool   `json:"is_member,omitempty"`
	CanSendMessages     bool   `json:"can_send_messages,omitempty"`
	CanSendAudios       bool   `json:"can_send_audios,omitempty"`
	CanSendDocuments    bool   `json:"can_send_documents,omitempty"`
	CanSendPhotos       bool   `json:"can_send_photos,omitempty"`
	CanSendVideos       bool   `json:"can_send_videos,omitempty"`
}

type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatPermissions ساختار دسترسی‌ها طبق مستندات رسمی بله
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendAudios         bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`
	CanSendVideos         bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
}

type ChatInviteLink struct {
	InviteLink  string `json:"invite_link"`
	Creator     *User  `json:"creator"`
	IsPrimary   bool   `json:"is_primary"`
	IsRevoked   bool   `json:"is_revoked"`
	Name        string `json:"name,omitempty"`
	ExpireDate  int    `json:"expire_date,omitempty"`
	MemberLimit int    `json:"member_limit,omitempty"`
}

type ResponseParameters struct {
	RetryAfter int `json:"retry_after,omitempty"`
}

type WebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	IPAddress            string `json:"ip_address,omitempty"`
	LastErrorDate        int    `json:"last_error_date,omitempty"`
	LastErrorMessage     string `json:"last_error_message,omitempty"`
}

type MessageId struct {
	MessageID int64 `json:"message_id"`
}

// --- Helpers ---

func NewInlineKeyboard(rows ...[]InlineKeyboardButton) *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

func NewKeyboard(rows ...[]KeyboardButton) *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{Keyboard: rows}
}

func mustMarshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// --- Transaction for Payment ---
type Transaction struct {
	ID        string `json:"id"`
	Status    string `json:"status"` // pending, paid, failed, rejected
	UserID    int64  `json:"userId"`
	Amount    int    `json:"amount"`
	CreatedAt int    `json:"createdAt"`
}
