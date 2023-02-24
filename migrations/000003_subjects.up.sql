CREATE TABLE IF NOT EXISTS `subjects` (
  `id` integer, 
  `created_at` datetime, 
  `updated_at` datetime, 
  `deleted_at` datetime, 
  `sid` integer UNIQUE, 
  `name` text, 
  `alias` text UNIQUE, 
  PRIMARY KEY (`id`)
)