package tools

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"sync"
)

var logger zap.Logger
var once = new(sync.Once)

func Logger() *zap.Logger {
	once.Do(func() {
		newLogger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = *newLogger
	})
	return &logger
}

func LoggingHttpResponse(response *http.Response, err error) {

	Logger().Info("response",
		zap.String("url", response.Request.URL.String()),
		zap.Int("status", response.StatusCode),
		zap.Any("header", response.Header),
		zap.String("rawResponse", func() string {
			bytes, _ := httputil.DumpResponse(response, true)
			return string(bytes)
		}()), zap.Error(err))
}
