GRANT ALL PRIVILEGES ON *.* TO 'user'@'%';

CREATE DATABASE IF NOT EXISTS `database`;

CREATE TABLE IF NOT EXISTS `database`.`guests` (
	`id` INT NOT NULL auto_increment,
	`name` VARCHAR(255),
	`accompanying_guests` SMALLINT,
	`time_arrived` TIMESTAMP NULL DEFAULT NULL,
	`is_arrived` BOOLEAN DEFAULT false,
	PRIMARY KEY (`id`)
) ENGINE InnoDB DEFAULT CHARSET = `utf8`;

CREATE TABLE IF NOT EXISTS `database`.`tables` (
	`id` INT NOT NULL auto_increment,
	`seats` SMALLINT DEFAULT 0,
	`guest_id` INT NULL ,
	PRIMARY KEY (`id`),
	CONSTRAINT `fk_guest` FOREIGN KEY (`guest_id`) REFERENCES `database`.`guests`(`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE InnoDB DEFAULT CHARSET = `utf8`;

-- CREATE TABLE `guestlist` (
--   `guest_id` INT,
--   `table_id` INT
--   FOREIGN KEY (guest_id) REFERENCES `guest`(id) ON DELETE CASCADE
--   FOREIGN KEY (table_id) REFERENCES `table`(id) ON DELETE CASCADE
--   ON DELETE SET NULL  
-- ) ENGINE InnoDB DEFAULT CHARSET = `utf8`;