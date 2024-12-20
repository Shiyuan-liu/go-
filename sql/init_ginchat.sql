CREATE Table `community`(
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name` longtext,
    `ownerid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `img` longtext,
    `desc` longtext,
    PRIMARY KEY(`id`),
    KEY `idx_community_deleted_at` (`deleted_at`)
)ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;

CREATE Table `contact`(
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `ownerid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `targetid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `type` BIGINT(20) DEFAULT NULL,
    `desc` longtext,
    PRIMARY KEY (`id`),
    KEY `idx_contact_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=185 DEFAULT CHARSET=utf8;

CREATE Table `group_basic`(
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name` longtext,
    `ownerid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `icon` longtext,
    `type` BIGINT(20) DEFAULT NULL,
    `desc` longtext,
    PRIMARY KEY (`id`),
    KEY `idx_group_basic_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE Table `message`(
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `userid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `targetid` BIGINT(20) UNSIGNED DEFAULT NULL,
    `type` BIGINT(20) DEFAULT NULL,
    `media` BIGINT(20) DEFAULT NULL,
    `content` longtext,
    `createTime` bigint(20) unsigned NULL AUTO_INCREMENT,
    `ReadTime` bigint(20) unsigned NULL AUTO_INCREMENT,
    `pic` longtext,
    `url` longtext,
    `desc` longtext,
    `amount` bigint(20) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_message_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE Table `user_basic`(
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `name` longtext,
    `password` longtext,
    `phone` longtext,
    `email` longtext,
    `identity` longtext,
    `client_ip` longtext,
    `client_port` longtext,
    `login_time` datetime(3) DEFAULT NULL,
    `heart_beat_time` datetime(3) DEFAULT NULL,
    `login_out_time` datetime(3) DEFAULT NULL,
    `is_log_out` tinyint(1) DEFAULT NULL,
    `device_info` longtext,
    `salt` longtext,
    `avatar` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_basic_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;