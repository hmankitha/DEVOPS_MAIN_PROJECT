package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("Test@12345"), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
