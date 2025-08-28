package configuration

import (
	"fmt"
	"log/slog"
	"net/http"
	"webhook-processor-ms/internal/infrastructure/commons/logger"

	"github.com/joho/godotenv"
)

var (
	Porta  = 8080
	Origin = ""
)

func LoadEnv() {
	if erro := godotenv.Load(); erro != nil {
		panic("Error ao carregar as variáveis de ambiente!")
	}
	slog.Info("Variáveis de ambiente carregadas com sucesso!")
}

// LoadLogger apenas para carregar logs personalizados
func LoadLogger() {
	custom_log := slog.New(logger.NewHandler(nil))
	slog.SetDefault(custom_log)
	slog.Info("Logger Carregado com sucesso!")
}

func LoadServer(routers http.Handler) {
	slog.Info(fmt.Sprintf("Servidor iniciado na porta %d", Porta))
	if erro := http.ListenAndServe(fmt.Sprintf(":%d", Porta), routers); erro != nil {
		panic(fmt.Sprintf("Error ao iniciar servidor %s", erro.Error()))
	}
}
