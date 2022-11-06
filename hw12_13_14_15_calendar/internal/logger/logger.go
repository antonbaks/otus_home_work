package logger

import (
	"os"
)

const (
	INFO  = "INFO"
	ERROR = "ERROR"
)

type Logger struct {
	level  string
	stderr *os.File
	stdout *os.File
}

func New(level string, stderr *os.File, stdout *os.File) *Logger {
	return &Logger{
		level:  level,
		stderr: stderr,
		stdout: stdout,
	}
}

func (l Logger) Info(msg string) {
	if l.level == INFO || l.level == ERROR {
		l.stdout.WriteString(msg + "\n")
	}
}

func (l Logger) Error(msg string) {
	if l.level == ERROR {
		l.stderr.WriteString(msg + "\n")
	}
}
