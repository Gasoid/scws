package vault

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/vault/api"
)

var (
	vaultClient *api.Client
)

// renewToken is intended to renew token
func renewToken(client *api.Client) {
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
		log.Println("Can't connect to vault: ", err.Error())
		return err
	}
	vaultClient.SetToken(token)

	go renewToken(vaultClient)

	return nil
}

func GetSecrets(path string) (map[string]string, error) {
	secretMap := map[string]string{}
	secret, err := vaultClient.Logical().Read(path)
	if err != nil {
		log.Println("Can't pull secrets: ", err.Error())
		return nil, err
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
