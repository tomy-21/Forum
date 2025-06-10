package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	errLoad := godotenv.Load(".env")
	if errLoad != nil {
		log.Fatalf("Erreur chargement fichier d'environemment - Impossible de lancer le programme \n\t Erreur : %s", errLoad.Error())
	}
}

func GetEnvWithDefault(key, defaultValue string) string {
	envVar, envErr := os.LookupEnv(key)
	if !envErr {
		return defaultValue
	}
	return envVar
}
