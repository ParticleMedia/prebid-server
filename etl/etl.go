package etl

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/prebid/prebid-server/config"
	model "github.com/prebid/prebid-server/etl/model"
	kafka "github.com/segmentio/kafka-go"
)

// The magic byte defined by Avro
const magicByte byte = 0x0

type DataProducer interface {
	ProduceOpenrtb2Auction(moa model.Msp_openrtb2_auction)
}

type KafkaDataProducer struct {
	auctionAvroSchema int
	kafkaWriter       *kafka.Writer
}

func NewKafkaDataProducer(cfg config.ETL) *KafkaDataProducer {
	return &KafkaDataProducer{
		auctionAvroSchema: cfg.AvroSchemaAuction,
		kafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaHost),
			Topic:    cfg.KafkaTopic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		},
	}
}

func (k *KafkaDataProducer) ProduceOpenrtb2Auction(moa model.Msp_openrtb2_auction) {
	// print(moa)

	var buf bytes.Buffer
	moa.Serialize(&buf)
	msg, _ := addSchema(k.auctionAvroSchema, buf.Bytes())

	if err := k.kafkaWriter.WriteMessages(context.Background(), kafka.Message{Value: msg}); err != nil {
		// Only validation errors would be reported in this case.
		log.Fatal("failed to write messages:", err)
	}
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

func print(moa model.Msp_openrtb2_auction) {
	out, _ := json.Marshal(moa)
	fmt.Println("** MSP SERVER LOG **", string(out))
}
