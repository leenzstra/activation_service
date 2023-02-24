CREATE TABLE IF NOT EXISTS `license_subject_classes` (
  `license_id` integer, 
  `subject_class_id` integer, 
  PRIMARY KEY (
    `license_id`, `subject_class_id`
  ), 
  CONSTRAINT `fk_license_subject_classes_license` FOREIGN KEY (`license_id`) REFERENCES `licenses`(`id`), 
  CONSTRAINT `fk_license_subject_classes_subject_class` FOREIGN KEY (`subject_class_id`) REFERENCES `subject_classes`(`id`)
)
