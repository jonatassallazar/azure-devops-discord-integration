package notify

import (
	"context"
	"errors"
	"log"
	"sync"
)

// Dispatcher fans a Message out to every Sink configured for one event
// category (e.g. all sinks wired up for pull request notifications).
type Dispatcher struct {
	Sinks []Sink
}

// Send delivers msg to every sink concurrently. It returns an error only if
// every sink failed; a partial failure is logged and swallowed so a single
// outbound delivery hiccup doesn't fail the inbound webhook.
func (d *Dispatcher) Send(ctx context.Context, msg Message) error {
	if len(d.Sinks) == 0 {
		return nil
	}

	errsCh := make(chan error, len(d.Sinks))

	var wg sync.WaitGroup
	for _, sink := range d.Sinks {
		wg.Add(1)
		go func(s Sink) {
			defer wg.Done()
			errsCh <- s.Send(ctx, msg)
		}(sink)
	}
	wg.Wait()
	close(errsCh)

	var errs []error
	for err := range errsCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	switch {
	case len(errs) == 0:
		return nil
	case len(errs) < len(d.Sinks):
		log.Printf("notify: %d/%d sinks failed for %q: %v", len(errs), len(d.Sinks), msg.Title, errors.Join(errs...))
		return nil
	default:
		return errors.Join(errs...)
	}
}
