package logger

import (
	"io"
	"os"
	"path"

	"github.com/cxfksword/fnsync-desktop/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	consoleWriter := NewConsoleWriter(os.Stderr)
	// output log to file
	logFileWriter := createRollingLogFile()
	if logFileWriter != nil {
		formatWriter := NewFileFormatWriter(logFileWriter)
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, formatWriter)
		log.Logger = log.With().Caller().Logger().Output(multiWriter)
	} else {
		log.Logger = log.With().Caller().Logger().Output(consoleWriter)
	}
}

func createRollingLogFile() io.Writer {
	if config.App.Log == "" {
		log.Debug().Msg("Not specific log file path, will ignore")
		return nil
	}

	logDir := path.Dir(config.App.Log)
	if err := os.MkdirAll(logDir, 0744); err != nil {
		log.Error().Err(err).Msg("can't create log directory")
		return nil
	}

	log.Info().Msgf("Log file: %s", config.App.Log)
	return &lumberjack.Logger{
		Filename:   config.App.Log,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //daysdays
	}
}

func SetDebugLevel() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
