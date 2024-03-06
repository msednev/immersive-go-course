package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func indexByteN(s string, c byte, n int) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			count++
		}
		if count == n {
			return i
		}
	}
	return -1
}

func parseCronFile(reader io.Reader) ([]Job, error) {
	var jobs []Job
	var args []string
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		splitPos := indexByteN(line, ' ', 5)
		if splitPos == -1 {
			return nil, fmt.Errorf("failed to parse \"%s\"", line)
		}
		spec := line[:splitPos]
		cmdAndArgs := strings.Split(line[splitPos+1:], " ")
		cmd := cmdAndArgs[0]
		if len(cmdAndArgs) > 1 {
			args = cmdAndArgs[1:]
		}
		job := Job{
			Command: cmd,
			Args:    args,
			Spec:    spec,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func createTopic(p *kafka.Producer, topic string) {
	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		log.Fatalf("Failed to create new admin client from producer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		log.Fatalf("ParseDuration(60s): %v", err)
	}
	results, err := a.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     2,
			ReplicationFactor: 1,
		}},
		kafka.SetAdminOperationTimeout(maxDur),
	)
	if err != nil {
		log.Fatalf("Admin Client request error: %v", err)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Fatalf("Failed to create topic: %v", result.Error)
		}
		fmt.Printf("%v\n", result)
	}
	a.Close()
}

func main() {
	topic := "cron"
	if len(os.Args) != 2 {
		log.Fatalf("expected one command-line argument, got %d", len(os.Args)-1)
	}
	cronFileHandle, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("cannot open file: %v", err)
	}
	jobs, err := parseCronFile(cronFileHandle)
	if err != nil {
		log.Fatalf("cannot parse provided file: %v", err)
	}
	hostname, _ := os.Hostname()
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         hostname,
		"acks":              "all",
	})
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}

	createTopic(producer, topic)

	cronRunner := cron.New()
	for _, job := range jobs {
		msg, err := json.Marshal(job)
		if err != nil {
			log.Fatalf("cannot serialize cron job: %v", err)
		}
		cronRunner.AddFunc(job.Spec, func() {
			err = producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          msg,
				Key:            []byte(uuid.New().String()),
			}, nil)
		})
	}
	cronRunner.Start()

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
}
