package ratelimit

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Message struct {
	UserID string
	TimeAt time.Time
	ChatID string
}
type MessageRateLimiter struct {
	messageMap map[string][]int64
	mutex      sync.RWMutex
	limit      int
	interval   int
	done       chan bool
}

func GetKey(userID string, chatID string) string {
	return strings.Join([]string{userID, chatID}, ":")
}

// AddMessage Add message to map of messages for chatId or chatId:userId
func (m *MessageRateLimiter) AddMessage(msgIdentifier string) {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	if _, ok := m.messageMap[msgIdentifier]; !ok {
		m.messageMap[msgIdentifier] = make([]int64, 0)
	}
	m.messageMap[msgIdentifier] = append(m.messageMap[msgIdentifier], time.Now().Unix())
}

// cron Removes old messages
func (m *MessageRateLimiter) garbageCollector() {
	keys := make([]string, 0)
	m.mutex.RLock()
	for k := range m.messageMap {
		keys = append(keys, k)
	}
	m.mutex.RUnlock()
	for _, k := range keys {
		m.mutex.Lock()
		var newMk []int64
		for _, v := range m.messageMap[k] {
			//i := i
			v := v
			if !(v+int64(m.interval) < time.Now().Unix()) {
				newMk = append(newMk, v)
			}
		}
		m.messageMap[k] = newMk
		m.mutex.Unlock()
	}
}

// StartCollector
func (m *MessageRateLimiter) StartCollector() {
	ticker := time.NewTicker(time.Duration(m.interval) * time.Second)

	for {
		select {
		case <-m.done:

			ticker.Stop()
			fmt.Println("Ticker stopped")
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
			m.garbageCollector()
		}
	}

}
func (m *MessageRateLimiter) StopCollector() {
	m.done <- true
}
func (m *MessageRateLimiter) IsRateLimited(msgIdentifier string) bool {
	defer m.mutex.RUnlock()
	m.mutex.RLock()
	if _, ok := m.messageMap[msgIdentifier]; !ok {
		return false
	}
	l := len(m.messageMap[msgIdentifier])
	fmt.Printf("len: %v\n", l)
	return m.limit < l

}
func NewMessageRateLimiter(limit int, intervalSeconds int) *MessageRateLimiter {
	return &MessageRateLimiter{
		limit:      limit,
		interval:   intervalSeconds,
		messageMap: make(map[string][]int64),
		done:       make(chan bool),
	}
}
