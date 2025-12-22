package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/segmentio/kafka-go"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type KafkaProducer struct {
    writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) interfaces.MessageProducer {
    writer := &kafka.Writer{
        Addr:                   kafka.TCP(brokers...),
        Balancer:               &kafka.LeastBytes{},
        AllowAutoTopicCreation: true,
        Async:                  true, // Асинхронная отправка
    }
    
    return &KafkaProducer{
        writer: writer,
    }
}

func (p *KafkaProducer) Send(ctx context.Context, topic string, key, value []byte) error {
    message := kafka.Message{
        Topic: topic,
        Key:   key,
        Value: value,
    }
    
    return p.writer.WriteMessages(ctx, message)
}

func (p *KafkaProducer) SendJSON(ctx context.Context, topic string, data interface{}) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal json: %w", err)
    }
    
    key := []byte(fmt.Sprintf("%v", data))
    return p.Send(ctx, topic, key, jsonData)
}

func (p *KafkaProducer) Close() error {
    return p.writer.Close()
}