package httpsrv

import (
	"sync"
	"testing"
)

type mockWatcher struct {
	id               string
	receivedMessages []string
	mu               sync.Mutex
}

func (m *mockWatcher) GetWatcherId() string {
	return m.id
}

func (m *mockWatcher) Send(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.receivedMessages = append(m.receivedMessages, msg)
}

func (m *mockWatcher) Messages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.receivedMessages
}

type mockServer struct {
	watchers     map[string]*mockWatcher
	watchersLock sync.RWMutex
}

func (s *mockServer) notifyWatchers(msg string) {
	s.watchersLock.RLock()
	defer s.watchersLock.RUnlock()
	for _, watcher := range s.watchers {
		watcher.Send(msg)
	}
}

func TestNotifyWatchers(t *testing.T) {
	server := &mockServer{
		watchers: make(map[string]*mockWatcher),
	}

	watcher1 := &mockWatcher{id: wid1}
	watcher2 := &mockWatcher{id: wid2}

	server.watchers[watcher1.GetWatcherId()] = watcher1
	server.watchers[watcher2.GetWatcherId()] = watcher2

	server.notifyWatchers(test_message)

	if len(watcher1.Messages()) != 1 || watcher1.Messages()[0] != test_message {
		t.Errorf("Expected watcher1 to receive %v, but got %v", test_message, watcher1.Messages())
	}

	if len(watcher2.Messages()) != 1 || watcher2.Messages()[0] != test_message {
		t.Errorf("Expected watcher2 to receive %v, but got %v", test_message, watcher2.Messages())
	}
}
