package server

type Objects struct {
	Users map[string]User `json:"users"`
	Apps  map[string]App  `json:"apps"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (o User) GetName() string {
	return o.Name
}

func (o User) GetType() string {
	return "User"
}

type App struct {
	Name    string  `json:"name"`
	Created int     `json:"created"`
	Price   float64 `json:"price"`
}

func (o App) GetName() string {
	return o.Name
}

func (o App) GetType() string {
	return "App"
}
