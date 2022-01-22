package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	file string
	logg *logrus.Logger
}

func New(file, level, formatter string) (*Logger, error) {
	log := logrus.New()

	switch file {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		fmt.Println("File: ", file)
		fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, fmt.Errorf("invalid log filename: %w", err)
		}

		log.SetOutput(fd)
	}

	levelID, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log.SetLevel(levelID)

	switch formatter {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text_simple":
		log.SetFormatter(&SimpleTextFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{})
	}

	logger := &Logger{
		file: file,
		logg: log,
	}

	return logger, nil
}

func (l *Logger) Debug(msg string, params ...interface{}) {
	l.logg.Debugf(msg, params...)
}

func (l *Logger) Info(msg string, params ...interface{}) {
	l.logg.Infof(msg, params...)
}

func (l *Logger) Warn(msg string, params ...interface{}) {
	l.logg.Warnf(msg, params...)
}

func (l *Logger) Error(msg string, params ...interface{}) {
	l.logg.Errorf(msg, params...)
}

func (l *Logger) LogHTTPRequest(r *http.Request, code, length int) {
	l.logg.Infof(
		"%s [%v] %s %s %s %d %d %q",
		r.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 +0600"),
		r.Method,
		r.URL.String(),
		r.Proto,
		code,
		length,
		r.UserAgent(),
	)
}

type SimpleTextFormatter struct{}

func (f *SimpleTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s\t%s\n", entry.Level, entry.Message)

	return []byte(msg), nil
}
