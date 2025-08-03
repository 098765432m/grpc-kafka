package cmd

import (
	"github.com/098765432m/grpc-kafka/common/consts"
	rating_app "github.com/098765432m/grpc-kafka/rating/internal/app"
	rating_database "github.com/098765432m/grpc-kafka/rating/internal/database"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "Rating",
	Short: "Running rating service",
	Run: func(cmd *cobra.Command, args []string) {

		// Connect to database
		conn, err := rating_database.Connect()
		if err != nil {
			zap.S().Fatalf("Failed to connec to database: %v", err)
		}

		// Start Grpc Server
		grpcServer := rating_app.NewGrpcServer(consts.RATING_GRPC_PORT, conn)
		go grpcServer.Run()

		// Start Http server
		httpServer := rating_app.NewHttpServer(consts.RATING_HTTP_PORT, conn)
		httpServer.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Fatal(err)
	}
}
