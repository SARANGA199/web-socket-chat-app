package wsmanager

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	connections map[string][]*websocket.Conn
	mutex       sync.RWMutex
}

var Manager = &WebSocketManager{
	connections: make(map[string][]*websocket.Conn),
}

// AddConnection adds a new WebSocket connection for a user.
func (m *WebSocketManager) AddConnection(userID string, conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.connections[userID] = append(m.connections[userID], conn)
}

// RemoveConnection removes a specific WebSocket connection for a user.
func (m *WebSocketManager) RemoveConnection(userID string, conn *websocket.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	connections := m.connections[userID]
	for i, c := range connections {
		if c == conn {
			// Close the connection before removing it
			c.Close()
			m.connections[userID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// Remove the user key if no connections remain
	if len(m.connections[userID]) == 0 {
		delete(m.connections, userID)
	}
}

// BroadcastMessage sends a message to all WebSocket connections of a user.
func (m *WebSocketManager) BroadcastMessage(userID string, message interface{}) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, conn := range m.connections[userID] {
		if err := conn.WriteJSON(message); err != nil {
			// Handle error (closing and removing the bad connection)
			conn.Close()
		}
	}
}
