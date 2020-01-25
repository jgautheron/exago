package main

import (
	"github.com/jgautheron/exago/internal/server"
)

func main() {
	server.InitializeConfig()

	s, err := server.New()
	if err != nil {
		panic(err)
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
