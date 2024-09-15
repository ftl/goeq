package goeq_test

import (
	"fmt"

	"github.com/ftl/goeq"
)

type FrequencyListener interface {
	FrequencyChanged(frequency float64) // FrequencyListeners will be notified when the frequency changed.
}

type Transceiver struct {
	events *goeq.Queue
}

func (t *Transceiver) emitFrequencyChanged(frequency float64) {
	t.events.Publish(goeq.Message(func(listener FrequencyListener) {
		listener.FrequencyChanged(frequency)
	}))
}

type FrequencyLabel struct{}

func (l *FrequencyLabel) SetText(text string) {
	fmt.Println("Frequency:", text)
}

func (l *FrequencyLabel) FrequencyChanged(frequency float64) {
	l.SetText(fmt.Sprintf("%0.2f", frequency))
}

func Example() {
	events := goeq.New()
	trx := &Transceiver{ // trx will produce FrequencyChanged events
		events: events,
	}
	label := &FrequencyLabel{} // label will consume FrequencyChanged events
	events.Subscribe(label)

	trx.emitFrequencyChanged(14030000.0)

	// Output:
	// Frequency: 14030000.00
}
