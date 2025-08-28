package security

var (
	SECRET_KEY string
)

func GetSecretKeyConfig(secretKey string) {
	SECRET_KEY = secretKey
}
