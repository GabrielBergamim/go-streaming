package producer_impl

import (
	"context"
	"encoding/json"
	"log"

	"github.com/example/go-streaming/processor/domain/message"
	"github.com/example/go-streaming/processor/domain/movie"
	"github.com/segmentio/kafka-go"
)

type InfoProducerImpl struct {
	kafkaProducer *kafka.Writer
}

func NewInfoProducerImpl(kafkaProducer *kafka.Writer) *InfoProducerImpl {
	return &InfoProducerImpl{
		kafkaProducer: kafkaProducer,
	}
}

func (p *InfoProducerImpl) Send(mv *movie.Movie) error {
    defer p.kafkaProducer.Close()

	var event = message.MovieInfoEvent {
		ID: mv.ID.String(),
		Name: mv.Title,
	}

	jsonValue, err := json.Marshal(&event)

	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte("key"),
		Value: jsonValue,
	}

	return p.kafkaProducer.WriteMessages(context.Background(), message)
}

