package main

import (
	"github.com/098765432m/grpc-kafka/common/utils"
	rating_cmd "github.com/098765432m/grpc-kafka/rating/cmd"
)

func main() {
	utils.Init()

	rating_cmd.Execute()
}
