// NATS JetStream Go 客户端示例
// 展示生产者、消费者和流管理

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Order 订单结构
type Order struct {
	ID        string    `json:"id"`
	Customer  string    `json:"customer"`
	Items     []Item    `json:"items"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func main() {
	ctx := context.Background()

	// 连接到 NATS
	nc, err := nats.Connect(
		nats.DefaultURL,
		nats.Name("Order Service"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected to %s", nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 创建 JetStream 上下文
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// 创建流
	if err := createStream(ctx, js); err != nil {
		log.Fatal(err)
	}

	// 启动消费者（并发）
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		consumeOrders(ctx, js)
	}()

	go func() {
		defer wg.Done()
		publishOrders(ctx, js)
	}()

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}

// createStream 创建订单流
func createStream(ctx context.Context, js jetstream.JetStream) error {
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:            "ORDERS",
		Subjects:        []string{"orders.*"},
		Storage:         jetstream.FileStorage,
		Replicas:        3,
		Retention:       jetstream.LimitsPolicy,
		MaxMsgs:         1_000_000,
		MaxBytes:        5 * 1024 * 1024 * 1024, // 5GB
		MaxAge:          30 * 24 * time.Hour,    // 30天
		DuplicateWindow: 5 * time.Minute,
	})
	if err != nil {
		// 流可能已存在
		stream, err = js.Stream(ctx, "ORDERS")
		if err != nil {
			return err
		}
	}

	info, _ := stream.Info(ctx)
	log.Printf("Stream: %s (msgs: %d, bytes: %d)",
		info.Config.Name, info.State.Msgs, info.State.Bytes)

	return nil
}

// publishOrders 发布订单消息
func publishOrders(ctx context.Context, js jetstream.JetStream) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			order := Order{
				ID:       fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
				Customer: "customer@example.com",
				Items: []Item{
					{ProductID: "PROD-001", Name: "Widget", Quantity: 2, Price: 29.99},
					{ProductID: "PROD-002", Name: "Gadget", Quantity: 1, Price: 49.99},
				},
				Total:     109.97,
				Status:    "created",
				CreatedAt: time.Now(),
			}

			data, _ := json.Marshal(order)

			// 发布消息（幂等，使用 MsgID）
			ack, err := js.Publish(ctx, "orders.created", data,
				jetstream.WithMsgID(order.ID),
				jetstream.WithExpectStream("ORDERS"),
			)
			if err != nil {
				log.Printf("Publish error: %v", err)
				continue
			}

			log.Printf("Published order %s (seq: %d)", order.ID, ack.Sequence)
		}
	}
}

// consumeOrders 消费订单消息
func consumeOrders(ctx context.Context, js jetstream.JetStream) {
	// 创建消费者
	cons, err := js.CreateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Name:          "order-processor",
		Durable:       "order-processor",
		FilterSubject: "orders.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		// 消费者可能已存在
		cons, err = js.Consumer(ctx, "ORDERS", "order-processor")
		if err != nil {
			log.Printf("Consumer error: %v", err)
			return
		}
	}

	// 消费消息
	consContext, err := cons.Consume(func(msg jetstream.Msg) {
		var order Order
		if err := json.Unmarshal(msg.Data(), &order); err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			msg.Nak() // 否定确认，重新投递
			return
		}

		log.Printf("Processing order %s (total: %.2f)", order.ID, order.Total)

		// 模拟处理
		time.Sleep(500 * time.Millisecond)

		// 确认消息
		if err := msg.Ack(); err != nil {
			log.Printf("Failed to ack: %v", err)
		}
	})
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}
	defer consContext.Stop()

	log.Println("Consumer started...")
	<-ctx.Done()
}

// 使用 KeyValue 存储配置
func useKeyValue(ctx context.Context, js jetstream.JetStream) {
	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "CONFIG",
		Storage: jetstream.FileStorage,
		TTL:     time.Hour,
	})
	if err != nil {
		log.Printf("KV create error: %v", err)
		return
	}

	// 存储配置
	if _, err := kv.PutString(ctx, "app.theme", "dark"); err != nil {
		log.Printf("KV put error: %v", err)
		return
	}

	// 读取配置
	entry, err := kv.Get(ctx, "app.theme")
	if err != nil {
		log.Printf("KV get error: %v", err)
		return
	}

	log.Printf("Config value: %s", string(entry.Value()))

	// 监听配置变化
	watcher, err := kv.WatchAll(ctx)
	if err != nil {
		log.Printf("KV watch error: %v", err)
		return
	}

	for entry := range watcher.Updates() {
		if entry == nil {
			continue
		}
		log.Printf("Config changed: %s = %s", entry.Key(), string(entry.Value()))
	}
}
