package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	. "kafka-cron/utils"
	"log"
	"os/exec"
	"time"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka1:19092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	defer c.Close()
	if err != nil {
		log.Fatalf("cannot create a new consumer: %v", err)
	}
	if err := c.Subscribe("cron", nil); err != nil {
		log.Fatalf("cannot subscribe to topic: %v", err)
	}

	for {
		var job Job
		msg, err := c.ReadMessage(30 * time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		} else {
			fmt.Println(err)
			continue
		}
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			fmt.Printf("Cannot deserialize message: %v\n", msg.Value)
		}
		cmd := exec.Command("sh", "-c", job.Command)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Cannot execute command %v: %v\n", cmd, err)
		}
	}
}
