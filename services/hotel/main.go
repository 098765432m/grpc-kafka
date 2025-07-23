package main

import (
	"github.com/098765432m/grpc-kafka/hotel/cmd"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func main() {
	// Init logger
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	if err := gotenv.Load(".env"); err != nil {
		zap.S().Fatal("Error loading .env file: ", err)
	}
	viper.AutomaticEnv()

	cmd.Execute()
}
