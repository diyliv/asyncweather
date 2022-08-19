package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func (s *server) StartHTTP() {
	s.logger.Info(fmt.Sprintf("Starting HTTP Server on port %s\n", s.cfg.HttpServer.Port))

	http.Handle("/register", s.handlers.Register())
	http.Handle("/login", s.handlers.Login())
	http.Handle("/api/forecast", s.handlers.GetForecast())

	go func() {
		if err := http.ListenAndServe(s.cfg.HttpServer.Port, nil); err != nil {
			panic(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	<-done
}
