package main

import (
	"os"

	"github.com/eriktate/divulge/http"
	"github.com/eriktate/divulge/mock"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting...")
	host := getEnvString("DIVULGE_HOST", "localhost")
	ps := &mock.PostService{}
	server := http.New(host, 8080, http.BuildRouter(ps, logger))
	logger.Error(server.Start())
}

func getEnvString(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
