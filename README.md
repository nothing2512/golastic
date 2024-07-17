# Golastic
golang elastic search plugin

## Usage
```sh
go get -u github.com/nothing2512/golastic
```

## Example
```go
package main

import (
    "fmt"
    "github.com/nothing2512/golastic"
)

type Test struct {
	ID   int    `json:"id" gl:"id"`
	Name string `json:"name" gl:"name"`
}

func (*Test) TableName() string {
	return "tests"
}

func main() {
	err := golastic.Connect("http://0.0.0.0:9200")
	if err != nil {
		panic(err)
	}
	data := []Test{}
	golastic.Save(&Test{1, "Fulan"})
	golastic.Save(&Test{2, "Fulan"})
	golastic.Save(&Test{3, "Fulan"})
	golastic.Update(&Test{1, "Fulanah"})
	golastic.Delete("tests", 2)
	golastic.Search(&data, "tests", "fulan", "name", "message")
	for _, x := range data {
		fmt.Println(x.ID, x.Name)
	}
}
```