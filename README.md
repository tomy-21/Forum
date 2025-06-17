# 🧵 Forum - Projet Ynov Aix-en-Provence (B1)

Ce dépôt contient le projet **Forum** développé en Go avec une architecture **MVC**, dans le cadre du Bachelor 1 à Ynov Aix-en-Provence.

## 👥 Équipe

BOURDOT Romain – Développement Backend / BDD
OUADHI Tomy – Développement Frontend / Contrôleurs

## 🚀 Fonctionnalités Principales

* **Authentification sécurisée :**

  * Inscription avec contraintes (mot de passe fort, email unique)
  * Connexion via email ou pseudo
  * Sessions gérées par **JWT**
  * Mots de passe hachés avec `bcrypt`

* **Gestion des contenus :**

  * Création de sujets et réponses
  * Consultation libre des sujets
  * Upload d’image dans les messages

* **Interactivité & Modération :**

  * Likes/Dislikes sans doublons
  * Tri des messages (date/popularité)
  * Suppression de ses propres messages/sujets

* **Navigation fluide :**

  * Pagination des sujets
  * Système de catégories
  * Barre de recherche par titre

* **Administration :**

  * Espace `/admin` sécurisé
  * Ban/déban utilisateurs
  * Gestion des contenus et statut des sujets

## ⚙️ Installation

### Prérequis

* Go ≥ 1.18
* MySQL

### Étapes

1. **Cloner le dépôt**

   ```bash
   git clone [URL_DU_DEPOT]
   cd [DOSSIER_DU_PROJET]
   ```

2. **Configurer MySQL**

   ```sql
   CREATE DATABASE IF NOT EXISTS forum_db;
   ```

   ```bash
   mysql -u [USER] -p forum_db < ./database/migration.sql
   # (Optionnel) Données de test :
   mysql -u [USER] -p forum_db < ./database/injections.sql
   ```

3. **Créer `.env` à la racine :**

   ```env
   DB_USER="..."
   DB_PASSWORD="..."
   DB_HOST="localhost"
   DB_PORT="3306"
   DB_NAME="forum_db"
   JWT_SECRET_KEY="CLE_ULTRA_SECRETE"
   ```

4. **Lancer le serveur**

   ```bash
   go run .
   ```

📍 Accès à l’application : [http://localhost:8080](http://localhost:8080)

## 🗺️ Routes Principales

### Vues (HTML)

| Route                 | Description              | Accès       |
| --------------------- | ------------------------ | ----------- |
| `/`                   | Accueil (sujets paginés) | Public      |
| `/login`, `/register` | Authentification         | Public      |
| `/category/{id}`      | Sujets par catégorie     | Public      |
| `/topic/{id}`         | Affichage d’un sujet     | Public      |
| `/profil`             | Profil utilisateur       | Authentifié |
| `/admin`              | Dashboard admin          | Admin       |

### Actions (POST/GET)

| Route                            | Action                       | Accès       |
| -------------------------------- | ---------------------------- | ----------- |
| `/login`, `/register`, `/logout` | Authentification             | Public      |
| `/category/{id}/topics/create`   | Créer un sujet               | Authentifié |
| `/topic/{id}/reply`, `/react`    | Répondre / réagir            | Authentifié |
| `/message/{id}/delete`           | Supprimer message            | Authentifié |
| `/admin/users/ban/{id}`          | Ban/Déban utilisateur        | Admin       |
| `/admin/topics/status/{id}`      | Changer le statut d’un sujet | Admin       |

---
