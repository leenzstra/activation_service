CREATE TABLE IF NOT EXISTS `licenses` (
    `id` integer PRIMARY KEY,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime,
    `key` text UNIQUE,
    `max_uses` integer,
    `contacts` text,
    `expiration` datetime
)