package pubsub

//go:generate mockgen -package pubsub -source=pubsub.go -destination=mocks_pubsub.go

import (
	"fmt"

	"github.com/Cheasezz/testForOzon/pkg/gSyncMap"
)

type IPubSub interface {
	Subscribe(keyId string) Subscriber
	Unsubscribe(keyId string, sub Subscriber)
	Publish(event CommentEvent)
}

type CommentEvent struct {
	KeyId   string
	Comment interface{}
}

type Subscriber chan CommentEvent

type PubSub struct {
	// ключ – postID, значение – список подписчиков
	subscribers *gSyncMap.GSyncMap[string, []Subscriber]
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: gSyncMap.NewGenericSyncMap[string, []Subscriber](),
	}
}

// Subscribe подписывается на события для конкретного postID.
func (ps *PubSub) Subscribe(keyId string) Subscriber {
	sub := make(Subscriber, 10)
	// Получаем текущее значение или создаем новое
	existing, loaded := ps.subscribers.LoadOrStore(keyId, []Subscriber{sub})
	if loaded {
		// Если значение уже есть, то existing — это срез подписчиков.
		subs := existing
		// Создаем новый срез, добавляя нашего нового подписчика.
		newSubs := append(subs, sub)
		ps.subscribers.Store(keyId, newSubs)
	}
	return sub
}

// Unsubscribe удаляет подписчика.
func (ps *PubSub) Unsubscribe(keyId string, sub Subscriber) {

	value, err := ps.subscribers.Load(keyId)
	if err != nil {
		return
	}
	subs := value
	for i, s := range subs {
		if s == sub {
			newSubs := append(subs[:i], subs[i+1:]...)
			ps.subscribers.Store(keyId, newSubs)
			close(s)
			break
		}
	}
}

// Publish рассылает событие всем подписчикам, зарегистрированным на данный postID.
func (ps *PubSub) Publish(event CommentEvent) {
	value, err := ps.subscribers.Load(event.KeyId)
	if err != nil {
		return
	}
	subs := value

	for _, sub := range subs {
		select {
		case sub <- event:
			fmt.Println("Event sent")
		default:
			fmt.Println("Subscriber channel is full, skipping event")
		}
	}
}
