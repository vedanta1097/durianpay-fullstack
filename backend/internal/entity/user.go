package entity

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

const (
	RoleCS        = "cs"
	RoleOperation = "operation"
)

func IsSupportedRole(role string) bool {
	return role == RoleCS || role == RoleOperation
}
