package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User 用户模型
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserStore 用户存储接口
type UserStore interface {
	Get(id string) (*User, error)
	Create(user *User) error
	List() ([]*User, error)
}

// InMemoryUserStore 内存存储实现
type InMemoryUserStore struct {
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	store := &InMemoryUserStore{
		users: make(map[string]*User),
	}
	// 添加示例数据
	store.users["1"] = &User{
		ID:        "1",
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
	}
	store.users["2"] = &User{
		ID:        "2",
		Username:  "bob",
		Email:     "bob@example.com",
		CreatedAt: time.Now(),
	}
	return store
}

func (s *InMemoryUserStore) Get(id string) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *InMemoryUserStore) Create(user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	user.CreatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryUserStore) List() ([]*User, error) {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

// UserHandler HTTP 处理器
type UserHandler struct {
	store UserStore
}

func NewUserHandler(store UserStore) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "user-service",
		"timestamp": time.Now().Unix(),
	})
}

func (h *UserHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "user-service",
	})
}

func main() {
	// 设置日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 设置 Gin 模式
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建存储
	store := NewInMemoryUserStore()
	handler := NewUserHandler(store)

	// 创建路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware())

	// 健康检查端点
	r.GET("/health", handler.HealthCheck)
	r.GET("/ready", handler.ReadinessCheck)

	// API 端点
	api := r.Group("/api/v1")
	{
		api.GET("/users", handler.ListUsers)
		api.GET("/users/:id", handler.GetUser)
		api.POST("/users", handler.CreateUser)
	}

	// 创建 HTTP 服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("User service started on port %s", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭，给予 5 秒处理现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// loggingMiddleware 日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %s %d %v",
			c.Request.Method,
			path,
			c.ClientIP(),
			status,
			latency,
		)
	}
}
