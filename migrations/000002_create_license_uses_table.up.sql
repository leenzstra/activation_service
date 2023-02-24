CREATE TABLE IF NOT EXISTS `license_uses` (
    `id` integer PRIMARY KEY,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime,
    `license_id` integer,
    `machine_info_hash` text,
    CONSTRAINT `fk_licenses_license_uses` FOREIGN KEY (`license_id`) REFERENCES `licenses`(`id`)
)