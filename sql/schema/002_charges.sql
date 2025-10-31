-- +goose Up
CREATE TABLE `charges` (
    `id` int AUTO_INCREMENT PRIMARY KEY,
    `payment_intent_id` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `amount` int NOT NULL,
    `user_email` varchar(255) NOT NULL,
    `created_at` timestamp NOT NULL,
    `updated_at` timestamp NOT NULL
);

-- +goose Down
DROP TABLE charges;
