package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Order struct {
	ID       uuid.UUID `json:"id"`
	Product  string    `json:"product"`
	Quantity int       `json:"quantity"`
}

type OrderCreatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	Product string    `json:"product"`
}

type OutboxMessage struct {
	ID      uuid.UUID `json:"id"`
	Topic   string    `json:"topic"`
	Message []byte    `json:"message"`
}

func createOrder(ctx context.Context, pool *pgxpool.Pool, order Order) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert the order into the orders table
	_, err = tx.Exec(ctx, "INSERT INTO orders (id, product, quantity) VALUES ($1, $2, $3)", order.ID, order.Product, order.Quantity)
	if err != nil {
		return err
	}

	// Create the outbox message
	event := OrderCreatedEvent{
		OrderID: order.ID,
		Product: order.Product,
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Insert the outbox message
	_, err = tx.Exec(ctx, "INSERT INTO outbox (id, topic, message, state) VALUES ($1, $2, $3, 'pending')", uuid.New(), "order.created", msgBytes)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func relay(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT id, topic, message
		FROM outbox
		WHERE state = 'pending'
		order by created_at
		limit 1
		for update skip locked
	`)
	if err != nil {
		return err
	}

	// Process the retrieved outbox messages
	var msg OutboxMessage

	if !rows.Next() {
		rows.Close()
		return nil // No pending messages
	}
	if err := rows.Scan(&msg.ID, &msg.Topic, &msg.Message); err != nil {
		return err
	}

	rows.Close()

	// Publish the message to Redis
	if err := rdb.Publish(ctx, msg.Topic, msg.Message).Err(); err != nil {
		return err
	}
	log.Printf("Published message %s to topic %s", msg.ID, msg.Topic)

	// Mark the message as processed
	if _, err := tx.Exec(ctx, "UPDATE outbox SET state = 'processed', processed_at = NOW() WHERE id = $1", msg.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func startRelay(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := relay(context.Background(), pool, rdb); err != nil {
			log.Printf("Relay error: %v", err)
		}
	}
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:123456@localhost:5432/mydb")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	go startRelay(ctx, pool, rdb)

	order1 := Order{
		ID:       uuid.New(),
		Product:  "Widget",
		Quantity: 10,
	}
	if err := createOrder(ctx, pool, order1); err != nil {
		log.Printf("Failed to create order: %v", err)
	}

	log.Printf("Order created with ID: %s", order1.ID)

	time.Sleep(5 * time.Second)
	order2 := Order{
		ID:       uuid.New(),
		Product:  "Gadget",
		Quantity: 5,
	}
	if err := createOrder(ctx, pool, order2); err != nil {
		log.Printf("Failed to create order: %v", err)
	}

	log.Printf("Order created with ID: %s", order2.ID)
	time.Sleep(5 * time.Second)
}
