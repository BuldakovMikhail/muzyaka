package muzyaka

import (
	"context"
	"github.com/go-chi/chi/v5"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	httpSwagger "github.com/swaggo/http-swagger"
	postgres2 "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
	"src/internal/config"
	delivery2 "src/internal/domain/album/delivery"
	middleware2 "src/internal/domain/album/middleware"
	postgres3 "src/internal/domain/album/repository/postgres"
	usecase2 "src/internal/domain/album/usecase"
	"src/internal/domain/auth/delivery"
	"src/internal/domain/auth/middleware"
	"src/internal/domain/auth/usecase"
	delivery3 "src/internal/domain/merch/delivery"
	middleware3 "src/internal/domain/merch/middleware"
	postgres5 "src/internal/domain/merch/repository/postgres"
	usecase4 "src/internal/domain/merch/usecase"
	delivery5 "src/internal/domain/musician/delivery"
	postgres4 "src/internal/domain/musician/repository/postgres"
	usecase3 "src/internal/domain/musician/usecase"
	delivery6 "src/internal/domain/playlist/delivery"
	middleware5 "src/internal/domain/playlist/middleware"
	postgres7 "src/internal/domain/playlist/repository/postgres"
	usecase6 "src/internal/domain/playlist/usecase"
	delivery7 "src/internal/domain/track/delivery"
	middleware7 "src/internal/domain/track/middleware"
	"src/internal/domain/track/repository/minio"
	postgres8 "src/internal/domain/track/repository/postgres"
	usecase8 "src/internal/domain/track/usecase"
	delivery8 "src/internal/domain/user/delivery"
	middleware6 "src/internal/domain/user/middleware"
	"src/internal/domain/user/repository/postgres"
	usecase7 "src/internal/domain/user/usecase"
	jwt2 "src/internal/lib/jwt"
	"src/internal/lib/logger/handlers/slogpretty"
	"time"

	_ "src/docs" // docs is generated by Swag CLI, you have to import it.
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// TODO: продублировать методы регистариции для админа, не давать возиожность вводить роль
// TODO: мб убрать musician_id из AddAlbumWithTracks
// TODO: AddAlbumWithTracks убрать id трека
// TODO: Еще не забыть про то что сорс трека должен автоматически генерироваться
// TODO: убрать к чертям везде musician_id, челу лучшее вообще про это не знать

// TODO: есть баг, что удаляя треки с помощью триггера они не удаляются в S3 (поправить в ппо)

// TODO: нет метода для связи пользователя с музыкантом
// TODO: добавить жанр мб?
// TODO: перепроверить все эндпоинты

// @title Muzyaka API
// @version 1.0
// @description API Server for musical service

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func App() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("Logger init")

	dsn := "host=localhost user=postgres password=123 dbname=postgres port=5432"
	db, err := gorm.Open(postgres2.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error())
	}

	endpoint := "localhost:9000"
	client, err := minio2.New(endpoint, &minio2.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		logger.Error(err.Error())
	}

	ctx := context.Background()

	err = client.MakeBucket(ctx, minio.TrackBucket, minio2.MakeBucketOptions{})
	if err != nil {
		logger.Error(err.Error())
	}

	tokenProvider := jwt2.NewTokenProvider("secret", time.Hour)

	userRep := postgres.NewUserRepository(db)
	albumRep := postgres3.NewAlbumRepository(db)
	trackStorage := minio.NewTrackStorage(client)
	musicianRep := postgres4.NewMusicianRepository(db)
	merchRep := postgres5.NewMerchRepository(db)
	playlistRep := postgres7.NewPlaylistRepository(db)
	trackRep := postgres8.NewTrackRepository(db)

	encryptor := usecase.NewEncryptor()
	authUseCase := usecase.NewAuthUseCase(tokenProvider, userRep, encryptor)
	musicianUseCase := usecase3.NewMusicianUseCase(musicianRep)
	albumUseCase := usecase2.NewAlbumUseCase(albumRep, trackStorage)
	merchUseCase := usecase4.NewMerchUseCase(merchRep)
	playlistUseCase := usecase6.NewPlaylistUseCase(playlistRep, trackRep)
	userUseCase := usecase7.NewUserUseCase(userRep, trackRep, encryptor)
	trackUseCase := usecase8.NewTrackUseCase(trackRep, trackStorage)

	musicianMiddleware := (func(h http.Handler) http.Handler {
		return middleware.CheckMusicianLevelPermissions(h, authUseCase)
	})

	adminMiddleware := (func(h http.Handler) http.Handler {
		return middleware.CheckAdminLevelPermissions(h, authUseCase)
	})

	userMiddleware := (func(h http.Handler) http.Handler {
		return middleware.CheckUserLevelPermissions(h, authUseCase)
	})

	checkForMusicianId := (func(h http.Handler) http.Handler {
		return middleware2.CheckIsUserRelatedToMusician(h, musicianUseCase)
	})

	checkForUserId := (func(h http.Handler) http.Handler {
		return middleware6.CheckIsUserRelated(h, userUseCase)
	})

	checkIsAlbumRelated := (func(h http.Handler) http.Handler {
		return middleware2.CheckAlbumOwnership(h, albumUseCase, musicianUseCase)
	})

	checkIsPlaylistRelated := (func(h http.Handler) http.Handler {
		return middleware5.CheckPlaylistOwnership(h, playlistUseCase)
	})

	checkIsMerchRelated := (func(h http.Handler) http.Handler {
		return middleware3.CheckMerchOwnership(h, merchUseCase, musicianUseCase)
	})

	checkIsTrackRelated := (func(h http.Handler) http.Handler {
		return middleware7.CheckTrackOwnership(h, albumUseCase, musicianUseCase)
	})

	basicAuthMiddleware := (func(h http.Handler) http.Handler {
		return middleware.JwtParseMiddleware(h, authUseCase)
	})

	// TODO: Handlers
	router := chi.NewRouter()

	//auth
	router.Post("/api/auth/sign-up/user", delivery.SignUp(authUseCase))
	router.Post("/api/auth/sign-up/admin", delivery.SignUpAdmin(authUseCase))
	router.Post("/api/auth/sign-in", delivery.SignIn(authUseCase))
	router.Post("/api/auth/sign-up/musician", delivery.SignUpMusician(authUseCase))
	router.Get("/api/get-me", delivery8.GetMe(musicianUseCase))

	// album
	router.Group(func(r chi.Router) {
		r.Use(musicianMiddleware)
		r.With(checkForMusicianId).Post("/api/musician/{musician_id}/album", delivery2.AddAlbumWithTracks(albumUseCase))

		r.Group(func(r chi.Router) {
			r.Use(checkIsAlbumRelated)
			r.Post("/api/album/{id}/tracks", delivery2.CreateTrack(albumUseCase))
			r.Delete("/api/album/{id}", delivery2.DeleteAlbum(albumUseCase))
			r.Put("/api/album/{id}", delivery2.UpdateAlbum(albumUseCase))
		})
	})

	// merch
	router.Group(func(r chi.Router) {
		r.Use(musicianMiddleware)
		r.With(checkForMusicianId).Post("/api/musician/{musician_id}/merch", delivery3.MerchCreate(merchUseCase))

		r.Group(func(r chi.Router) {
			r.Use(checkIsMerchRelated)
			r.Delete("/api/merch/{id}", delivery3.DeleteMerch(merchUseCase))
			r.Put("/api/merch/{id}", delivery3.UpdateMerch(merchUseCase))
		})
	})

	// musician
	router.Group(func(r chi.Router) {
		r.Use(musicianMiddleware)
		r.With(checkForMusicianId).Put("/api/musician/{musician_id}", delivery5.UpdateMusician(musicianUseCase))
		r.With(checkForMusicianId).Delete("/api/musician/{musician_id}", delivery5.DeleteMusician(musicianUseCase))
	})

	// playlist

	router.Group(func(r chi.Router) {
		r.Use(userMiddleware)

		r.With(checkForUserId).Post("/api/user/{user_id}/playlist", delivery6.PlaylistCreate(playlistUseCase))
		r.Group(func(r chi.Router) {
			r.Use(checkIsPlaylistRelated)
			r.Put("/api/playlist/{id}", delivery6.UpdatePlaylist(playlistUseCase))
			r.Delete("/api/playlist/{id}", delivery6.DeletePlaylist(playlistUseCase))
			r.Post("/api/playlist/{id}/track", delivery6.AddTrack(playlistUseCase))
			r.Delete("/api/playlist/{id}/track/{track_id}", delivery6.DeleteTrack(playlistUseCase))
		})
	})

	// Track
	router.Group(func(r chi.Router) {
		r.Use(musicianMiddleware)
		r.Use(checkIsTrackRelated)
		r.Put("/api/track/{id}", delivery7.UpdateTrack(trackUseCase))
		r.Delete("/api/track/{id}", delivery7.DeleteTrack(trackUseCase))
	})

	// User
	router.Group(func(r chi.Router) {
		r.Use(basicAuthMiddleware)
		r.Use(checkForUserId)
		r.Get("/api/user/{user_id}", delivery8.GetUser(userUseCase))
		r.Put("/api/user/{user_id}", delivery8.UpdateUser(userUseCase))
		r.Delete("/api/user/{user_id}", delivery8.DeleteUser(userUseCase))
	})

	// admin only
	router.Group(func(r chi.Router) {
		r.Use(adminMiddleware)
		r.Post("/api/musician", delivery5.CreateMusician(musicianUseCase))
	})

	// Likes
	router.Group(func(r chi.Router) {
		r.Use(userMiddleware)
		r.Use(checkForUserId)
		r.Post("/api/user/{user_id}/favorite", delivery8.Like(userUseCase))
		r.Delete("/api/user/{user_id}/favorite", delivery8.Dislike(userUseCase))
		r.Get("/api/user/{user_id}/favorite", delivery8.GetAllLiked(userUseCase))
	})

	// Other opened requests
	router.Group(func(r chi.Router) {
		r.Use(basicAuthMiddleware)
		r.Get("/api/track/{id}", delivery7.GetTrack(trackUseCase))
		r.Get("/api/playlist/{playlist_id}/track", delivery6.GetAllTracks(playlistUseCase))
		r.Get("/api/playlist/{id}", delivery6.GetPlaylist(playlistUseCase))
		r.Get("/api/musician/{musician_id}", delivery5.GetMusician(musicianUseCase))
		r.Get("/api/album/{id}/tracks", delivery2.GetAllTracks(albumUseCase))
		r.Get("/api/musician/{musician_id}/merch", delivery3.GetAllMerchForMusician(merchUseCase))
		r.Get("/api/album/{id}", delivery2.GetAlbum(albumUseCase))
		r.Get("/api/merch/{id}", delivery3.GetMerch(merchUseCase))
		r.Get("/api/musician/{musician_id}/album", delivery2.GetAllAlbumForMusician(albumUseCase))
	})

	// Swagger
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":8080", router)
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