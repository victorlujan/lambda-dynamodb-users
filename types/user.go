package types

type User struct {
	ID   string `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
	Phone string `json:"phone"`

}