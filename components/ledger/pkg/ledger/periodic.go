package ledger

import (
	"context"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

type periodic struct {
	stopChan chan chan struct{}
	sem      chan struct{}
	fn       func(ctx context.Context) error
}

func (pc *periodic) Trigger(ctx context.Context) error {
	logging.FromContext(ctx).Debugf("Trigger new")
	defer func() {
		logging.FromContext(ctx).Debugf("Terminated!")
	}()
	defer func() {
		pc.sem <- struct{}{}
	}()
	err := pc.fn(ctx)
	if err != nil {
		logging.FromContext(ctx).Error("Error with procedure: %s", err)
	}
	return err
}

func (pc *periodic) Run(ctx context.Context) error {
	for {
		go pc.Trigger(ctx)
		select {
		case <-pc.sem:
			select {
			case <-time.After(500 * time.Millisecond):
			case ch := <-pc.stopChan:
				close(ch)
				return nil
			}
		case ch := <-pc.stopChan:
			close(ch)
			return nil
		}
	}
	return nil
}

func (pc *periodic) Stop() {
	ch := make(chan struct{})
	pc.stopChan <- ch
	<-ch
}

func newPeriodic(fn func(ctx context.Context) error) *periodic {
	return &periodic{
		fn:       fn,
		stopChan: make(chan chan struct{}),
		sem:      make(chan struct{}, 1),
	}
}
