package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/SousaGLucas/apps/banking/account/app/service/api"
	"github.com/SousaGLucas/apps/banking/account/app/service/api/v1"
	"github.com/SousaGLucas/apps/banking/account/domain/usecases/account"
	"github.com/SousaGLucas/apps/banking/account/gateways/auth"
	"github.com/SousaGLucas/apps/banking/account/gateways/ledger"
	"github.com/SousaGLucas/apps/banking/account/gateways/pg"
	"github.com/SousaGLucas/apps/banking/account/gateways/pg/migrations"
)

type config struct {
	ServerAddress string `conf:"env:SERVER_ADDRESS,default:0.0.0.0:3000"`

	DatabaseHost     string `conf:"env:DATABASE_HOST,default:0.0.0.0"`
	DatabasePort     string `conf:"env:DATABASE_PORT,default:5432"`
	DatabaseDatabase string `conf:"env:DATABASE_DATABASE,default:postgres"`
	DatabaseUser     string `conf:"env:DATABASE_USER,default:postgres"`
	DatabasePassword string `conf:"env:DATABASE_PASSWORD,required"`

	BankingAuthURL   string `conf:"env:BANKING_AUTH_URL,required"`
	BankingLedgerURL string `conf:"env:BANKING_LEDGER_URL,required"`
}

func main() {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("starting logger: %s", err))
	}
	defer func() {
		_ = logger.Sync()
	}()

	mainLogger := logger.With(
		zap.Int("go_max_procs", runtime.GOMAXPROCS(0)),
		zap.Int("runtime_num_cpu", runtime.NumCPU()),
	)

	var cfg config
	if help, err := conf.Parse("", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		mainLogger.Error("loading configuration file", zap.Error(err))
		return
	}

	mainLogger.Info("starting service")
	if err := Main(logger, cfg); err != nil {
		mainLogger.Error("failed to run service", zap.Error(err))
	}
}

func Main(logger *zap.Logger, cfg config) error {
	pgPool, err := newPGPool(cfg)
	if err != nil {
		return fmt.Errorf("creating postgres connection: %w", err)
	}

	authClient, err := auth.NewClient(auth.Config{BaseURL: cfg.BankingAuthURL})
	if err != nil {
		return fmt.Errorf("creating auth client: %w", err)
	}

	ledgerClient, err := ledger.NewClient(ledger.Config{BaseURL: cfg.BankingLedgerURL})
	if err != nil {
		return fmt.Errorf("creating ledger client: %w", err)
	}

	createAccountUseCase := account.CreateAccountUseCase{
		DB:     pg.AccountsRepository{DB: pgPool},
		Auth:   authClient,
		Ledger: ledgerClient,
	}
	getAccountUseCase := account.GetAccountUseCase{
		DB: pg.AccountsRepository{DB: pgPool},
	}
	getBalanceUseCase := account.GetBalanceUseCase{
		DB:     pg.AccountsRepository{DB: pgPool},
		Ledger: ledgerClient,
	}
	depositUseCase := account.DepositUseCase{
		DB:     pg.AccountsRepository{DB: pgPool},
		Ledger: ledgerClient,
	}
	withdrawUseCase := account.WithdrawUseCase{
		DB:     pg.AccountsRepository{DB: pgPool},
		Ledger: ledgerClient,
	}

	router := api.NewServer(logger)

	apiV1 := v1.API{
		CreateAccountHandler: v1.CreateAccountHandler(createAccountUseCase),
		GetAccountHandler:    v1.GetAccountHandler(getAccountUseCase),
		GetBalanceHandler:    v1.GetBalanceHandler(getBalanceUseCase),
		DepositHandler:       v1.DepositHandler(depositUseCase),
		WithdrawHandler:      v1.WithdrawHandler(withdrawUseCase),
	}
	apiV1.Routes(router)

	srv := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errGroup, _ := errgroup.WithContext(context.Background())

	errGroup.Go(func() error {
		logger.Info(fmt.Sprintf("starting http server: %s", cfg.ServerAddress))
		_ = srv.ListenAndServe()
		return nil
	})

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutting down http server: %w", err)
	}

	logger.Info("http server stopped")
	return nil
}

func newPGPool(cfg config) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(pgConnString(cfg))
	if err != nil {
		return nil, fmt.Errorf("parsing postgres config: %w", err)
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("creating postgres connection: %w", err)
	}

	err = runMigrations(cfg)
	if err != nil {
		return nil, fmt.Errorf("run postgres migrations: %w", err)
	}

	return pgPool, nil
}

func runMigrations(cfg config) error {
	pgxCfg, err := pgx.ParseConfig(pgConnString(cfg))
	if err != nil {
		return fmt.Errorf("parsing pgx config: %w", err)
	}

	driver, err := postgres.WithInstance(stdlib.OpenDB(*pgxCfg), &postgres.Config{
		DatabaseName: cfg.DatabaseDatabase,
	})
	if err != nil {
		return fmt.Errorf("creating postgres driver: %w", err)
	}

	source, err := httpfs.New(http.FS(migrations.Migrations), ".")
	if err != nil {
		return fmt.Errorf("creating migrations source: %w", err)
	}

	instance, err := migrate.NewWithInstance("httpfs", source, cfg.DatabaseDatabase, driver)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("uping migrations: %w", err)
	}

	sourceErr, databaseErr := instance.Close()
	if sourceErr != nil {
		return fmt.Errorf("closing db source: %w", err)
	}
	if databaseErr != nil {
		return fmt.Errorf("closing migrations connection: %w", err)
	}

	return nil
}

func pgConnString(cfg config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseDatabase,
	)
}
