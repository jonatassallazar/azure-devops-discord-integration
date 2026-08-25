package notify

import "context"

// Sink delivers a Message to one outbound destination, e.g. Discord or
// Google Chat.
type Sink interface {
	Send(ctx context.Context, msg Message) error
}
