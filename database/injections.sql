-- Insertion des rôles
INSERT INTO `roles` (`role_id`, `name`) VALUES
(1, 'Administrateur'),
(2, 'Modérateur'),
(3, 'Utilisateur');

-- Insertion des utilisateurs
-- Mot de passe fictif : 'password123' (dans une vraie application, il faudrait le hacher !)
INSERT INTO `Utilisateurs` (`user_id`, `role_id`, `name`, `email`, `password`) VALUES
(1, 1, 'AliceAdmin', 'alice.admin@example.com', '$2y$10$NotARealHashJustForExample1'),
(2, 2, 'BobLeModo', 'bob.modo@example.com', '$2y$10$NotARealHashJustForExample2'),
(3, 3, 'CharlieUser', 'charlie.user@example.com', '$2y$10$NotARealHashJustForExample3'),
(4, 3, 'DavidDev', 'david.dev@example.com', '$2y$10$NotARealHashJustForExample4'),
(5, 3, 'EveGamer', 'eve.gamer@example.com', '$2y$10$NotARealHashJustForExample5');

-- Insertion des catégories
INSERT INTO `categories` (`category_id`, `name`, `description`, `created_at`) VALUES
(1, 'Technologie', 'Discussions sur le matériel, les logiciels et les nouvelles technologies.', '2024-10-01 10:00:00'),
(2, 'Jeux Vidéo', 'Tout sur les jeux PC, consoles et mobiles. Partagez vos astuces et vos critiques !', '2024-10-01 11:00:00'),
(3, 'Discussion Générale', 'Un espace pour parler de tout et de rien.', '2024-10-01 12:00:00');

-- Insertion des forums
INSERT INTO `forums` (`forum_id`, `category_id`, `name`, `description`) VALUES
(1, 1, 'Hardware & Composants', 'Discutez des dernières cartes graphiques, processeurs, etc.'),
(2, 1, 'Développement Web', 'Pour les passionnés de HTML, CSS, JavaScript, PHP, et plus encore.'),
(3, 2, 'PC Gaming', 'L\'univers du jeu sur ordinateur.'),
(4, 3, 'Le Bistrot', 'Discussions libres et hors-sujet.');

-- Insertion des sujets (topics)
-- Note: status = 1 pour Ouvert, 0 pour Fermé
INSERT INTO `sujet` (`topic_id`, `forum_id`, `user_id`, `title`, `status`) VALUES
(1, 2, 4, 'Quel est votre framework JS préféré en 2025 ?', 1),
(2, 3, 5, 'Recommandations de RPG pour PC', 1),
(3, 4, 3, 'Votre film préféré de l\'année ?', 1),
(4, 1, 4, 'Aide pour upgrade ma config PC', 0); -- Sujet fermé

-- Insertion des messages
-- Sujet 1: Framework JS
INSERT INTO `messages` (`message_id`, `topic_id`, `user_id`, `content`) VALUES
(1, 1, 4, 'Salut tout le monde ! Je me lance dans un nouveau projet et j\'hésite entre React, Vue et Svelte. Quels sont vos avis et retours d\'expérience ?'),
(2, 1, 1, 'Personnellement, je trouve que Svelte est vraiment innovant et offre des performances incroyables pour des projets de petite et moyenne taille.'),
(3, 1, 4, 'Intéressant ! Et pour des plus gros projets, tu resterais sur React ?'),
(4, 1, 2, 'Pour les grosses applications d\'entreprise, l\'écosystème de React est imbattable. C\'est un choix plus sûr.');

-- Sujet 2: RPG PC
INSERT INTO `messages` (`message_id`, `topic_id`, `user_id`, `content`) VALUES
(5, 2, 5, 'Hello les gamers ! Je cherche un bon RPG pour m\'occuper cet hiver. Des idées ? J\'ai déjà fait The Witcher 3 et Skyrim.'),
(6, 2, 3, 'Si tu aimes les classiques, je te conseille fortement Disco Elysium. L\'écriture est absolument géniale !'),
(7, 2, 5, 'Oh oui j\'en ai entendu parler ! Merci pour la suggestion, je vais regarder ça.');

-- Sujet 3: Films
INSERT INTO `messages` (`message_id`, `topic_id`, `user_id`, `content`) VALUES
(8, 3, 3, 'Alors, quel film vous a le plus marqué cette année ? Pour moi, c\'est sans hésiter "Dune: Part Two". Une claque visuelle !');

-- Insertion des réactions (likes/dislikes)
INSERT INTO `reaction` (`feedback_id`, `user_id`, `message_id`, `type`) VALUES
(1, 1, 6, 'like'),   -- Alice aime le message de Charlie sur Disco Elysium
(2, 4, 2, 'like'),   -- David aime la suggestion de Svelte par Alice
(3, 5, 6, 'like'),   -- Eve aime aussi la suggestion pour Disco Elysium
(4, 3, 4, 'dislike'); -- Charlie n'est pas d'accord avec le message de Bob sur React

-- Insertion des réponses (table `Reponses` pour des réponses directes à un message)
-- Ceci simule un système de "répondre à"
INSERT INTO `Reponses` (`reply_id`, `reply_to_id`, `content`) VALUES
(1, 2, 'Merci pour ton avis détaillé sur Svelte, Alice ! C\'est très éclairant.'), -- David (user 4) répond directement au message 2 d'Alice
(2, 6, 'Je plussoie ! Disco Elysium est un chef-d\'œuvre.'); -- Un autre utilisateur pourrait répondre directement à la suggestion de Charlie