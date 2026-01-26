package events

import (
	"fmt"
	"sync"
)

type Message struct{
	Topic string
	Data any
}

type Subscriber chan Message

type Broker struct {
	subscribers map[string]map[Subscriber]bool
	mu sync.RWMutex
}

func NewBroker() *Broker {

	return &Broker{
		subscribers: make(map[string]map[Subscriber]bool),
	}
}

func (b *Broker) Subscribe(topic string) (sub Subscriber, unsub func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	sub = make(Subscriber, 10)

	if _, isTopicInitialized := b.subscribers[topic]; !isTopicInitialized {
		b.subscribers[topic] = make(map[Subscriber]bool)
	}

	b.subscribers[topic][sub] = true

	// Create an unsubscribe function that closes over the channel and topic
	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		// Check if the topic and subscriber still exist
		if subscribers, found := b.subscribers[topic]; found {
			if _, ok := subscribers[sub]; ok {
				// Remove the subscriber from the map
				delete(subscribers, sub)
				// Close the subannel to signal the subscriber's goroutine to exit
				close(sub)
			}
		}
	}

	return sub, unsubscribe
}

func (b *Broker) Publish(topic string, data any) {
	b.mu.RLock() // Use a Read-lock for publishing
	defer b.mu.RUnlock()

	msg := Message{Topic: topic, Data: data}

	if subscribers, found := b.subscribers[topic]; found {
		for ch := range subscribers {
			go func(c Subscriber) {
				select {
				case c <- msg:
					// Message sent
				default:
					fmt.Printf("WARNING: Subscriber channel full, blocking broker while waiting for space, topic: %s\n", topic)
				}
			}(ch)
		}
	}
}

