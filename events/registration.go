package events

type EventHandler interface {
    Handle(c interface{}, evt Event)
}

var handlers = map[string]EventHandler{}

func Register(eventType string, h EventHandler) {
    handlers[eventType] = h
}

func GetHandler(eventType string) EventHandler {
    return handlers[eventType]
}
