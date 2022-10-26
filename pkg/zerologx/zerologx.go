package zerologx

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Zerolog struct {
	logger *zerolog.Logger
}

func New(level string, w io.Writer) Zerolog {
	if w == nil {
		w = os.Stdout
	}
	var lvl zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		lvl = zerolog.ErrorLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "debug":
		lvl = zerolog.DebugLevel
	default:
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)

	skipFrameCount := 3
	l := zerolog.New(w).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return Zerolog{
		logger: &l,
	}
}

func (l *Zerolog) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Zerolog) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Zerolog) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Zerolog) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.msg("error", message, args...)
}

func (l *Zerolog) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Zerolog) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Zerolog) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
