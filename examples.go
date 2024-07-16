package golastic

import "fmt"

type tweet struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
}

func examples() {
	err := Connect("http://0.0.0.0:9200")
	if err != nil {
		panic(err)
	}
	data := []tweet{}
	Save("twitter", tweet{ID: 1, User: "fulanah", Message: "Hello World!"})
	Save("twitter", tweet{ID: 2, User: "fulanah", Message: "Hello World!"})
	Save("twitter", tweet{ID: 3, User: "fulanah", Message: "Hello World!"})
	Delete("twitter", "id", 1)
	Search(&data, "twitter", "fulan", "user", "message")
	for _, x := range data {
		fmt.Println(x.User, x.Message)
	}
}
