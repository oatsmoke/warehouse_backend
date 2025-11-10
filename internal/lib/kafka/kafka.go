package kafka

import (
	"context"
	"encoding/json"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/segmentio/kafka-go"
)

const topic = "actions"

func SendMessage(text string) {
	broker := "localhost:9092"
	value := struct {
		Action string `json:"action"`
	}{
		Action: text,
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		writer := &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: topic,
		}
		defer writer.Close()

		marshalled, err := json.Marshal(&value)
		if err != nil {
			logger.Error("failed to marshalling message", err)
		}

		msg := kafka.Message{
			Key:   []byte("1"),
			Value: marshalled,
		}

		if err := writer.WriteMessages(ctx, msg); err != nil {
			logger.Error("failed to write message", err)
		}

		logger.Info(text)
	}()

	return
}
