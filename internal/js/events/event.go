package events

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warning"
	LevelError Level = "error"
)

type Event struct {
	Level   Level
	Message string
	Error   error
}

type Events []Event

func (e Events) EventsCountWithLevel(l Level) (result int) {
	for _, event := range e {
		if event.Level == l {
			result++
		}
	}

	return
}

func (e Events) HasEventsWithLevel(l Level) bool {
	for _, event := range e {
		if event.Level == l {
			return true
		}
	}

	return false
}
