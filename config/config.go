package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Environment   string `json:"environment"`
	ApiPort       string `json:"api-port"`
	ApiName       string `json:"api-name"`
	Email         string `json:"email"`
	EmailPassword string `json:"email-password"`
	EmailDomain   string `json:"email-domain"`
	EmailServer   string `json:"email-server"`
	EmailPort     string `json:"email-port"`
	Site          string `json:"site"`
}

/* ****************************************************
**	Exported functions
** ***************************************************/

func Get(path string) Configuration {
	// Tries to open the configuration file
	file, err := os.Open(path)
	if err != nil {
		// If it can't open the configuration file
		log.Println("Arquivo de configuração inválido.")
		log.Fatalln(err)
	}

	// Tries to decode the json file
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		// If it can't decode the configuration file
		log.Println("Error decoding the configuration file...")
		log.Fatalln(err)
	}

	return conf
}
