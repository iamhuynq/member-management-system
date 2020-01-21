USE vista;
ALTER SCHEMA vista DEFAULT CHARACTER SET utf8;

/**
- members table
**/
CREATE TABLE `members` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `login_id` varchar(100) NOT NULL,
    `password` varchar(200) NOT NULL,
    `name` varchar(40) NOT NULL,
    `role_type` int(11) DEFAULT 1 COMMENT 'member role 1:normal, 2:admin',
    `birthday` date DEFAULT NULL,
    `picture_url` varchar(400) DEFAULT NULL,
    `gender_type` int(11) DEFAULT 3 COMMENT 'gender type 1:male, 2:female, 3:other',
    `comment` varchar(1000) DEFAULT NULL,
    `status` int(11) DEFAULT 0 COMMENT 'status 0:in, 1:out',
    `department_id` int(11) DEFAULT 0,
    `company_id` int(11) NOT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- team_member table
**/
CREATE TABLE `team_member` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `team_id` int(11) NOT NULL,
  `member_id` int(11) NOT NULL,
  `position` INT(11) NULL DEFAULT 1 COMMENT 'position 1:member, 2:leader',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- teams table
**/
CREATE TABLE `teams` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(40) NOT NULL,
    `picture_url` varchar(400) DEFAULT NULL COMMENT 'team icon URL',
    `description` varchar(1000) DEFAULT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- companies table
**/
CREATE TABLE `companies` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(40) NOT NULL,
    `color` varchar(11) NOT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- departments table
**/
CREATE TABLE `departments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `company_id` int(11) NOT NULL,
  `color` varchar(11) NOT NULL,
  `status` int(11) DEFAULT 0 COMMENT 'department status 0:active, 1:inactive',
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- sns table
**/
CREATE TABLE `sns` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `member_id` int(11) NOT NULL,
    `github` varchar(100) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- seat_master table
**/
CREATE TABLE `seat_master` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(40) NOT NULL,
  `seat_master` varchar(1000) NOT NULL,
  `company_id` int(11) NOT NULL,
  `status` int(11) DEFAULT 1 COMMENT 'status 0:inactive, 1:active',
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- seats table
**/
CREATE TABLE `seats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL,
  `table_number` int(11) NOT NULL,
  `seat_master_id` int(11) NOT NULL,
  `row` int(11) NOT NULL,
  `col` int(11) NOT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

/**
- default company: TMH
- default member: login_id:vista, password:vista
**/
INSERT INTO `vista`.`companies` (`name`, `color`) VALUES ('TMH', '#f08080');
INSERT INTO `vista`.`members` (`login_id`, `password`, `name`, `role_type`, `picture_url`, `company_id`) VALUES ('vista', '$2y$14$PP0o64wdIsNc5zsBANqn3.9Y20wDlruA2NB.9LPPcbT07byCR1/tq', 'vista', '2', 'webroot/img/avatar_empty.png', '1');
