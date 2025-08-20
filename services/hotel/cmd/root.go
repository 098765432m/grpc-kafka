package cmd

import (
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	hotel_app "github.com/098765432m/grpc-kafka/hotel/internal/app/hotel"
	room_app "github.com/098765432m/grpc-kafka/hotel/internal/app/room"
	room_type_app "github.com/098765432m/grpc-kafka/hotel/internal/app/room-type"
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
		hotelGrpcServer := hotel_app.NewGrpcServer(consts.HOTEL_GRPC_PORT, conn)
		go hotelGrpcServer.Run()

		roomTypeGrpcServer := room_type_app.NewGrpcServer(consts.ROOM_TYPE_GRPC_PORT, conn)
		go roomTypeGrpcServer.Run()

		roomGrpcServer := room_app.NewGrpcServer(consts.ROOM_GRPC_PORT, conn)
		roomGrpcServer.Run()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
