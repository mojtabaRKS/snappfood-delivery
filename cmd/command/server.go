package command

import (
	"context"
	"fmt"
	"log"
	"snappfood/internal/api/rest"
	order_handler "snappfood/internal/api/rest/handlers/order"
	vendor_handler "snappfood/internal/api/rest/handlers/vendor"
	"snappfood/internal/config"
	"snappfood/internal/infra/mysql"
	"snappfood/internal/infra/redis"
	"snappfood/internal/repositories"
	order_service "snappfood/internal/services/order"
	"sync"

	"github.com/spf13/cobra"
)

type Server struct{}

func (cmd Server) Command(ctx context.Context, cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "run server",
		Run: func(_ *cobra.Command, _ []string) {
			cmd.main(cfg, ctx)
		},
	}
}

func (cmd *Server) main(cfg *config.Config, ctx context.Context) {
	var wg sync.WaitGroup

	db, err := mysql.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to mysql database %s", err.Error())
		return
	}
	gormDB, err := mysql.NewGormWithInstance(db, cfg.AppDebug)
	if err != nil {
		log.Fatalf("failed to connect to mysql database: err : %s", err.Error())
		return
	}

	redisClient, err := redis.NewClient(ctx, *cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis, %s", err.Error())
	}

	err = mysql.Migrate(db)
	if err != nil {
		log.Fatalf("mysql migration failed: %s", err.Error())
	}

	server := rest.New()

	// intialize repositories
	orderRepo := repositories.NewOrderRepo(gormDB)
	delayReportRepo := repositories.NewDelayRepo(gormDB)
	tripRepo := repositories.NewTripRepo(gormDB)

	// intialize services
	orderService := order_service.New(orderRepo, delayReportRepo, tripRepo, redisClient)

	// intialize handlers
	orderHandler := order_handler.New(orderService)
	vendorHandler := vendor_handler.New(orderService)

	server.SetupAPIRoutes(orderHandler, vendorHandler)
	if err := server.Serve(ctx, fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
