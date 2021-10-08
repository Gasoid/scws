package vault

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/vault/api"
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

func VaultInit(address, path, token string) error {
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		log.Println("Can't connect to vault: ", err.Error())
		return err
	}
	client.SetToken(token)
	// vault = &vaultServer{
	// 	client: client,
	// }
	secret, err := client.Logical().Read(path)
	if err != nil {
		log.Println("Can't pull secrets: ", err.Error())
		return err
	}
	if secret == nil {
		return errors.New("secret data doesn't exist")
	}
	data := secret.Data["data"].(map[string]interface{})
	go renewToken(client)
	if data != nil {
		return nil
	}
	// for k, v := range data {
	// 	envVars.secrets[k] = v.(string)
	// }
	return nil
}
