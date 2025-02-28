package main

import (
	"chatapp/wsmanager"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define WebSocket & MongoDB Variables
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var client *mongo.Client
var collection *mongo.Collection

// Each message will have an ID, sender, receiver, content, and timestamp.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MpID      string             `bson:"mp_id" json:"mp_id"`
	MID       string             `bson:"m_id" json:"m_id"`
	Content   string             `bson:"content" json:"content"`
	SenderID  string             `bson:"senderId" json:"senderId"`
	Status    int                `bson:"status" json:"status"`
	IsUnread  bool               `bson:"is_unread" json:"is_unread"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

// Handle WebSocket Connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading WebSocket:", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("No userId provided")
		return
	}

	wsmanager.Manager.AddConnection(userID, conn)
	defer wsmanager.Manager.RemoveConnection(userID, conn)

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		eventType, ok := msg["type"].(string)
		if !ok {
			log.Println("Invalid event type")
			continue
		}

		log.Println("Received message:", msg)

		switch eventType {
		case "message":
			handleMessage(msg)
		}
	}
}

// Handle Messages
func handleMessage(msg map[string]interface{}) {
	senderId, ok1 := msg["senderId"].(string)
	receiverId, ok2 := msg["receiverId"].(string)
	content, ok3 := msg["content"].(string)

	if !ok1 || !ok2 || !ok3 {
		log.Println("Invalid message format")
		return
	}

	id := primitive.NewObjectID()
	newMessage := Message{
		ID:        id,
		SenderID:  senderId,
		MpID:      receiverId,
		MID:       fmt.Sprintf("%s_%s", id.Hex(), senderId),
		Content:   content,
		Status:    1,
		IsUnread:  true,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	//Saves messages in MongoDB
	_, err := collection.InsertOne(context.TODO(), newMessage)
	if err != nil {
		log.Println("Error saving message:", err)
		return
	}

	//Broadcasts messages to sender & receiver
	wsmanager.Manager.BroadcastMessage(receiverId, msg)
	wsmanager.Manager.BroadcastMessage(senderId, msg)
}

// Allows frontend apps to connect to WebSocket
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Starts the WebSocket server
func main() {
	var err error

	// Replace with your MongoDB connection string
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with your MongoDB database and collection
	collection = client.Database("chatapp").Collection("messages")

	http.HandleFunc("/ws", enableCORS(handleConnections))

	log.Println("WebSocket server started on ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
