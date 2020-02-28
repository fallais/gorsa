package gorsa

// TokenResponse is a JWT response
type TokenResponse struct {
	ID           string   `json:"id"`
	Roles        []string `json:"roles"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}
