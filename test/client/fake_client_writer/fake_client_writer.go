package fakeclientwriter

import (
	"context"
)

type Client interface {
}

type Message struct {
	UserID int64
	Text   string
}

type FakeClientWriter struct {
	messages []Message
}

func New() *FakeClientWriter {
	return &FakeClientWriter{
		messages: make([]Message, 0),
	}
}

func (c *FakeClientWriter) Write(ctx context.Context, text string, userID int64) error {
	c.messages = append(c.messages, Message{
		UserID: userID,
		Text:   text,
	})

	return nil
}

func (c FakeClientWriter) GetMessages() []Message {
	return c.messages
}
