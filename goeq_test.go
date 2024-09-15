package goeq

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type IFA interface{ A() }

type IFB interface{ B() }

type IFC interface{ C() }

type TA struct{ result *[]string }

func (a *TA) A() { *a.result = append(*a.result, "A") }

type TB struct{ result *[]string }

func (b *TB) B() { *b.result = append(*b.result, "B") }

type TAB struct{ result *[]string }

func (ab *TAB) A() { *ab.result = append(*ab.result, "a") }
func (ab *TAB) B() { *ab.result = append(*ab.result, "b") }

type TABEmbedded struct {
	TA
	TB
}

func TestPublish(t *testing.T) {
	result := make([]string, 0)

	q := New()
	a := &TA{result: &result}
	b := &TB{result: &result}
	ab := &TAB{result: &result}
	abe := &TABEmbedded{TA: TA{result: &result}, TB: TB{result: &result}}

	q.Subscribe(a)
	q.Subscribe(b)
	q.Subscribe(ab)
	q.Subscribe(abe)

	q.Publish(Event(func(item IFA) {
		item.A()
	}))
	assert.Equal(t, []string{"A", "a", "A"}, result)

	q.Publish(Event(func(item IFB) {
		item.B()
	}))
	assert.Equal(t, []string{"A", "a", "A", "B", "b", "B"}, result)

	q.Publish(Event(func(item IFC) {
		item.C()
	}))
	assert.Equal(t, []string{"A", "a", "A", "B", "b", "B"}, result)
}

func TestPublishConfined(t *testing.T) {
	asyncRunner := NewAsyncRunner()
	runner := func(f func()) {
		time.Sleep(2 * time.Millisecond)
		asyncRunner(f)
	}

	wg := new(sync.WaitGroup)

	result := make([]string, 0)

	q := NewConfined(runner)
	a := &TA{result: &result}
	b := &TB{result: &result}
	ab := &TAB{result: &result}
	abe := &TABEmbedded{TA: TA{result: &result}, TB: TB{result: &result}}

	q.Subscribe(a)
	q.Subscribe(b)
	q.Subscribe(ab)
	q.Subscribe(abe)

	wg.Add(3)
	q.Publish(Event(func(item IFA) {
		item.A()
		wg.Done()
	}))
	assert.NotEqual(t, []string{"A", "a", "A"}, result)
	wg.Wait()
	assert.Equal(t, []string{"A", "a", "A"}, result)

	wg.Add(3)
	q.Publish(Event(func(item IFB) {
		item.B()
		wg.Done()
	}))
	assert.NotEqual(t, []string{"A", "a", "A", "B", "b", "B"}, result)
	wg.Wait()
	assert.Equal(t, []string{"A", "a", "A", "B", "b", "B"}, result)
}
