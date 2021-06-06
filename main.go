package main

import (
	"fmt"
	"log"

	"github.com/Akshit8/go-vault/secret"
)

func main() {
	vaultPath := "/secret"                // NOTE: use config tool of your choice to load this value
	vaultToken := "myroot"                // NOTE: use config tool of your choice to load this value
	vaultAddress := "http://0.0.0.0:8300" // NOTE: use config tool of your choice to load this value

	provider, err := secret.NewVaultProvider(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		log.Fatalln("Couldn't load vault provider", err)
	}

	DBUsernameKey := "/database:username"
	DBPasswordKey := "/database:password"

	DBUsernameValue, err := provider.Get(DBUsernameKey)
	if err != nil {
		log.Fatalln("Couldn't load username", err)
	}

	DBPasswordValue, err := provider.Get(DBPasswordKey)
	if err != nil {
		log.Fatalln("Couldn't load password", err)
	}

	fmt.Println("DB username:", DBUsernameValue)
	fmt.Println("DB password:", DBPasswordValue)
}
