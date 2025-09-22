package cmd

import (
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	image_app "github.com/098765432m/grpc-kafka/image/internal/app"
	image_database "github.com/098765432m/grpc-kafka/image/internal/database"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "image",
	Short: "Image service for hotel management system",
	Run: func(cmd *cobra.Command, args []string) {

		// Cloudinary SDK
		// cld, err := cloudinary.New()
		// if err != nil {
		// 	zap.S().Infoln("Failed to initialize Cloudinary client: ", err)
		// }

		// img, err := cld.Image()

		conn, err := image_database.Connect()
		if err != nil {
			zap.S().Fatal("Failed to connected to database")
		}

		// Start grpc server
		grpcServer := image_app.NewGrpcServer(consts.IMAGE_GRPC_PORT, conn)
		go grpcServer.Run()

		//Start Http server
		httpServer := image_app.NewHttpServer(consts.IMAGE_HTTP_PORT, conn)
		httpServer.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
