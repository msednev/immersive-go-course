package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type CronSpec struct {
	Minute uint64 `json:"minute,omitempty"`
	Hour   uint64 `json:"hour,omitempty"`
	Dom    uint64 `json:"dom,omitempty"`
	Month  uint64 `json:"month,omitempty"`
	Dow    uint64 `json:"dow,omitempty"`
}

type Job struct {
	Command string `json:"command,omitempty"`
	CronSpec
}

func parseCronFile(reader io.Reader) ([]Job, error) {
	var jobs []Job
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lastInd := strings.LastIndex(line, " ")
		if lastInd == -1 {
			return nil, fmt.Errorf("failed to parse \"%s\"", line)
		}
		cronExpr := line[:lastInd]
		command := line[lastInd+1:]
		sched, err := parser.Parse(cronExpr)
		if err != nil {
			return nil, err
		}
		specShed, ok := sched.(*cron.SpecSchedule)
		if ok == false {
			return nil, fmt.Errorf("type assertion failed")
		}
		job := Job{
			Command: command,
			CronSpec: CronSpec{
				specShed.Minute,
				specShed.Hour,
				specShed.Dom,
				specShed.Month,
				specShed.Dow,
			},
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
			NumPartitions:     1,
			ReplicationFactor: 3,
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
	msg, err := json.Marshal(jobs)
	if err != nil {
		log.Fatalf("cannot serialize cron jobs: %v", err)
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
	deliveryChan := make(chan kafka.Event, 100)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, deliveryChan)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("successfully produced record to topic %s partition [%d] @ offset %v\n",
						ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
}
