package database

import (
	"context"
	"fmt"
	"log"
	"time" // Para o timeout de conexão

	"github.com/go-redis/redis/v8"
)

var CACHE redis.Client

// RedisClient é o cliente Redis configurado.
type RedisClient struct {
	client *redis.Client
}

func InitializeCache() {
	address := fmt.Sprintf(`%s:%s`, REDIS_ADDR, REDIS_PORT)
	redisDB := 0
	redisConnTimeout := 5 * time.Second
	redis_client, erro := newRedisClient(address, REDIS_PASSWORD, redisDB, redisConnTimeout)
	if erro != nil {
		panic(erro)
	}
	CACHE = *redis_client.client
	log.Println("Redis inicializado com sucesso")
}

func newRedisClient(addr, password string, db int, timeout time.Duration) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		// Adicione opções de pool de conexões e timeouts se necessário para produção
		PoolSize:     10, // Número máximo de conexões ociosas e ativas
		PoolTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Testa a conexão
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao Redis em %s: %w", addr, err)
	}

	log.Printf("Conectado com sucesso ao Redis em %s\n", addr)
	return &RedisClient{client: rdb}, nil
}

// GetClient retorna a instância interna do cliente redis-go.
// Pode ser útil para operações mais avançadas, mas tente usar os métodos do repositório.
func (rc *RedisClient) GetClient() *redis.Client {
	return rc.client
}

// Close fecha a conexão com o Redis.
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}
