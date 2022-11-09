package fakeclientreader

import (
	"context"
	"time"
)

type Message struct {
	UserID int64
	Date   time.Time
	Text   string
}

type FakeClientReader struct {
	messages []Message
	duration time.Duration
}

func New(messages []Message, duration time.Duration) *FakeClientReader {
	return &FakeClientReader{
		messages: messages,
		duration: duration,
	}
}

func (c FakeClientReader) Read(ctx context.Context, callback func(context.Context, int64, time.Time, string)) {
	for _, message := range c.messages {
		time.Sleep(c.duration)

		callback(ctx, message.UserID, message.Date, message.Text)
	}
}
