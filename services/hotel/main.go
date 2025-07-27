package main

import (
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/098765432m/grpc-kafka/hotel/cmd"
)

func main() {

	utils.Init()

	cmd.Execute()
}
