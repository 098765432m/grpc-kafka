package main

import (
	"github.com/098765432m/grpc-kafka/common/utils"
	"github.com/098765432m/grpc-kafka/user/cmd"
)

func main() {
	utils.Init()

	cmd.Execute()
}
