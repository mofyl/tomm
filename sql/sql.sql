CREATE DATABASE platform DEFAULT CHARSET utf8;

CREATE TABLE `platform`.`platform_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `memo` varchar(256) NOT NULL COMMENT '简介',
  `app_key` varchar(32) NOT NULL COMMENT '唯一标识',
  `secret_key` varchar(32) NOT NULL COMMENT '私钥',
  `index_url` varchar(256) NOT NULL COMMENT '主页url',
  `channel_name` varchar(32) NOT NULL COMMENT '平台名字',
  `sign_url` varchar(256) NOT NULL COMMENT '回调url',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `deleted` int(1) NOT NULL COMMENT '删除标记 1 表示未删除  2表示删除',
  `deleted_time` bigint(20) COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`),
  KEY `idx_channel_name` (`channel_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `platform`.`mm_user_auth_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_key` varchar(32) NOT NULL COMMENT'三方平台唯一标识',
  `create_time` bigint(20) NOT NULL COMMENT '创建时间',
  `mm_user_id` varchar(256) NOT NULL COMMENT 'mm授权用户的ID',
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



CREATE TABLE `platform`.`admin_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `login_name` varchar(32) NOT NULL COMMENT '三方平台管理员的登录名',
  `login_pwd` varchar(32) NOT NULL COMMENT '三方平台管理员的密码  这里会使用base64编码',
  `name` varchar(32) NOT NULL COMMENT '用户名',
  `number` varchar(32) NOT NULL COMMENT '工号',
  `created` bigint(32) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  key idx_login_name(`login_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `platform_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(32) NOT NULL COMMENT '角色的名字',
  `platform_ids` varchar(256) NOT NULL COMMENT '该角色可访问的平台ID，若有多个使用逗号分隔',
  `create_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;








