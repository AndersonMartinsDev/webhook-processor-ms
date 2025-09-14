package database

var (
	REDIS_ADDR     = ""
	REDIS_PORT     = ""
	REDIS_PASSWORD = ""
)

func SetRedisEnv(address, port, pass string) {
	REDIS_ADDR = address
	REDIS_PORT = port
	REDIS_PASSWORD = pass
}
