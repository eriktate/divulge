package main

import (
	"github.com/eriktate/divulge/http"
	"github.com/eriktate/divulge/mock"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting...")
	ps := &mock.PostService{}
	server := http.New("localhost", 1337, http.BuildRouter(ps, logger))
	logger.Error(server.Start())
}
