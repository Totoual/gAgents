package gAgents

import (
	"sync"

	"golang.org/x/net/context"
)

type Event struct {
	Type         string
	Payload      interface{}
	ResponseChan chan interface{}
}

type EventDispatcher struct {
	subscribers map[string][]func(event Event)
	ctx         context.Context
	mu          sync.RWMutex
}

func NewEventDispatcher(ctx context.Context) *EventDispatcher {
	return &EventDispatcher{
		subscribers: make(map[string][]func(event Event)),
		ctx:         ctx,
	}
}

func (ed *EventDispatcher) Subscribe(eventType string, handler func(event Event)) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	ed.subscribers[eventType] = append(ed.subscribers[eventType], handler)
}

func (ed *EventDispatcher) Publish(event Event) {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

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
