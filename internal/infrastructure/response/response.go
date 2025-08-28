package response

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

// JSON retorna uma resposta em JSON para a requisição
func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if dados != nil {
		if erro := json.NewEncoder(w).Encode(dados); erro != nil {
			log.Fatal(erro)
		}
	}

}

// Erro retorna um erro em formato JSON
func Erro(w http.ResponseWriter, statusCode int, erro error) {

	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
	slog.Error(erro.Error())
}

// JSON retorna uma resposta em JSON para a requisição
func Response(w http.ResponseWriter, statusCode int, dados interface{}, erro error) {
	if erro != nil {
		Erro(w, statusCode, erro)
		return
	}
	JSON(w, statusCode, dados)
}
