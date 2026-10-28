package main

import (
	"gemini-api/internal/app/server"
	"log/slog"
)

func main() {
	slog.Info("Start gemini-api Service ....")
	err := server.StartServer()
	if err != nil {
		slog.Error("Start gemini-api Service Failed", err)
	}
}
