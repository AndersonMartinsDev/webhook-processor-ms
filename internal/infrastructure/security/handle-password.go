package security

import "golang.org/x/crypto/bcrypt"

type HandlePassword struct {
}

func NewHandlerPassword() *HandlePassword {
	return &HandlePassword{}
}

// Hash recebe uma string e coloca um hash nela
func (handle HandlePassword) Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

// Verificar senha compara hash com string
func (handle HandlePassword) VerificarSenha(senhaString string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(senhaString), []byte(hash))
}
