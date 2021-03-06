package main

import (
	"fmt"
	"github.com/sudiyi/sdy/aliware/kafka"
	"github.com/tidwall/gjson"
	"os"
	"os/signal"
	"syscall"
)

var demoConfigConsumerString = `
{
  "topics": ["demo"],
  "servers": ["kafka1:9092"],
  "consumerId": "demo-consumer-group",
}
`

func main() {
	results := gjson.GetMany(demoConfigConsumerString, "servers", "consumerId", "topics")
	servers := results[0].Array()
	var s []string
	s = append(s, servers[0].String())
	consumerId, topics := results[1].String(), results[2].Array()
	var t []string
	t = append(t, topics[0].String())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)

	client := kafka.New(s, true)
	consumer, err := client.NewConsumer(consumerId, t, `oldest`)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			fmt.Printf(
				"Topic: %s, Key: %s, Partition: %d, Offset: %d, Content: %s \n", msg.Topic(),
				msg.Key(), msg.Partition(), msg.Offset(), string(msg.Value()),
			)
			consumer.Commit(msg)
		case <-signals:
			fmt.Println("Stop consumer server...")
			consumer.Close()
			return
		}
	}
}
