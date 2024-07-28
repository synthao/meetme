package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	pb "github.com/synthao/meetme/gen/go/imgproc"
	"github.com/synthao/meetme/internal/client/grpc"
	"github.com/synthao/meetme/internal/config"
	userInfra "github.com/synthao/meetme/internal/infrastructure"
	userApp "github.com/synthao/meetme/internal/user/application"
	userInt "github.com/synthao/meetme/internal/user/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cnf := config.LoadConfig()

	db := sqlx.MustConnect("postgres", config.GetDSN())
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to db", err.Error())
	}

	log.Println("connected to db successfully")

	logger, err := newLogger(cnf.LogLevel)
	if err != nil {
		return
	}

	applyMigrations(logger, db.DB)

	logger.Info("connecting to gRPC server", zap.String("address", os.Getenv("IMGPROC_GRPC_SERVER_ADDRESS")))

	grpcClient := grpc.NewClient(cnf.GrpcServerAddress, logger)

	grpcClient.Connect()
	defer grpcClient.Close()

	var (
		grpcImgClient = pb.NewImageProcessingServiceClient(grpcClient)
		mainRouter    = mux.NewRouter()
	)

	var (
		userRepo    = userInfra.NewRepository(db)
		userService = userApp.NewService(userRepo, grpcImgClient)
		userHandler = userInt.NewHandler(userService)
	)

	mainRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	mainRouter.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"APP_VERSION": cnf.AppVersion,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(data)
	}).Methods("GET")

	userInt.InitRoutes(mainRouter, userHandler)

	srv := &http.Server{
		Addr:    cnf.ServerAddress,
		Handler: mainRouter,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)

	g.Go(func() error {
		log.Println("Starting server on", cnf.ServerAddress)
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Fatal("exit reason: ", err)
	}
}

func newLogger(lvl string) (*zap.Logger, error) {
	atomicLogLevel, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		return nil, err
	}

	atom := zap.NewAtomicLevelAt(atomicLogLevel.Level())
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
	), nil
}

func applyMigrations(logger *zap.Logger, db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("cannot create postgres driver", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // migration files location
		"postgres",          // database name
		driver,
	)
	if err != nil {
		logger.Fatal("could not create the migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal("could not apply migrations", zap.Error(err))
	}

	logger.Info("migrations applied successfully")
}
