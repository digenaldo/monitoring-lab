package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	operationsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mongodb_operations_total",
			Help: "Total number of MongoDB operations",
		},
	)

	mongoLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "mongodb_operation_duration_seconds",
			Help:    "MongoDB operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(operationsCounter)
	prometheus.MustRegister(mongoLatency)

	// Note: Process and Go runtime metrics are automatically registered by promhttp.Handler()
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongo:27017"
	}

	// Connect to MongoDB
	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Test connection
	err = client.Ping(connectCtx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	db := client.Database("monitoring")
	collection := db.Collection("events")

	// Start background goroutine for periodic operations with a non-expiring context
	periodicCtx := context.Background()
	go periodicOperations(periodicCtx, client, collection)

	// HTTP server
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func periodicOperations(ctx context.Context, client *mongo.Client, collection *mongo.Collection) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Ping operation
			pingStart := time.Now()
			var result bson.M
			err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result)
			pingDuration := time.Since(pingStart).Seconds()

			if err != nil {
				log.Printf("Ping error: %v", err)
			} else {
				operationsCounter.Inc()
				mongoLatency.Observe(pingDuration)
				log.Printf("Ping successful, latency: %.3fs", pingDuration)
			}

			// Insert operation
			insertStart := time.Now()
			doc := bson.M{
				"timestamp": time.Now(),
				"source":    "go-app",
				"message":   "Periodic event",
			}
			_, err = collection.InsertOne(ctx, doc)
			insertDuration := time.Since(insertStart).Seconds()

			if err != nil {
				log.Printf("Insert error: %v", err)
			} else {
				operationsCounter.Inc()
				mongoLatency.Observe(insertDuration)
				log.Printf("Insert successful, latency: %.3fs", insertDuration)
			}

		case <-ctx.Done():
			return
		}
	}
}
