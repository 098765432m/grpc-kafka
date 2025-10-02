package utils

import (
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func Init() {
	// Init logger
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	if err := gotenv.Load("../../.env"); err != nil {
		zap.S().Fatal("Error loading .env file: ", err)
	}
	viper.AutomaticEnv()
}
