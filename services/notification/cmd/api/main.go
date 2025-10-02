package main

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))

	seeds := []string{"localhost:????"}

	// 1. Config
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		panic(err)
	}

	zap.S().Infoln("Kafka client connected")

	adminClient := kadm.NewClient(kafkaClient)

	topics := []string{
		"hotel.booking.created",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = adminClient.CreateTopics(ctx, -1, -1, nil, topics...)
	if err != nil {
		panic(err)
	}
}
