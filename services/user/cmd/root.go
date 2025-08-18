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

		conn, err := user_database.Connect()
		if err != nil {
			zap.S().Fatal("Failed to connect to database: ", err.Error())
		}

		// Start grpc server
		grpcServer := app.NewGrpcServer(consts.USER_GRPC_PORT, conn)
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
