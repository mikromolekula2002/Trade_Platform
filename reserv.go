package tradeplatform

// Устанавливаем контекст, а также инициализируем редис и проверяем подключение
ctx := context.Background()
redisClient := cache.NewRedisClient(cfg)

// Проверка подключения
_, err := redisClient.Client.Ping(ctx).Result()
if err != nil {
	loger.Logrus.Errorf("ошибка подключения Redis: %v.\n Продолжаем работу без Редис.", err)
}
loger.Logrus.Debug("Redis успешно подключен.")

// Устанавливаем контекст, а также инициализируем редис и проверяем подключение
ctx := context.Background()
redisClient := cache.NewRedisClient(cfg)

// Проверка подключения
_, err := redisClient.Client.Ping(ctx).Result()
if err != nil {
	loger.Logrus.Errorf("ошибка подключения Redis: %v.\n Продолжаем работу без Редис.", err)
}
loger.Logrus.Debug("Redis успешно подключен.")

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
)

// Redis interface defines the methods for interacting with Redis.
type Redis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Increment(ctx context.Context, key string) error
}

// Структура с подключением к Redis
type RedisClient struct {
	Client *redis.Client
}

// Инит Редиса
func NewRedisClient(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &RedisClient{Client: client}
}

// Получение значения с Редиса
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("cache: Get - ошибка получения из Redis.\n ERROR: %v", err)
	}
	return val, nil
}

// Сохранение значения в Редис
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("cache: Set - ошибка установки в Redis.\n ERROR: %v", err)
	}
	return nil
}

// Удаление ключей в Redis
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	err := r.Client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("cache: Del - ошибка удаления из Redis.\n ERROR: %v", err)
	}
	return nil
}

func (r *RedisClient) Increment(ctx context.Context, key string) error {
	err := r.Client.Incr(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("cache: Incr - ошибка инкремента в Redis.\n ERROR: %v", err)
	}
	return nil
}

Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
} `yaml:"redis"`

type UserService struct {
	repo   repository.PostgreSQL
	cache  cache.RedisClient
	jwtKey string
}

// Инит нашей структуры с базой данных
func NewUserService(repo repository.PostgreSQL, jwtkey string, cache cache.RedisClient) *UserService {
	return &UserService{
		repo:   repo,
		cache:  cache,
		jwtKey: jwtkey,
	}
}

// Сервисная логика, подсчет и автоинкремент просмотров объявления
func (s *UserService) IncrViews(ctx context.Context, adsID string) error {
	// инкрементим число просмотров объявления
	err := s.cache.Increment(ctx, adsID)
	if err != nil {
		return fmt.Errorf("userservice.go: IncrViews:\n%v", err)
	}

	return nil
}