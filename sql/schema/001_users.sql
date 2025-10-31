-- +goose Up
CREATE TABLE `users` (
    `id` int AUTO_INCREMENT PRIMARY KEY,
    `email` varchar(255) NOT NULL,
    `hwid` varchar(255) DEFAULT NULL,
    `created_at` timestamp NOT NULL,
    `updated_at` timestamp NOT NULL
);

INSERT INTO `users` (email) VALUES ('ramontanoeiro@gmail.com');

-- +goose Down
DROP TABLE users;
