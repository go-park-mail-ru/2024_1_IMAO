package models

const (
	SessionIDLen = 32
)

type UnauthorizedUser struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type CSRFToken struct {
	TokenBody string
}

type Session struct {
	UserID uint32
	Value  string
}

type DBInsertionUser struct {
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type AdditionalUserData struct {
	User        User   `json:"user"`
	Avatar      string `json:"avatarImg"`
	IsAuth      bool   `json:"isAuth"`
	CartNum     uint   `json:"cartNum"`
	FavNum      uint   `json:"favNum"`
	PhoneNumber string `json:"phoneNumber"`
}

type AuthResponse struct {
	User   User   `json:"user"`
	Avatar string `json:"avatarImg"`
	IsAuth bool   `json:"isAuth"`
}
