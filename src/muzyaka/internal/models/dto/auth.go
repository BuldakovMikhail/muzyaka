package dto

type SignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type GetMeResponse struct {
	UserId     uint64 `json:"user_id"`
	Role       string `json:"role"`
	MusicianId uint64 `json:"musician_id,omitempty"`
}

type SignUpResponse struct {
	Token string `json:"token"`
}

type SignUp struct {
	UserInfo
}

type SignUpMusician struct {
	UserInfo
	MusicianWithoutId
}
