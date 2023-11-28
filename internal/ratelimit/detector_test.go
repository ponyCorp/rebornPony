package ratelimit

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewMessageRateLimiter(t *testing.T) {
	type args struct {
		limit           int
		intervalSeconds int
	}
	tests := []struct {
		name string
		args args
		want *MessageRateLimiter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessageRateLimiter(tt.args.limit, tt.args.intervalSeconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageRateLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageRateLimiter_IsRateLimited(t *testing.T) {

	m := NewMessageRateLimiter(10, 5)
	for i := 0; i < 11; i++ {
		m.AddMessage("test")
	}

	if !m.IsRateLimited("test") {
		t.Errorf("MessageRateLimiter.IsRateLimited() = %v, want %v", false, true)
	}
	// type fields struct {
	// 	limit    int
	// 	interval int
	// }
	// type args struct {
	// 	msgIdentifier string
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	args   args
	// 	want   bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		m := NewMessageRateLimiter(tt.fields.limit, tt.fields.interval)
	// 		if got := m.IsRateLimited(tt.args.msgIdentifier); got != tt.want {
	// 			t.Errorf("MessageRateLimiter.IsRateLimited() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
func TestMessageRateLimiter_RateDetectWithCron(t *testing.T) {
	m := NewMessageRateLimiter(10, 1)
	wg := &sync.WaitGroup{}
	don := make(chan bool)
	wg.Add(3)
	go func() {
		m.StartCollector()
		wg.Done()
	}()
	go func() {
		for {

			select {
			case <-don:
				fmt.Println("done collector")
				defer wg.Done()
				return
			case <-time.After(10 * time.Millisecond):
				//for i := 0; i < 111; i++ {
				if m.IsRateLimited("test") {
					fmt.Println("detect rate limit. waiting...")
					time.Sleep(1 * time.Second)
					fmt.Println("continue...")
					continue
				}
				m.AddMessage("test")
				//m.AddMessage("test")
				//m.AddMessage("test")
				//m.AddMessage("test")
				fmt.Println("add messages")
			}
		}
		//}

	}()
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("stop collector")
		m.StopCollector()
		don <- true
		wg.Done()
	}()
	wg.Wait()

}
