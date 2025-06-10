package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	// Récupération des parametres lié à la base de données (stocké dans un fichier d'environemment)
	// Possibilité d'ajouter des valeur par défaut
	user := GetEnvWithDefault("DB_USER", "")
	pwd := GetEnvWithDefault("DB_PWD", "")
	host := GetEnvWithDefault("DB_HOST", "")
	port := GetEnvWithDefault("DB_PORT", "")
	name := GetEnvWithDefault("DB_NAME", "")

	if user == "" || host == "" || port == "" || name == "" {
		return nil, fmt.Errorf(" Erreur connection base de données - Donneés de connexions manquantes")
	}

	// Préparation de la chaîne de connexion à la base de données
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, name)

	// Mise en place de la connexion
	dbContext, dbContextErr := sql.Open("mysql", connectionString)
	if dbContextErr != nil {
		return nil, fmt.Errorf(" Erreur connection base de données - Erreur : \n\t %s", dbContextErr.Error())
	}

	// Test de ping la base de données
	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		return nil, fmt.Errorf(" Erreur ping base de données - Erreur : \n\t %s", pingErr.Error())
	}

	return dbContext, nil
}
