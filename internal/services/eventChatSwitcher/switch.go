package eventchatswitcher

import "github.com/ponyCorp/rebornPony/internal/models"

type EventChatSwitcher struct {
	chats   map[int64]bool
	adapter Adapter
}
type Adapter interface {
	GetChatByChatID(id int64) (models.Chat, error)
}

func New(adapter Adapter) *EventChatSwitcher {
	return &EventChatSwitcher{
		chats:   make(map[int64]bool),
		adapter: adapter,
	}
}
func (e *EventChatSwitcher) ChatIsDisabled(chatId int64) (bool, error) {

	val, ok := e.chats[chatId]
	if ok {
		return val, nil
	}
	chat, err := e.adapter.GetChatByChatID(chatId)
	if err != nil {
		return false, err
	}
	e.chats[chatId] = chat.IgnoreEvents
	return chat.IgnoreEvents, nil
}

// RefreshChatInfo
func (e *EventChatSwitcher) RefreshChatInfo(chatId int64) error {

	chat, err := e.adapter.GetChatByChatID(chatId)
	if err != nil {
		return err
	}
	e.chats[chatId] = chat.IgnoreEvents
	return nil
}
