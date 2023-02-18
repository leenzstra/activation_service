CREATE TABLE IF NOT EXISTS `subject_classes` (
  `id` integer, 
  `created_at` datetime, 
  `updated_at` datetime, 
  `deleted_at` datetime, 
  `subject_id` integer, 
  `class` text, 
  PRIMARY KEY (`id`), 
  CONSTRAINT `fk_subjects_subject_classes` FOREIGN KEY (`subject_id`) REFERENCES `subjects`(`id`)
)
