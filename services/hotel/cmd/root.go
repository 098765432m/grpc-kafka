package cmd

import (
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/hotel/internal/app"
	"github.com/098765432m/grpc-kafka/hotel/internal/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hotel",
	Short: "Hotel service command line interface",
	Long:  `This is the command line interface for the hotel service, allowing you to manage hotel operations.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Connect to database
		conn, err := database.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.HOTEL_GRPC_PORT, conn)
		go grpcServer.Run()

		// Start Http server
		httpServer := app.NewHttpServer(consts.HOTEL_HTTP_PORT, conn)
		httpServer.Run()

	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
