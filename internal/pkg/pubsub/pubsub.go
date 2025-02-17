package pubsub

import (
	"fmt"

	gSyncMap "github.com/Cheasezz/testForOzon/internal/pkg/gemericSyncMap"
)

type CommentEvent struct {
	PostID string
	Comment interface{}
}

type Subscriber chan CommentEvent

type PubSub struct {
	// ключ – postID, значение – список подписчиков
	subscribers *gSyncMap.GenericMap[string, []Subscriber] 
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: &gSyncMap.GenericMap[string, []Subscriber]{},
	}
}

// Subscribe подписывается на события для конкретного postID.
func (ps *PubSub) Subscribe(postID string) Subscriber {
	sub := make(Subscriber, 10)
	// Получаем текущее значение или создаем новое
	existing, loaded := ps.subscribers.LoadOrStore(postID, []Subscriber{sub})
	if loaded {
		// Если значение уже есть, то existing — это срез подписчиков.
		subs := existing
		// Создаем новый срез, добавляя нашего нового подписчика.
		newSubs := append(subs, sub)
		ps.subscribers.Store(postID, newSubs)
	}
	return sub
}

// Unsubscribe удаляет подписчика.
func (ps *PubSub) Unsubscribe(postID string, sub Subscriber) {

	value, err := ps.subscribers.Load(postID)
	if err != nil {
		return
	}
	subs := value
	for i, s := range subs {
		if s == sub {
			newSubs := append(subs[:i], subs[i+1:]...)
			ps.subscribers.Store(postID, newSubs)
			close(s)
			break
		}
	}
}

// Publish рассылает событие всем подписчикам, зарегистрированным на данный postID.
func (ps *PubSub) Publish(event CommentEvent) {
	value, err := ps.subscribers.Load(event.PostID)
	if err != nil {
		return
	}
	subs := value
	fmt.Printf("From Publish pubsub: %v", subs)
	for _, sub := range subs {
		select {
		case sub <- event:
			fmt.Println("Event sent")
		default:
			fmt.Println("Subscriber channel is full, skipping event")
		}
	}
}
