package server

type Objects struct {
	Users map[string]User `json:"users"`
	Apps  map[string]App  `json:"apps"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u User) GetName() string {
	return u.Name
}

func (u User) GetType() string {
	return "User"
}

type App struct {
	Name    string  `json:"name"`
	Created int     `json:"created"`
	Price   float64 `json:"price"`
}

func (u App) GetName() string {
	return u.Name
}

func (u App) GetType() string {
	return "App"
}
