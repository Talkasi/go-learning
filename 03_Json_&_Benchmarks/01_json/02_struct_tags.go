package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int `json:"user_id,string"`
	Username string
	Address  string `json:",omitempty"` // имя такое же, не включать в json, если пусто
	Comnpany string `json:"-"`			// не включать в json
}

func main() {
	u := &User{
		ID:       42,
		Username: "rvasily",
		Address:  "test",
		Comnpany: "Mail.Ru Group",
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}
