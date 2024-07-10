package golastic

import "fmt"

type tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func examples() {
	err := Connect("http://0.0.0.0:9200")
	if err != nil {
		panic(err)
	}
	data := []tweet{}
	Save("twitter", tweet{User: "fulanah", Message: "Hello World!"})
	Search(&data, "twitter", "fulan", "user", "message")
	for _, x := range data {
		fmt.Println(x.User, x.Message)
	}
}
