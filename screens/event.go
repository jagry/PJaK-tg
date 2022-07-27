package screens

func NewEvent(screen Interface, text string) Event {
	return Event{Screen: screen, Text: text}
}

type Event struct {
	Screen Interface
	Text   string
}
