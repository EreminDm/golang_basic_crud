-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               10.4.8-MariaDB - mariadb.org binary distribution
-- Server OS:                    Win64
-- HeidiSQL Version:             10.2.0.5599
-- --------------------------------------------------------
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Creating user
CREATE USER 'test'@'%' IDENTIFIED BY 'test';
GRANT SHOW DATABASES, SELECT, PROCESS, EXECUTE, ALTER ROUTINE, ALTER, SHOW VIEW, CREATE TABLESPACE, CREATE ROUTINE, CREATE, DELETE, CREATE VIEW, CREATE TEMPORARY TABLES, INDEX, EVENT, DROP, TRIGGER, REFERENCES, INSERT, FILE, CREATE USER, UPDATE, RELOAD, LOCK TABLES, SHUTDOWN, REPLICATION SLAVE, REPLICATION CLIENT, SUPER ON *.* TO 'test'@'%';
FLUSH PRIVILEGES;

-- Dumping database structure for person
CREATE DATABASE IF NOT EXISTS `person` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `person`;

-- Dumping structure for table person.person
CREATE TABLE IF NOT EXISTS `person` (
  `id` binary(50) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
  `name` char(50) NOT NULL,
  `last_name` char(50) NOT NULL,
  `phone` char(50) DEFAULT NULL,
  `email` char(50) DEFAULT NULL,
  `year_od_birth` year(4) DEFAULT NULL,
  KEY `Index 1` (`id`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Data exporting was unselected.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `person`;
DROP DATABASE IF EXISTS `person`;