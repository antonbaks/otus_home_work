package logger

import "fmt"

const (
	INFO  = "INFO"
	ERROR = "ERROR"
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{
		level: level,
	}
}

func (l Logger) Info(msg string) {
	if l.level == INFO || l.level == ERROR {
		fmt.Println(msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level == ERROR {
		fmt.Println(msg)
	}
}
