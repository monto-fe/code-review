-- 选择数据库
USE ucode_review;

-- 创建 t_namespace 表
CREATE TABLE IF NOT EXISTS `t_namespace` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace` varchar(128) NOT NULL DEFAULT '' COMMENT '项目组标识',
  `parent` varchar(128) NOT NULL DEFAULT '' COMMENT '所属父项目组',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '项目组中文名',
  `describe` varchar(128) NOT NULL DEFAULT '' COMMENT '描述',
  `operator` varchar(128) NOT NULL DEFAULT '' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `namespace` (`namespace`),
  KEY `parent` (`parent`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_login_token 表
CREATE TABLE IF NOT EXISTS `t_login_token` (
  `id` int NOT NULL AUTO_INCREMENT,
  `token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `user` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `expired_at` int NOT NULL,
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_resource 表
CREATE TABLE IF NOT EXISTS `t_resource` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace` varchar(128) NOT NULL DEFAULT '' COMMENT '所属项目组',
  `category` varchar(128) NOT NULL DEFAULT '' COMMENT '资源分类',
  `resource` varchar(128) NOT NULL DEFAULT '' COMMENT '资源标识',
  `properties` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资源属性 (JSON Schema)',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '资源标识',
  `describe` text NOT NULL COMMENT '描述',
  `operator` varchar(128) NOT NULL DEFAULT '' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `unique_namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_role 表
CREATE TABLE IF NOT EXISTS `t_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace` varchar(128) NOT NULL DEFAULT '' COMMENT '所属项目组',
  `role` varchar(255) NOT NULL DEFAULT '' COMMENT '角色标识',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '角色中文名',
  `describe` varchar(512) NOT NULL DEFAULT '' COMMENT '描述',
  `operator` varchar(512) NOT NULL DEFAULT '' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique-namespace_name` (`namespace`,`name`),
  UNIQUE KEY `unique-namespace_role` (`namespace`,`role`),
  KEY `unique-namespace_id` (`namespace`,`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_role_permission 表
CREATE TABLE IF NOT EXISTS `t_role_permission` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace` varchar(128) NOT NULL DEFAULT '' COMMENT '所属项目组',
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  `resource_id` int unsigned NOT NULL COMMENT '资源ID',
  `describe` text NOT NULL COMMENT '描述',
  `operator` varchar(128) NOT NULL DEFAULT '' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`namespace`,`role_id`,`resource_id`),
  KEY `role_id` (`role_id`),
  KEY `resource_id` (`resource_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_user 表
CREATE TABLE IF NOT EXISTS `t_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'user id',
  `namespace` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `o_id` int DEFAULT NULL COMMENT 'oa id',
  `user` varchar(255) NOT NULL COMMENT '用户英文名',
  `name` varchar(255) NOT NULL COMMENT '用户中文名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `job` varchar(512) NOT NULL DEFAULT '' COMMENT '职位名称',
  `phone_number` varchar(512) NOT NULL DEFAULT '' COMMENT '手机号码',
  `email` varchar(512) NOT NULL DEFAULT '' COMMENT '邮箱',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `delete_time` int DEFAULT NULL,
  PRIMARY KEY (`id`),
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- 创建 t_user_role 表
CREATE TABLE IF NOT EXISTS `t_user_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `namespace` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '所属项目组',
  `user` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户英文名',
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '当前状态(-1: 禁用；1: 可用；)',
  `operator` varchar(512) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '创建人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `namespace` (`namespace`),
  KEY `user` (`user`),
  KEY `role_id` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- gitlab table --
CREATE TABLE IF NOT EXISTS `t_gitlab_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `api` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'GitLab API 名称',
  `token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'GitLab 访问令牌',
  `webhook_url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'webhook url',
  `webhook_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'webhook name',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态(-1: 禁用；1: 可用)',
  `source_branch` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '源分支',
  `target_branch` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '目标分支',
  `gitlab_version` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'GitLab 版本',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `gitlab_url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'GitLab 服务器地址',
  `prompt` text COLLATE utf8mb4_general_ci COMMENT 'AI prompt信息',
  `webhook_status` tinyint NOT NULL DEFAULT '0' COMMENT 'webhook状态(0: 未设置, 1: 已设置)',
  `project_ids` text COLLATE utf8mb4_general_ci COMMENT '项目ID列表，逗号分隔',
  `project_ids_synced` tinyint(1) NOT NULL DEFAULT '0' COMMENT '项目ID是否已同步(0: 未同步, 1: 已同步)',
  `rule_check_status` tinyint NOT NULL DEFAULT '0' COMMENT '规则检查状态(0: 未检查, 1: 检查中, 2: 已检查)',
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `api` (`api`),
  KEY `status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='GitLab 信息表';


-- 表3: t_common_rule
CREATE TABLE IF NOT EXISTS `t_common_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规则名称',
  `language` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '语言类型 (Python, Java, JavaScript, Go, Ruby, C++, other)',
  `rule` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '规则内容',
  `description` text COLLATE utf8mb4_general_ci COMMENT '规则描述',
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `language` (`language`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='通用规则表';


-- 表4: t_custom_rule
CREATE TABLE IF NOT EXISTS `t_custom_rule` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `project_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '项目名称',
  `project_id` int unsigned NOT NULL COMMENT '项目 ID',
  `rule` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '项目自定义规则内容',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '规则状态 (1: enable, -1: disabled)',
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='自定义规则表';


-- 表5: t_ai_message
CREATE TABLE IF NOT EXISTS `t_ai_message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `project_id` varchar(255) NOT NULL COMMENT '项目 ID',
  `merge_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合并请求的 ID',
  `merge_url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合并请求的 URL',
  `ai_model` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'AI 模型名称',
  `rule` tinyint NOT NULL COMMENT '规则类型 (1: common, 2: custom)',
  `rule_id` int unsigned NOT NULL COMMENT '规则 ID',
  `result` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '检测结果',
  `passed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否通过检测 (0: 否, 1: 是)',
  `human_rating` tinyint NOT NULL DEFAULT '0' COMMENT '人工评分',
  `remark` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '备注',
  `checked_by` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '检测员（为空表示自动检测）',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作人',
  PRIMARY KEY (`id`),
  KEY `project_id` (`project_id`),
  KEY `rule_id` (`rule_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='AI 检测消息表';

-- t_ai_config
CREATE TABLE IF NOT EXISTS `t_ai_config` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(64) NOT NULL COMMENT '配置名称',
  `api_url` VARCHAR(255) NOT NULL COMMENT 'API 请求地址',
  `api_key` VARCHAR(255) NOT NULL COMMENT 'API 密钥',
  `model` VARCHAR(100) NOT NULL COMMENT 'AI 模型名称',
  `is_active` TINYINT(1) NOT NULL DEFAULT 2 COMMENT '是否启用该配置（1: 启用，2: 禁用）',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_model` (`model`),
  KEY `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI 配置表';

-- 为优化查询性能添加索引
CREATE INDEX idx_project_id ON t_ai_message(project_id);
CREATE INDEX idx_rule_id ON t_ai_message(rule_id);
CREATE INDEX idx_status ON t_gitlab_info(status);