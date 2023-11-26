package eventtypes

const (
	//custom based
	Message            = "Message"
	ReplyMessage       = "ReplyMessage"
	ReplyToBotText     = "ReplyToBotText"
	ReplyToBotCallback = "ReplyToBotCallback"
	CommandMessage     = "CommandMessage"
	EditedMessage      = "EditedMessage"
	ChannelPost        = "ChannelPost"
	MemberJoinInChat   = "MemberJoinInChat"
	LeaveChatMember    = "LeaveChatMember"
	Undefined          = "Undefined"
	//for midddlware
	AllUpdateTypes = "AllUpdateTypes"
	//services
	EditedChannelPost  = "EditedChannelPost"
	InlineQuery        = "InlineQuery"
	ChosenInlineResult = "ChosenInlineResult"
	CallbackQuery      = "CallbackQuery"
	ShippingQuery      = "ShippingQuery"
	PreCheckoutQuery   = "PreCheckoutQuery"
	Poll               = "Poll"
	PollAnswer         = "PollAnswer"
	MyChatMember       = "MyChatMember"
	ChatMember         = "ChatMember"
	ChatJoinRequest    = "ChatJoinRequest"
)
