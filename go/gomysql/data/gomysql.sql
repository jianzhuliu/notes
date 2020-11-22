-- 测试表

create database gomysql default character set utf8mb4 collate utf8mb4_0900_ai_ci;

use `gomysql`;

DROP TABLE IF EXISTS `columns`;

create table `columns` (
	`id` int(11) unsigned not null auto_increment,
	`name` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
	`phone` char(11) NOT NULL DEFAULT '' COMMENT '手机号',
	`gender` enum('male','female','unknow') NOT NULL DEFAULT 'unknow' COMMENT '员工性别',
	`status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 1:enable, 0:disable, -1:deleted',
	`info` text COMMENT '描述信息',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	
	primary key (`id`)
)engine=innodb default charset=utf8mb4;

insert into `columns` (`name`, `phone`, `gender`, `status`, `info`) values
("name1", "12345678912", "male", 1, "information"),
("name2", "12345678913", "female", 0, "information"),
("name3", "12345678914", "female", 1, "information"),
("name4", "12345678915", "male", 1, "information");
