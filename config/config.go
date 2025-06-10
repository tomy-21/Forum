package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Le driver mysql
)

// InitDB initialise et retourne une connexion à la base de données.
func InitDB() (*sql.DB, error) {
	// Récupération des variables d'environnement
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// CORRECTION : Ajoutez "?parseTime=true" à la fin de la chaîne de connexion.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	// Ouverture de la connexion à la base de données
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de la connexion à la BDD: %w", err)
	}

	// Vérification de la validité de la connexion
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("erreur lors du ping de la BDD: %w", err)
	}

	log.Println("Connexion à la base de données réussie !")
	return db, nil
}
