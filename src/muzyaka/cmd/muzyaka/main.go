package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"src/internal/config"
	"src/internal/cron/outbox_producer/usecase"
	"src/internal/lib/kafka"
	"src/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("Logger init")

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=32771"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error())
	}
	//
	pr, err := kafka.NewProducer("localhost:29092")
	op := usecase.New(pr, db, "outbox")
	op.ProduceMessages()
	//
	//rep := postgres.NewAlbumRepository(db)
	//err = rep.DeleteAlbumOutbox(1)
	//_, err = rep.AddAlbum(&models.Album{
	//	Id:    0,
	//	Name:  "Aboba2",
	//	Cover: "aboba2.org",
	//	Type:  "LP",
	//})
	//
	//if err != nil {
	//	logger.Error(err.Error())
	//}

	//_, err = rep.AddTrackToAlbumOutbox(1, &models.Track{
	//	Id:         0,
	//	Source:     "test",
	//	Producers:  nil,
	//	Authors:    nil,
	//	Performers: nil,
	//	Name:       "test",
	//	Genre:      "test",
	//})
	//
	//if err != nil {
	//	logger.Error(err.Error())
	//}

	// TODO: Handlers

	// TODO: server run
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
