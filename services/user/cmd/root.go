package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotels"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "user",
	Short: "User service command line interface",
	Run: func(cmd *cobra.Command, args []string) {
		// Dial grpc
		hotelConn := NewGrpcClient(consts.HOTEL_GRPC_PORT)
		defer hotelConn.Close()

		hotelClient := hotels.NewHotelServiceClient(hotelConn)

		fmt.Println("123123")
		hotel, err := hotelClient.GetHotel(
			context.Background(),
			&hotels.GetHotelRequest{Id: "czxczxczxc"},
		)

		fmt.Println(hotel, err)

		// // Start grpc server
		// grpcServer := app.NewGrpcServer(consts.USER_GRPC_PORT)
		// go grpcServer.Run()

		// // Start HTTP server
		// httpServer := app.NewHttpServer(consts.USER_HTTP_PORT)
		// httpServer.Run()
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
