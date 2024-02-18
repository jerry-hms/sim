# 用户表
CREATE TABLE `sim_user` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
    `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
    `nickname` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
    `mobile` varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

# 聊天会话表
CREATE TABLE `sim_im_session` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) NOT NULL COMMENT '用户id',
    `rel_id` int(11) NOT NULL COMMENT '关系id',
    `session_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话名称',
    `sep_svr` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息序列号',
    `last_sender_info` json NOT NULL COMMENT '最后一条消息的发送人信息',
    `last_message` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '最后一条消息内容',
    `unread_message_number` smallint(6) NOT NULL DEFAULT '0' COMMENT '未读消息数量',
    `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天会话表';

# 聊天消息表
CREATE TABLE `sim_im_session_message` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `rel_id` int(11) NOT NULL COMMENT '会话关系id',
    `message_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息id，唯一值',
    `sep_svr` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息序列号',
    `sender_id` int(11) NOT NULL COMMENT '发送人id',
    `sender` json NOT NULL COMMENT '发送人信息',
    `receiver_id` int(11) NOT NULL COMMENT '接收者id',
    `send_content` json NOT NULL COMMENT '发送内容',
    `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读 0 = 未读，1 = 已读',
    `created_at` timestamp NULL DEFAULT NULL,
    `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表';

# 聊天会话关系表
CREATE TABLE `sim_im_session_relation` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned NOT NULL COMMENT '用户id',
    `relation_id` int(11) unsigned NOT NULL COMMENT '关系人方id',
    `scene` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话类型 friend = 好友聊天',
    `sep_svr` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '该会话的最新消息的序列号',
    `created_at` datetime NOT NULL COMMENT '创建时间',
    `updated_at` datetime NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天会话关系表';