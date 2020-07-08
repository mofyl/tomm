
CREATE TABLE `platform_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `memo` varchar(256) NOT NULL,
  `app_key` varchar(32) NOT NULL,
  `secret_key` varchar(32) NOT NULL,
  `index_url` varchar(256) NOT NULL,
  `channel_name` varchar(32) NOT NULL,
  `sign_url` varchar(256) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  `deleted` int(1) NOT NULL ,
  `deleted_time` bigint(20),
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`),
  KEY `idx_channel_name` (`channel_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `code_infos` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_key` varchar(32) NOT NULL,
  `create_time` bigint(20) NOT NULL,
  `code` int(1) NOT NULL ,
  `mm_user_id` bigint(20),
  PRIMARY KEY (`id`),
  KEY `idx_app_key` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


