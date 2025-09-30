package rating_cmd

import (
	"strconv"

	"github.com/098765432m/grpc-kafka/common/consts"
	rating_app "github.com/098765432m/grpc-kafka/rating/internal/app"
	rating_database "github.com/098765432m/grpc-kafka/rating/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "Rating",
	Short: "Running rating service",
	Run: func(cmd *cobra.Command, args []string) {

		redisPortStr := viper.GetString("REDIS_PORT")
		if redisPortStr == "" {
			panic("There is no redis port")
		}
		redisPort, err := strconv.Atoi(redisPortStr)
		if err != nil {
			zap.S().Fatal("Invalid redis port format")
		}

		// Connect to database
		conn, err := rating_database.Connect()
		if err != nil {
			zap.S().Fatalf("Failed to connec to database: %v", err)
		}

		// Start Grpc Server
		grpcServer := rating_app.NewGrpcServer(consts.RATING_GRPC_PORT, redisPort, conn)
		go grpcServer.Run()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Fatal(err)
	}
}
