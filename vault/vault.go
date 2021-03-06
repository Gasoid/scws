package vault

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/vault/api"
)

var (
	vaultClient vaultService
)

type vaultService interface {
	Logical() *api.Logical
	SetToken(v string)
	Auth() *api.Auth
	Sys() *api.Sys
}

// renewToken is intended to renew token
func renewToken(client vaultService) {
	token := client.Auth().Token()
	secret, err := token.LookupSelf()
	if err != nil {
		log.Println("RenewToken: couldn't lookup", err.Error())
		return
	}
	if renewable, err := secret.TokenIsRenewable(); !renewable {
		log.Println("RenewToken: it seems token is not renewable", err)
		return
	}
	for {
		d, err := secret.TokenTTL()
		if err != nil {
			log.Println("RenewToken: couldn't get token TTL", err)

			break
		}
		timer := time.NewTimer(d / 2)
		<-timer.C
		_, err = token.RenewSelf(int(d / 2))
		if err != nil {
			log.Println("RenewToken: couldn't renew token", err)

			continue
		}
		log.Println("RenewToken: token has been renewed successfully")
	}
}

func Init(address, token string) error {
	log.Println("Vault address:", address)
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	vaultClient = client
	if err != nil {
		return fmt.Errorf("vault.api.NewClient failed: %v", err)
	}
	vaultClient.SetToken(token)
	_, err = vaultClient.Sys().Health()
	if err != nil {
		return fmt.Errorf("can't connect to vault: %v", err)
	}

	go renewToken(vaultClient)

	return nil
}

func Secrets(path string) (map[string]string, error) {
	secretMap := map[string]string{}
	secret, err := vaultClient.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault.Logical().Read() can't read secrets: %v", err)
	}
	if secret == nil {
		return nil, errors.New("secret data doesn't exist")
	}
	// if secret.Renewable {
	// 	api.
	// }
	if _, ok := secret.Data["data"]; !ok {
		return nil, errors.New("secret data doesn't exist")
	}
	data := secret.Data["data"].(map[string]interface{})

	for k, v := range data {
		secretMap[k] = v.(string)
	}
	return secretMap, nil
}
