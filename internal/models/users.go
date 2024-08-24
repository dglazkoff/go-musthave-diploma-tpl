package models

// модели наверное надо будет разбивать, потому что разные поля нужно возвращать в разных запросах ??

// хранить UserRequest в хендлере где используется
/*
response := map[string]any{
"login":user.Login,
"password":user.Password,
}
*/
type User struct {
	Login    string  `json:"login"`
	Password string  `json:"-"`
	Balance  float64 `json:"balance"`
}
