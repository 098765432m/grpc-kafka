package booking_cmd

import (
	"log"

	booking_app "github.com/098765432m/grpc-kafka/booking/internal/app"
	booking_database "github.com/098765432m/grpc-kafka/booking/internal/database"
	"github.com/098765432m/grpc-kafka/common/consts"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "booking",
	Short: "Booking service command line interface",
	Run: func(cmd *cobra.Command, args []string) {

		// Connect to database
		conn, err := booking_database.Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v\n", err)
		}

		grpcServer := booking_app.NewBookingGrpcServer(consts.BOOKING_GRPC_PORT, conn)
		grpcServer.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.S().Fatal(err)
	}
}
