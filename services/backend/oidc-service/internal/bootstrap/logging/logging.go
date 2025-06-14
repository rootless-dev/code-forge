package logging

import (
	"github.com/phuslu/log"
	"os"
)

func init() {
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeField:  "date",
		TimeFormat: "2006-01-02",
		Writer:     &log.IOWriter{Writer: os.Stdout},
	}

	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}
}
