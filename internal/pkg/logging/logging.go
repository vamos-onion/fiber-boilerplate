package logging

import (
	"context"
	"io"
	"os"

	"fiber-boilerplate/internal/pkg/util"

	flog "github.com/gofiber/fiber/v2/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger :
var Logger *zerolog.Logger

// FiberLogger :
var FiberLogger flog.AllLogger

func write(zEvent *zerolog.Event, skip int, format string, values ...interface{}) {
	if len(values) > 0 {
		zEvent.Caller(skip).Msgf(format, values...)
	} else {
		zEvent.Caller(skip).Msg(format)
	}
}

// Trace :
func Trace(format string, values ...interface{}) {
	write(Logger.Trace().Str("func", util.Prev.Func()), 2, format, values...)
}

// TraceSQL :
func TraceSQL(format string, values ...interface{}) {
	f := util.PPrev.Func()
	write(Logger.Trace().Str("func", f), 4, util.String.Concat("[", f, "] ", format), values...)
}

// Debug :
func Debug(format string, values ...interface{}) {
	write(Logger.Debug().Str("func", util.Prev.Func()), 2, format, values...)
}

// Info :
func Info(format string, values ...interface{}) {
	write(Logger.Info().Str("func", util.Prev.Func()), 2, format, values...)
}

// Warn :
func Warn(err error, format string, values ...interface{}) {
	write(Logger.Warn().Str("func", util.Prev.Func()).Err(err), 2, format, values...)
}

// Error :
func Error(err error, format string, values ...interface{}) {
	write(Logger.Error().Str("func", util.Prev.Func()).Err(err), 2, format, values...)
}

// Fatal :
func Fatal(err error, format string, values ...interface{}) {
	write(Logger.Fatal().Str("func", util.Prev.Func()).Err(err), 2, format, values...)
}

// Panic :
func Panic(err error, format string, values ...interface{}) {
	write(Logger.Panic().Str("func", util.Prev.Func()).Err(err), 2, format, values...)
}

// zerolog <-> fiber/log 어댑터
type fiberZL struct{ l zerolog.Logger }

// Logger (plain)
func (f fiberZL) Trace(v ...interface{}) { f.l.Trace().Msgf("%v", v...) }
func (f fiberZL) Debug(v ...interface{}) { f.l.Debug().Msgf("%v", v...) }
func (f fiberZL) Info(v ...interface{})  { f.l.Info().Msgf("%v", v...) }
func (f fiberZL) Warn(v ...interface{})  { f.l.Warn().Msgf("%v", v...) }
func (f fiberZL) Error(v ...interface{}) { f.l.Error().Msgf("%v", v...) }
func (f fiberZL) Fatal(v ...interface{}) { f.l.Fatal().Msgf("%v", v...) }
func (f fiberZL) Panic(v ...interface{}) { f.l.Panic().Msgf("%v", v...) }

// FormatLogger (*f)
func (f fiberZL) Tracef(s string, v ...interface{}) { f.l.Trace().Msgf(s, v...) }
func (f fiberZL) Debugf(s string, v ...interface{}) { f.l.Debug().Msgf(s, v...) }
func (f fiberZL) Infof(s string, v ...interface{})  { f.l.Info().Msgf(s, v...) }
func (f fiberZL) Warnf(s string, v ...interface{})  { f.l.Warn().Msgf(s, v...) }
func (f fiberZL) Errorf(s string, v ...interface{}) { f.l.Error().Msgf(s, v...) }
func (f fiberZL) Fatalf(s string, v ...interface{}) { f.l.Fatal().Msgf(s, v...) }
func (f fiberZL) Panicf(s string, v ...interface{}) { f.l.Panic().Msgf(s, v...) }

// WithLogger (*w)
func (f fiberZL) Tracew(msg string, kv ...interface{}) { withKV(f.l.Trace(), kv...).Msg(msg) }
func (f fiberZL) Debugw(msg string, kv ...interface{}) { withKV(f.l.Debug(), kv...).Msg(msg) }
func (f fiberZL) Infow(msg string, kv ...interface{})  { withKV(f.l.Info(), kv...).Msg(msg) }
func (f fiberZL) Warnw(msg string, kv ...interface{})  { withKV(f.l.Warn(), kv...).Msg(msg) }
func (f fiberZL) Errorw(msg string, kv ...interface{}) { withKV(f.l.Error(), kv...).Msg(msg) }
func (f fiberZL) Fatalw(msg string, kv ...interface{}) { withKV(f.l.Fatal(), kv...).Msg(msg) }
func (f fiberZL) Panicw(msg string, kv ...interface{}) { withKV(f.l.Panic(), kv...).Msg(msg) }

// SetLevel :
func (f fiberZL) SetLevel(lv flog.Level) {
	switch lv {
	case flog.LevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case flog.LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case flog.LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case flog.LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case flog.LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case flog.LevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case flog.LevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
}

// SetOutput :
func (f *fiberZL) SetOutput(w io.Writer) { f.l = f.l.Output(w) }

func (f fiberZL) WithContext(ctx context.Context) flog.CommonLogger { return f }

func withKV(e *zerolog.Event, kv ...interface{}) *zerolog.Event {
	for i := 0; i+1 < len(kv); i += 2 {
		if k, ok := kv[i].(string); ok {
			e = e.Interface(k, kv[i+1])
		}
	}
	return e
}

func init() {
	var writer io.Writer = os.Stdout
	zLogger := zerolog.New(os.Stdout)
	zLogger = log.Output(zerolog.ConsoleWriter{Out: writer})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = zLogger.With().Timestamp().Logger()
	Logger = &log.Logger
	FiberLogger = &fiberZL{l: *Logger}
	flog.SetLogger(FiberLogger)
}
