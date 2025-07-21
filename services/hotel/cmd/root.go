package cmd

import (
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/hotel/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hotel",
	Short: "Hotel service command line interface",
	Long:  `This is the command line interface for the hotel service, allowing you to manage hotel operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.HOTEL_GRPC_PORT)
		go grpcServer.Run()

		// Start Http server
		httpServer := app.NewHttpServer(consts.HOTEL_HTTP_PORT)
		httpServer.Run()
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
