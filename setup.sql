CREATE DATABASE IF NOT EXISTS recipe_db;

CREATE USER IF NOT EXISTS 'recipe_user'@'localhost' IDENTIFIED BY 'password';

GRANT ALL PRIVILEGES ON recipe_db.* TO 'recipe_user'@'localhost';

FLUSH PRIVILEGES;
