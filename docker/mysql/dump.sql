GRANT ALL PRIVILEGES ON *.* TO 'user'@'%';

CREATE TABLE `guests` (
	`id` INT NOT NULL auto_increment,
	`name` VARCHAR(255),
	`accompanying_guests` SMALLINT,
	`time_arrived` TIMESTAMP,
	PRIMARY KEY (`id`)
) ENGINE InnoDB DEFAULT CHARSET = `utf8`;

CREATE TABLE `tables` (
	`id` INT NOT NULL auto_increment,
	`seats` SMALLINT DEFAULT 0,
	`is_busy` BOOLEAN DEFAULT 0,
	`guest_id` INT REFERENCES `guest`(`id`) ON DELETE SET null,
	INDEX par_ind (`guest_id`),
	PRIMARY KEY (`id`)
) ENGINE InnoDB DEFAULT CHARSET = `utf8`;

-- CREATE TABLE `guestlist` (
--   `guest_id` INT,
--   `table_id` INT
--   FOREIGN KEY (guest_id) REFERENCES `guest`(id) ON DELETE CASCADE
--   FOREIGN KEY (table_id) REFERENCES `table`(id) ON DELETE CASCADE
--   ON DELETE SET NULL  
-- ) ENGINE InnoDB DEFAULT CHARSET = `utf8`;