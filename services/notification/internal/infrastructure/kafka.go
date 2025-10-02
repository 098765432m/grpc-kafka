package notification_infrastructure

import (
	"context"
	"encoding/json"

	notification_domain "github.com/098765432m/grpc-kafka/notification/internal/domain"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	cl *kgo.Client
}

func NewKafkaConsumer(brokers []string, topics, groupId string) (*KafkaConsumer, error) {

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topics),
		kgo.ConsumerGroup(groupId),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		cl: client,
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler func(notification_domain.BookingCreatedEvent)) {

	for {
		fetches := kc.cl.PollFetches(ctx)
		if fetches.IsClientClosed() {
			return
		}

		fetches.EachError(func(t string, p int32, err error) {
			zap.S().Infoln("fetch err topic %s partition %d: %v", t, p, err)
		})

		fetches.EachRecord(func(r *kgo.Record) {
			var event notification_domain.BookingCreatedEvent
			if err := json.Unmarshal(r.Value, &event); err != nil {
				zap.S().Errorf("Failed to unmarshal record: %v", err)
				return
			}

			// pass to handler
			handler(event)

			// commit offset
			if err := kc.cl.CommitRecords(ctx, r); err != nil {
				zap.S().Infoln("Failed to commit record: %v", err)
			}
		})
	}

}
