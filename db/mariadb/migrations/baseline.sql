-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- Dumping database structure for person
CREATE DATABASE IF NOT EXISTS `person` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `person`;

-- Dumping structure for table person.person
CREATE TABLE IF NOT EXISTS `person` (
  `id` char(50) NOT NULL,
  `name` char(50) NOT NULL,
  `last_name` char(50) NOT NULL,
  `phone` char(50) DEFAULT NULL,
  `email` char(50) DEFAULT NULL,
  `year_od_birth` year(4) DEFAULT NULL,
  KEY `Index 1` (`id`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `person`;
DROP DATABASE IF EXISTS `person`;