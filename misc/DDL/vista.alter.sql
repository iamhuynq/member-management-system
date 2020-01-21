-- ********************
-- SCHEMA
-- ********************
-- 14/06/2019
ALTER SCHEMA `vista`  DEFAULT CHARACTER SET utf8  DEFAULT COLLATE utf8_unicode_ci;


-- ********************
-- members
-- ********************
-- 10/06/2019
ALTER TABLE `vista`.`members` ADD COLUMN `login_id` VARCHAR(100) NOT NULL COMMENT 'login ID' AFTER `name`;
ALTER TABLE `vista`.`members` ADD COLUMN `password` VARCHAR(200) NOT NULL COMMENT 'password' AFTER `login_id`;

-- 14/06/2019
ALTER TABLE `vista`.`members` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;

-- 19/06/2019
ALTER TABLE `vista`.`members` ADD COLUMN `role_type` INT(20) NOT NULL DEFAULT 1 COMMENT 'member role 1:normal, 2:admin' AFTER `name`;

-- 09/07/2019
ALTER TABLE `vista`.`members`
DROP COLUMN `team_id`;

-- 19/08/2019
ALTER TABLE `vista`.`members` ADD COLUMN `department_id` INT(11) NOT NULL COMMENT 'department ID' AFTER `status`;

-- 21/08/2019
ALTER TABLE `vista`.`members`
CHANGE COLUMN `name` `name` VARCHAR(40) NOT NULL,
CHANGE COLUMN `gender_type` `gender_type` INT(11) NOT NULL DEFAULT 3 COMMENT 'gender type 1:male, 2:female, 3:other',
CHANGE COLUMN `status` `status` INT(11) NOT NULL DEFAULT 0 COMMENT 'status 0:enrolling, 1:quited',
CHANGE COLUMN `company_id` `company_id` INT(11) NOT NULL;

-- 11/09/2019
ALTER TABLE `vista`.`members`
CHANGE COLUMN `department_id` `department_id` INT(11) NOT NULL DEFAULT 0 ;

-- 13/09/2019
ALTER TABLE `vista`.`members`
CHANGE COLUMN `birthday` `birthday` date DEFAULT NULL;


-- ********************
-- team_member
-- ********************
-- 15/07/2019
ALTER TABLE `vista`.`team_member` ADD COLUMN `position` INT(11) NULL DEFAULT 1 AFTER `member_id`;


-- ********************
-- teams
-- ********************
-- 07/06/2019
ALTER TABLE `vista`.`teams` ADD COLUMN `picture_url` VARCHAR(400) COMMENT 'team iconURL' AFTER `name`;
ALTER TABLE `vista`.`teams` ADD COLUMN `description` VARCHAR(1000) COMMENT 'team description' AFTER `picture_url`;

-- 14/06/2019
ALTER TABLE `vista`.`teams` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;

-- 21/08/2019
ALTER TABLE `vista`.`teams` CHANGE COLUMN `name` `name` VARCHAR(40) NOT NULL;


-- ********************
-- companies
-- ********************
-- 14/06/2019
ALTER TABLE `vista`.`companies` CONVERT TO CHARACTER SET utf8 COLLATE utf8_unicode_ci;

-- 17/07/2019
ALTER TABLE `vista`.`companies` CHANGE COLUMN `color_type` `color` VARCHAR(20) NULL DEFAULT NULL ;

-- 21/08/2019
ALTER TABLE `vista`.`companies`
CHANGE COLUMN `name` `name` VARCHAR(40) NOT NULL,
CHANGE COLUMN `color` `color` VARCHAR(11) NOT NULL ;


-- ********************
-- departments
-- ********************
-- 21/08/2019
ALTER TABLE `vista`.`departments`
CHANGE COLUMN `name` `name` VARCHAR(40) NOT NULL,
CHANGE COLUMN `company_id` `company_id` INT(11) NOT NULL,
CHANGE COLUMN `color` `color` VARCHAR(11) NOT NULL ;


-- ********************
-- seat_master
-- ********************
-- 01/08/2019
ALTER TABLE `vista`.`seat_master` DROP COLUMN `door_position`;

-- 19/08/2019
ALTER TABLE `vista`.`seat_master`
CHANGE COLUMN `title` `title` VARCHAR(40) NOT NULL ;

-- 21/08/2019
ALTER TABLE `vista`.`seat_master`
CHANGE COLUMN `seat_master` `seat_master` varchar(1000) NOT NULL,
CHANGE COLUMN `company_id` `company_id` INT(11) NOT NULL,
CHANGE COLUMN `status` `status` INT(11) NOT NULL DEFAULT 1 COMMENT 'status 0:inactive, 1:active';


-- ********************
-- seats
-- ********************
-- 21/08/2019
ALTER TABLE `vista`.`seats`
CHANGE COLUMN `member_id` `member_id` INT(11) NOT NULL,
CHANGE COLUMN `seat_master_id` `seat_master_id` INT(11) NOT NULL;

-- 17/09/2019
ALTER TABLE `vista`.`seats`
CHANGE COLUMN `table_number` `table_number` INT(11) NOT NULL ,
CHANGE COLUMN `row` `row` INT(11) NOT NULL ,
CHANGE COLUMN `column` `col` INT(11) NOT NULL ;
