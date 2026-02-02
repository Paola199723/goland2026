package dto

type AuthTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
}

type AuthErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
