package logging

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	model "github.com/prebid/prebid-server/logging/model"
	kafka "github.com/segmentio/kafka-go"
)

const (
	// The magic byte defined by Avro
	magicByte byte = 0x0

	mspServerLogSchema   int = 1
	mspRequestSchema     int = 2
	mspResponseSchema    int = 3
	bidderRequestSchema  int = 4
	bidderResponseSchema int = 5

	topic string = "msp_server_log"
)

var kafkaWriter = &kafka.Writer{}

func Connect() {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
}

func Close() {
	if err := kafkaWriter.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func MSPServerLog(mspServerLog model.Msp_server_log) {
	var buf bytes.Buffer
	mspServerLog.Serialize(&buf)
	msg, _ := addSchema(mspServerLogSchema, buf.Bytes())
	produce(msg)
}

func addSchema(id int, msgBytes []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := buf.WriteByte(magicByte)
	if err != nil {
		return nil, err
	}
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, uint32(id))
	_, err = buf.Write(idBytes)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(msgBytes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func produce(msg []byte) {
	if err := kafkaWriter.WriteMessages(context.Background(), kafka.Message{Value: msg}); err != nil {
		// Only validation errors would be reported in this case.
		log.Fatal("failed to write messages:", err)
	}
}
