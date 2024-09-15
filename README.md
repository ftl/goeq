# goeq - A Simple Event Queue for Go

The `goeq` event queue leverages Go's flavor of interfaces to connect consumers and producers of events. The queue can run the event method in a certain context (for example in a certain goroutine).

An event is just a method defined in an interface:
```go
type FrequencyListener interface {
    FrequencyChanged(frequency float64)
}
```

To emit an event, the producer provides a callback function which will be called by the event queue with an instance of the particular interface:
```go
func emitFrequencyChanged(frequency float64) {
    queue.Publish(goeq.Event(func(listener FrequencyListener) {
        listener.FrequencyChanged(frequency)
    }))
}
```

Consumers must have subscribed to the event queue before the event was emitted:
```go
func (l *FrequencyLabel) subscribe() {
    queue.Subscribe(l)
}
```

To receive an event, the consumer needs to implement the interface which is used by the consumer to emit the event:
```go
func (l *FrequencyLabel) FrequencyChanged(frequency float64) {
    l.SetText(fmt.Sprintf("%0.2f", frequency))
}
```

## License
This library is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)