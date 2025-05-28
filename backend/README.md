根据下面的目录结构，生成对应的文件

CREATE TABLE `t_gitlab_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `api` varchar(100) NOT NULL DEFAULT '' COMMENT 'GitLab API 名称',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT 'GitLab 访问令牌',
  `webhook_url` varchar(255) NOT NULL DEFAULT '' COMMENT 'webhook url',
  `webhook_name` varchar(255) NOT NULL DEFAULT '' COMMENT 'webhook name',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态(-1: 禁用；1: 可用)',
  `gitlab_version` varchar(20) NOT NULL DEFAULT '' COMMENT 'GitLab 版本',
  `expired` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '过期时间',
  `gitlab_url` varchar(255) NOT NULL DEFAULT '' COMMENT 'GitLab 服务器地址',
  `operator` varchar(255) NOT NULL DEFAULT '' COMMENT '操作人',
  `create_time` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
  `source_branch` varchar(255) NOT NULL DEFAULT '' COMMENT '源分支',
  `target_branch` varchar(255) NOT NULL DEFAULT '' COMMENT '目标分支',
  `project_ids` text NOT NULL COMMENT '项目ID列表',
  `project_ids_synced` tinyint(1) NOT NULL DEFAULT 0 COMMENT '项目ID列表是否已同步完成',
  PRIMARY KEY (`id`),
  KEY `api` (`api`),
  KEY `status` (`status`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='GitLab 信息表';