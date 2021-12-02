package views

var (
	ErrInvalidCredentials = Error{Code: "invalid_credentials", Message: "Login ou senha incorretos."}
	ErrInvalidEmail       = Error{Code: "invalid_email", Message: "Esse e-mail já está cadastrado."}
	ErrInvalidUser        = Error{Code: "invalid_user", Message: "Esse usuário já está cadastrado."}
	ErrUnverifiedUser     = Error{Code: "unverified_user", Message: "Seu e-mail precisa ser verificado para utilizar o USPY."}
	ErrBannedUser         = Error{Code: "banned_user", Message: "Infelizmente sua conta foi banida."}
	ErrWrongPassword      = Error{Code: "invalid_password", Message: "Senha incorreta"}
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
