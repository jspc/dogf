package main

import (
	"fmt"
	
	"golang.org/x/crypto/bcrypt"
)

func main() {
	p, _ := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	
	fmt.Println(string(p))
}
