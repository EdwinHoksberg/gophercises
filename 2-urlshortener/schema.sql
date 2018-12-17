CREATE DATABASE `urlshortener`;

CREATE TABLE `urlshortener`.`paths` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `path` varchar(200) NOT NULL,
  `url` varchar(2000) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `path_UNIQUE` (`path`)
);

INSERT INTO `urlshortener`.`paths` (`path`, `url`) VALUES 
    ('/mysql-test', 'https://mysql.com'),
    ('/gitlab', 'https://gitlab.com/edwinhoksberg');
