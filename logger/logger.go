package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LogCtxKey int

const (
	LoggerCtxKey LogCtxKey = iota
)

var (
	LOG_KEY_CORRELATION_ID  = "correlationId"
	LOG_KEY_RESPONSE_TIME   = "responseTime"
	LOG_KEY_STACK_TRACE     = "stackTrace"
	LOG_KEY_ADDITIONAL_INFO = "additionalInfo"
	LOG_KEY_STATUS_CODE     = "status"
)

type LoggerWrapper interface {
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	AdditionalInfo(info map[string]interface{}) *loggerWrapper
	StackTrace(err error) *loggerWrapper
}

type loggerWrapper struct {
	logger         zerolog.Logger
	logbase        *zerolog.Event
	startTime      time.Time
	stackTrace     error
	messagePrefix  string
	additionalInfo map[string]interface{}
}

func Ctx(ctx context.Context) LoggerWrapper {
	logger := zerolog.Logger{}
	ctxVal := ctx.Value(LoggerCtxKey)
	if ctxVal != nil {
		logger = ctx.Value(LoggerCtxKey).(zerolog.Logger)
	}
	return &loggerWrapper{
		logger:         logger,
		startTime:      time.Now(),
		messagePrefix:  "",
		additionalInfo: nil,
		stackTrace:     nil,
	}
}

func ConfigureLogger(debugMode bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.LevelInfoValue = "INFO"
	zerolog.LevelDebugValue = "DEBUG"
	zerolog.LevelErrorValue = "ERROR"
	zerolog.LevelFatalValue = "FATAL"
	zerolog.LevelWarnValue = "WARN"
	zerolog.LevelPanicValue = "PANIC"
	zerolog.ErrorStackFieldName = "stackTrace"
}

func GetCorrelationIDLoggerCtx(ctx context.Context, cid string) context.Context {
	cidlog := zlog.With().Str(LOG_KEY_CORRELATION_ID, cid).Logger()
	return context.WithValue(ctx, LoggerCtxKey, cidlog)
}

func Debug(format string, args ...interface{}) {
	zlog.Debug().Msgf(format, args...)
}

// log info level without context
func Info(format string, args ...interface{}) {
	zlog.Info().Msgf(format, args...)
}

// log error level without context
func Error(format string, args ...interface{}) {
	zlog.Error().Msgf(format, args...)
}

// log fatal level without context
func Fatal(format string, args ...interface{}) {
	zlog.Fatal().Msgf(format, args...)
}

func (l *loggerWrapper) logFormatter(format string, args ...interface{}) {
	if l.messagePrefix != "" {
		format = l.messagePrefix + " " + format
	}

	if l.stackTrace != nil {
		l.AdditionalInfo(map[string]interface{}{
			LOG_KEY_STACK_TRACE: fmt.Sprintf("%+v", l.stackTrace),
		})
	}

	if len(l.additionalInfo) > 0 {
		l.logbase.Interface(LOG_KEY_ADDITIONAL_INFO, l.additionalInfo)
	}

	l.logbase.Msgf(format, args...)

	l.additionalInfo = nil
	l.stackTrace = nil
	l.messagePrefix = ""
}

func (l *loggerWrapper) Info(format string, args ...interface{}) {
	l.logbase = l.logger.Info()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Debug(format string, args ...interface{}) {
	l.logbase = l.logger.Debug()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Warn(format string, args ...interface{}) {
	l.logbase = l.logger.Warn()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Error(format string, args ...interface{}) {
	l.logbase = l.logger.Error()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) Fatal(format string, args ...interface{}) {
	l.logbase = l.logger.Fatal()
	l.logFormatter(format, args...)
}

func (l *loggerWrapper) AdditionalInfo(info map[string]interface{}) *loggerWrapper {
	if l.additionalInfo == nil {
		l.additionalInfo = info
	} else {
		for k, v := range info {
			l.additionalInfo[k] = v
		}
	}

	return l
}

func (l *loggerWrapper) StackTrace(err error) *loggerWrapper {
	l.stackTrace = err
	return l
}
