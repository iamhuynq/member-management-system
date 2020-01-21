/**
- members table
**/
CREATE TABLE `members` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(40) DEFAULT NULL,
    `birthday` date DEFAULT NULL,
    `picture_url` varchar(400) DEFAULT NULL,
    `gender_type` int(11) DEFAULT NULL,
    `comment` varchar(1000) DEFAULT NULL,
    `status` int(11) DEFAULT NULL,
    `team_id` int(11) DEFAULT NULL,
    `company_id` int(11) DEFAULT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

/**
- teams table
**/
CREATE TABLE `teams` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(40) DEFAULT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

/**
- companies table
**/
CREATE TABLE `companies` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(40) DEFAULT NULL,
    `color_type` int(11) DEFAULT NULL,
    `created` datetime DEFAULT CURRENT_TIMESTAMP,
    `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

/**
- sns table
**/
CREATE TABLE `sns` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `member_id` int(11) NOT NULL,
    `github` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/**
- team_member table
**/
CREATE TABLE `team_member` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `team_id` int(11) NOT NULL,
  `member_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=latin1

/**
- departments table
**/
CREATE TABLE `departments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company_id` int(11) DEFAULT NULL,
  `color` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` int(11) DEFAULT 0,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

/**
- seat_master table
**/
CREATE TABLE `seat_master` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `seat_master` varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL,
  `door_position` varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company_id` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

/**
- seats table
**/
CREATE TABLE `seats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) DEFAULT NULL,
  `table_number` int(11) DEFAULT NULL,
  `seat_master_id` int(11) DEFAULT NULL,
  `row` int(11) DEFAULT NULL,
  `column` int(11) DEFAULT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
