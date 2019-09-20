CREATE DATABASE IF NOT EXISTS person;

CREATE TABLE IF NOT EXISTS `person`.`personal_data` (
	`id` CHAR(50) NOT NULL,
	`name` CHAR(50) NOT NULL,
	`last_name` CHAR(50) NOT NULL,
	`phone` CHAR(50) NULL DEFAULT NULL,
	`email` CHAR(50) NULL DEFAULT NULL,
	`year_od_birth` INT(11) NULL DEFAULT NULL
);
