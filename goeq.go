// Package goeq provides a simple event queue that leverages Go's flavor of interfaces to connect consumers and producers of events.
package goeq

// A runner is used to run the event method in a certain context (for example in a certain goroutine).
type Runner func(func())

// The queue manages a list of consumers in order to notify them of certain events.
type Queue struct {
	consumers []any
	runner    Runner
}

// New returns a new event queue that runs the event method in the same goroutine as the `Publish` method is called.
func New() *Queue {
	return &Queue{}
}

// NewConfied returns a new event queue that uses the given runner to execute the event method.
func NewConfined(runner Runner) *Queue {
	return &Queue{
		runner: runner,
	}
}

// Subscribe registers the given consumer with this event queue.
func (q *Queue) Subscribe(consumer any) {
	q.consumers = append(q.consumers, consumer)
}

// Publish runs the given function for all subscribed consumers. Use `Publish` in conjunction with `Message` function
// to restrict execution to the relevant consumers (those that implement the relevant interface).
func (q *Queue) Publish(notify func(item any)) {
	for _, consumer := range q.consumers {
		q.run(func() {
			notify(consumer)
		})
	}
}

func (q *Queue) run(f func()) {
	if q.runner == nil {
		f()
		return
	}
	q.runner(f)
}

// Message is used to restrict the execution of the event method to those consumers that implement the relevant interface.
func Message[T any](f func(T)) func(any) {
	return func(item any) {
		t, ok := item.(T)
		if !ok {
			return
		}
		f(t)
	}
}

type goRunner struct {
	messages chan func()
}

// NewSyncRunner returns a runner that executes all event methods confined to the same goroutine.
// The runner returns after the event method was executed.
func NewSyncRunner() Runner {
	return newRunner(0)
}

// NewAsyncRunner returns a runner that executes all event methods confined to the same goroutine.
// The runner returns immediately and does not wait until the event method was executed.
func NewAsyncRunner() Runner {
	return newRunner(1)
}

func newRunner(buffer int) Runner {
	result := &goRunner{
		messages: make(chan func(), buffer),
	}

	go func() {
		for m := range result.messages {
			m()
		}
	}()

	return result.run
}

func (r *goRunner) run(f func()) {
	r.messages <- f
}
