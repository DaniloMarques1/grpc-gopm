package model

type Password struct {
	Id  string `json:"id"`
	Key string `json:"key"`
	Pwd string `json:"pwd"`
}

type PasswordRepository interface {
	Save(*Password) error
	FindAllKeys() ([]string, error)
	FindPassword(string) (*Password, error)
	Delete(string) error
	Update(*Password) error
}
