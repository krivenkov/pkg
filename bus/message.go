package bus

import "time"

type Message[T any] struct {
	Key   string
	Value MessageValue[T]
}

type MessageValue[T any] struct {
	CommandID string
	UserID    string
	CreatedAt time.Time
	Payload   *T
}

func (m *Message[T]) WriteCommandID(commandID string) {
	m.Value.CommandID = commandID
}

func (m *Message[T]) WriteUserID(userID string) {
	m.Value.UserID = userID
}

func (m *Message[T]) WriteCreatedAt(createdAt time.Time) {
	m.Value.CreatedAt = createdAt
}

func (m *Message[T]) WritePayload(payload *T) {
	m.Value.Payload = payload
}

func (m *Message[T]) WriteKey(key string) {
	m.Key = key
}

func (m *Message[T]) ReadCommandID() string {
	return m.Value.CommandID
}

func (m *Message[T]) ReadUserID() string {
	return m.Value.UserID
}

func (m *Message[T]) ReadCreatedAt() time.Time {
	return m.Value.CreatedAt
}

func (m *Message[T]) ReadPayload() *T {
	return m.Value.Payload
}

func (m *Message[T]) ReadKey() []byte {
	return []byte(m.Key)
}

func (m *Message[T]) ReadValue() MessageValue[T] {
	return m.Value
}
