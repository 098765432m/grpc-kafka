package cmd

import (
	"fmt"
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/image_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/098765432m/grpc-kafka/hotel/internal/app"
	"github.com/098765432m/grpc-kafka/hotel/internal/database"
	hotel_handler "github.com/098765432m/grpc-kafka/hotel/internal/handler/hotel"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/room-type"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/service/room-type"
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

		// repo
		hotelRepo := hotel_repo.New(conn)

		// dial grpc
		imageConn := utils.NewGrpcClient(fmt.Sprintf(":%d", consts.IMAGE_GRPC_PORT))
		imageClient := image_pb.NewImageServiceClient(imageConn)

		hotelService := hotel_service.NewHotelService(hotelRepo, imageClient)

		roomTypeRepo := room_type_repo.New(conn)
		roomTypeService := room_type_service.NewRoomTypeService(roomTypeRepo)

		roomRepo := room_repo.New(conn)
		roomService := room_service.NewRoomService(roomRepo)

		hotelHandler := hotel_handler.NewHotelHttpHandler(hotelService, roomTypeService, roomService)

		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.HOTEL_GRPC_PORT, hotelService)
		go grpcServer.Run()

		// Start Http server
		httpServer := app.NewHttpServer(consts.HOTEL_HTTP_PORT, hotelHandler)
		httpServer.Run()

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
