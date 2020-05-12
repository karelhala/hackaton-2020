package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/joho/godotenv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"net/http"
)

type Topics struct {
	Topic string `json:"topic"`
}

func readTopics() []Topics {
	file, _ := ioutil.ReadFile("topics.json")
	data := make([]Topics,0)

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func sendToListener(msg *kafka.Message) {
	var message SSEMessage
	message.accountNumber = "55"
	message.room = "test"
	message.msg = formatSSE("testing", string(msg.Value))
	go func() {
		for messageChannel := range messageChannels {
			messageChannel <- message
		}
	}()
	fmt.Printf("%% Message on %s:\n%s\n", msg.TopicPartition, string(msg.Value))
}

func main() {
	topicsConfig := readTopics()
	_ = godotenv.Load()
	apiVersion := os.Getenv("API_VERSION")
	var topics []string
	for i := 0; i < len(topicsConfig); i++ {
		topics = append(topics, topicsConfig[i].Topic)
	}

	if apiVersion == "" {
		apiVersion = "v1"
	}

	go func(){
		connectKafka(topics, sendToListener)
	}()

	http.HandleFunc(fmt.Sprintf("/api/notifier/%s/connect", apiVersion), listenHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
