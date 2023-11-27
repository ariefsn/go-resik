package logger

import (
	"os"

	"github.com/ariefsn/go-resik/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	mode := os.Getenv("MODE")
	if mode == "" || mode == "DEBUG" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func toDict(m common.M) *zerolog.Event {
	res := zerolog.Dict()
	for k, v := range m {
	SWT:
		switch t := v.(type) {
		case string:
			res.Str(k, t)
			break SWT
		case int:
			res.Int(k, t)
			break SWT
		case int8:
			res.Int8(k, t)
			break SWT
		case int16:
			res.Int16(k, t)
			break SWT
		case int32:
			res.Int32(k, t)
			break SWT
		case int64:
			res.Int64(k, t)
			break SWT
		case bool:
			res.Bool(k, t)
			break SWT
		case float32:
			res.Float32(k, t)
			break SWT
		case float64:
			res.Float64(k, t)
			break SWT
		default:
			res.Interface(k, t)
			break SWT
		}
	}

	return res
}

func Info(message string, additionalInfo ...common.M) {
	l := log.Info()

	if len(additionalInfo) > 0 {
		l.Dict("additionalInfo", toDict(additionalInfo[0]))
	}

	l.Msg(message)
}

func Warning(message string, additionalInfo ...common.M) {
	l := log.Warn()

	if len(additionalInfo) > 0 {
		l.Dict("additionalInfo", toDict(additionalInfo[0]))
	}

	l.Msg(message)
}

func Error(err error, additionalInfo ...common.M) {
	l := log.Error().Stack()

	if len(additionalInfo) > 0 {
		l.Dict("additionalInfo", toDict(additionalInfo[0]))
	}

	l.Err(err).Msg("")
}

func Fatal(err error, additionalInfo ...common.M) {
	l := log.Fatal().Stack()

	if len(additionalInfo) > 0 {
		l.Dict("additionalInfo", toDict(additionalInfo[0]))
	}

	l.Err(err).Msg("")
}
