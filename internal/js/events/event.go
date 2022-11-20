package events

type Event struct {
	Level   Level
	Message string
	Error   error
}

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warning"
	LevelError Level = "error"
)
