DROP TABLE IF EXISTS `homestay`;
CREATE TABLE `homestay` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `del_state` tinyint NOT NULL DEFAULT '0',
                            `version` bigint NOT NULL DEFAULT '0' COMMENT '版本号',
                            `title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
                            `cover` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '轮播图，第一张封面',
                            `intro` varchar(4069) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '介绍',
                            `location` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '位置',
                            `homestay_business_id` bigint NOT NULL DEFAULT '0' COMMENT '民宿店铺id',
                            `user_id` bigint NOT NULL DEFAULT '0' COMMENT '房东id，冗余字段',
                            `row_state` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:下架 1:上架',
                            `rating_stars` float NOT NULL DEFAULT 0 COMMENT '评分',
                            `price_before` bigint NOT NULL DEFAULT 0 COMMENT '民宿价格（分）',
                            `price_after`  bigint NOT NULL DEFAULT 0 COMMENT '',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='每一间民宿';