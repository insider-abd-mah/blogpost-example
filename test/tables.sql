DROP DATABASE IF EXISTS sample_db;
create database sample_db;

use sample_db;

CREATE TABLE IF NOT EXISTS `posts`
(
    `id`                INT(11) AUTO_INCREMENT PRIMARY KEY,
    `title`             VARCHAR(255) NOT NULL,
    `description`       VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
