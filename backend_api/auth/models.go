package auth

type AuthToken struct {
	Token string `json:"accesstoken"`
}

type LoginTokens struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type AccountAuth struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
