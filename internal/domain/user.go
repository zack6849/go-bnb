package domain

type User struct {
	UUIDIdentifier
	Email        string
	PasswordHash string
}
