package ParsBale

type ChatAction string

const (
	ChatActionTyping        ChatAction = "typing"
	ChatActionUploadPhoto   ChatAction = "upload_photo"
	ChatActionRecordVideo   ChatAction = "record_video"
	ChatActionUploadVideo   ChatAction = "upload_video"
	ChatActionRecordVoice   ChatAction = "record_voice"
	ChatActionUploadVoice   ChatAction = "upload_voice"
	ChatActionChooseSticker ChatAction = "choose_sticker"
)

type ParseMode string

const (
	ModeMarkdown ParseMode = "Markdown"
	ModeHTML     ParseMode = "HTML"
)

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type MessageEntityType string

const (
	EntityMention    MessageEntityType = "mention"
	EntityBotCommand MessageEntityType = "bot_command"
	EntityURL        MessageEntityType = "url"
	EntityEmail      MessageEntityType = "email"
	EntityBold       MessageEntityType = "bold"
	EntityItalic     MessageEntityType = "italic"
	EntityCode       MessageEntityType = "code"
	EntityPre        MessageEntityType = "pre"
	EntityTextLink   MessageEntityType = "text_link"
)
