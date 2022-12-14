package config

import (
	"encoding/json"
	"os"
)

const key = "APP_ENV"
const defaultEnv = "dev"

type Conf struct {
	Mongo       MongoConf       `json:"mongo"`
	JWT         JWTConf         `json:"jwt"`
	Wallet      WalletConf      `json:"wallet"`
	Transaction TransactionConf `json:"transaction"`
}

type MongoConf struct {
	URI        string         `json:"uri"`
	Database   string         `json:"database"`
	Collection CollectionConf `json:"collection"`
}

type CollectionConf struct {
	Wallet      string `json:"wallet"`
	Transaction string `json:"transaction"`
}

type JWTConf struct {
	ValidityDurationInMin int    `json:"validityDurationInMin"`
	Issuer                string `json:"issuer"`
	Secret                string `json:"secret"`
}

type WalletConf struct {
	InitialBalance float64 `json:"initialBalance"`
	MaxBalance     float64 `json:"maxBalance"`
	MinBalance     float64 `json:"minBalance"`
}

type TransactionConf struct {
	MaxAmount float64 `json:"maxAmount"`
	MinAmount float64 `json:"minAmount"`
}

func Read(path string) (Conf, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return Conf{}, err
	}

	var c Conf
	err = json.Unmarshal(contentBytes, &c)
	return c, err
}

func GetEnvOrDefault() string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}

	return defaultEnv
}
