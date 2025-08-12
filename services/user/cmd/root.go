package cmd

import (
	"fmt"
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/user/internal/app"
	user_database "github.com/098765432m/grpc-kafka/user/internal/database"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "User service command line interface",
	Run: func(cmd *cobra.Command, args []string) {
		// Dial grpc
		// hotelConn := NewGrpcClient(consts.HOTEL_GRPC_PORT)
		// defer hotelConn.Close()

		// hotelClient := hotel_pb.NewHotelServiceClient(hotelConn)

		// _, err := hotelClient.CreateHotel(context.Background(), &hotel_pb.CreateHotelRequest{Name: "Test Hotel"})
		// if err != nil {
		// 	log.Fatalf("Failed to create hotel: %v", err)
		// }

		// fmt.Println("Hotel created successfully")

		// hotels, err := hotelClient.GetAllHotels(context.Background(), &hotel_pb.GetAllHotelsRequest{})
		// if err != nil {
		// 	log.Fatalf("Failed to get all hotels: %v", err)
		// }

		// for _, hotel := range hotels.Hotels {
		// 	fmt.Printf("Hotel ID: %s, Name: %s\n", hotel.Id, hotel.Name)
		// }

		conn, err := user_database.Connect()
		if err != nil {
			zap.S().Fatal("Failed to connect to database: ", err.Error())
		}

		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.USER_GRPC_PORT)
		go grpcServer.Run()

		// Start HTTP server
		httpServer := app.NewHttpServer(consts.USER_HTTP_PORT, conn)
		httpServer.Run()
	},
}

func NewGrpcClient(addr int) *grpc.ClientConn {
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to grpc server on port %d: %v", addr, err)
	}

	return conn
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
