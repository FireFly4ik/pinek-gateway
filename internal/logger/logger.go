package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(productionType string) *zerolog.Logger {
	if productionType == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = "15:04:05 02.01.2006"

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		parts := strings.Split(file, "/")
		if len(parts) > 2 {
			file = strings.Join(parts[len(parts)-2:], "/")
		}
		return fmt.Sprintf("%s:%d", file, line)
	}

	var writer io.Writer = os.Stdout

	loggerContext := zerolog.New(writer).
		With().
		Caller().
		Timestamp().
		Logger()

	log.Info().Msg("logger setup complete")
	return &loggerContext
}

func RequestLogger(serviceName string) gin.HandlerFunc {
	logger := log.Logger.
		With().
		Str("service", serviceName).
		Logger()
	return func(c *gin.Context) {
		logger.Debug().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String())

		start := time.Now()
		defer func() {
			if time.Since(start) > time.Second*2 {
				log.Warn().
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.String()).
					Dur("elapsed_ms", time.Since(start)).
					Msg("long response time")
			}
		}()

		c.Next()
	}
}
