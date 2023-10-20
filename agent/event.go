package gAgents

import (
	"sync"

	"log"

	"golang.org/x/net/context"
)

type EventType string

type Event struct {
	Type         EventType
	Payload      interface{}
	ResponseChan chan interface{}
}

type EventDispatcher struct {
	subscribers map[EventType][]func(event Event)
	ctx         context.Context
	mu          sync.RWMutex
}

func NewEventDispatcher(ctx context.Context) *EventDispatcher {
	return &EventDispatcher{
		subscribers: make(map[EventType][]func(event Event)),
		ctx:         ctx,
	}
}

func (ed *EventDispatcher) Subscribe(eventType EventType, handler func(event Event)) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.subscribers[eventType] = append(ed.subscribers[eventType], handler)
}

func (ed *EventDispatcher) Publish(event Event) {
	ed.mu.RLock()
	defer ed.mu.RUnlock()
	log.Printf("New event published: %s", event.Type)
	handlers, found := ed.subscribers[event.Type]
	if !found {
		return
	}

	for _, handler := range handlers {

		localHandler := handler
		go func() {
			select {
			case <-ed.ctx.Done():
				return
			default:
				localHandler(event)
			}
		}()
	}
}
