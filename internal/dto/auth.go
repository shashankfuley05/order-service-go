package dto

type TokenRequest struct {
	UserID   string `json:"userId"  binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
