package main

import (
	booking_cmd "github.com/098765432m/grpc-kafka/booking/cmd"
	"github.com/098765432m/grpc-kafka/common/utils"
)

func main() {
	utils.Init()

	booking_cmd.Execute()
}
