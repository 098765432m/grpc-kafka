package main

import (
	"github.com/098765432m/grpc-kafka/booking/cmd"
	"github.com/098765432m/grpc-kafka/common/utils"
)

func main() {
	utils.Init()

	cmd.Execute()
}
