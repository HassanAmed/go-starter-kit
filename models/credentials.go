package models

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=0"`
}
