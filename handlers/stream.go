package handlers

import (
	"RTDS_API/config"
	"RTDS_API/services"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Response struct {
	Message string `json:"message"`
}

var resultsChannels = make(map[string]chan string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func StartStreamHandler(w http.ResponseWriter, r *http.Request) {
	streamID := uuid.New().String()
	log.Printf("[Stream %s] Initializing...", streamID)

	if err := services.InitKafkaProducer(config.BrokerAddress, streamID); err != nil {
		log.Printf("[Stream %s] Error initializing Kafka producer: %v", streamID, err)
		http.Error(w, "Failed to start stream", http.StatusInternalServerError)
		return
	}

	resultsChan := make(chan string)
	resultsChannels[streamID] = resultsChan

	go services.ConsumeMessages(streamID, resultsChan)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Stream started with ID: " + streamID})
}

func GetResultsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	resultsChan, exists := resultsChannels[streamID]
	if !exists {
		http.Error(w, "Stream results not available", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket setup failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for result := range resultsChan {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(result)); err != nil {
			break
		}
	}
}

func SendDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := services.ProduceMessage(streamID, string(body)); err != nil {
		http.Error(w, "Failed to send data to Kafka", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "Data sent to stream ID: " + streamID})
}
