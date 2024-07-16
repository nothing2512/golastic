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
	Save("twitter", 1, tweet{ID: 1, User: "fulanah", Message: "Hello World!"})
	Save("twitter", 2, tweet{ID: 2, User: "fulanah", Message: "Hello World!"})
	Save("twitter", 3, tweet{ID: 3, User: "fulanah", Message: "Hello World!"})
	Delete("twitter", 1)
	Update("twitter", 1, tweet{ID: 4, User: "Achmad", Message: "Halo Dunia"})
	Search(&data, "twitter", "fulan", "user", "message")
	for _, x := range data {
		fmt.Println(x.User, x.Message)
	}
}
