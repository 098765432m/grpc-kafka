package cmd

import (
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/image/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "image",
	Short: "Image service for hotel management system",
	Run: func(cmd *cobra.Command, args []string) {

		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.IMAGE_GRPC_PORT, nil)
		go grpcServer.Run()

		//Start Http server
		httpServer := app.NewHttpServer(consts.IMAGE_HTTP_PORT, nil)
		httpServer.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
