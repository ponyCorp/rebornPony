package eventtypes

type Event string

const (
	//custom based
	Message            Event = "Message"
	ReplyMessage       Event = "ReplyMessage"
	ReplyToBotText     Event = "ReplyToBotText"
	ReplyToBotCallback Event = "ReplyToBotCallback"
	CommandMessage     Event = "CommandMessage"
	EditedMessage      Event = "EditedMessage"
	ChannelPost        Event = "ChannelPost"
	MemberJoinInChat   Event = "MemberJoinInChat"
	LeaveChatMember    Event = "LeaveChatMember"
	Undefined          Event = "Undefined"
	//for midddlware
	AllUpdateTypes Event = "AllUpdateTypes"
	//services
	EditedChannelPost  Event = "EditedChannelPost"
	InlineQuery        Event = "InlineQuery"
	ChosenInlineResult Event = "ChosenInlineResult"
	CallbackQuery      Event = "CallbackQuery"
	ShippingQuery      Event = "ShippingQuery"
	PreCheckoutQuery   Event = "PreCheckoutQuery"
	Poll               Event = "Poll"
	PollAnswer         Event = "PollAnswer"
	MyChatMember       Event = "MyChatMember"
	ChatMember         Event = "ChatMember"
	ChatJoinRequest    Event = "ChatJoinRequest"
)

func (e Event) String() string {
	return string(e)
}
