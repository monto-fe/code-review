USE ucode_review;
-- t_namespace
INSERT INTO `ucode_review`.`t_namespace` (`id`, `namespace`, `parent`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (1, 'acl', '', '系统权限', '权限管理分类', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_resource
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (1, 'acl', 'API', 'CreateResource', '{}', '创建资源', '创建资源', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (2, 'acl', 'Data', 'GetUserList', '{}', '用户列表', '用户列表', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (3, 'acl', 'API', 'GetResources', '{}', '获取资源', '获取资源信息', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (4, 'acl', 'API', 'UpdateResource', '{}', '更新资源', '更新资源接口', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (5, 'acl', 'API', 'GetRoles', '{}', '获取角色', '获取角色', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (6, 'acl', 'API', 'AssertRolePermission', '{}', 'API角色分配资源权限', 'API角色分配资源权限', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_resource` (`id`, `namespace`, `category`, `resource`, `properties`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (7, 'acl', 'API', 'GetUserInfo', '{}', '获取用户信息', '获取用户信息', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_role
INSERT INTO `ucode_review`.`t_role` (`id`, `namespace`, `role`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (1, 'acl', 'admin', '超级管理员', '最高权限', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_role` (`id`, `namespace`, `role`, `name`, `describe`, `operator`, `create_time`, `update_time`) VALUES (2, 'acl', 'salePeople', '销售', '销售角色', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_role_permission
INSERT INTO `ucode_review`.`t_role_permission` (`id`, `namespace`, `role_id`, `resource_id`, `describe`, `operator`, `create_time`, `update_time`) VALUES (1, 'acl', 1, 1, '创建资源', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
INSERT INTO `ucode_review`.`t_role_permission` (`id`, `namespace`, `role_id`, `resource_id`, `describe`, `operator`, `create_time`, `update_time`) VALUES (2, 'acl', 1, 6, '获取角色', 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_user
INSERT INTO `ucode_review`.`t_user` (`id`, `namespace`, `o_id`, `user`, `name`, `password`, `job`, `phone_number`, `email`, `create_time`, `update_time`) VALUES (1, 'acl', NULL, 'admin', '超级管理员', '25d55ad283aa400af464c76d713c07ad', 'admin', '', 'admin@admin.com', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_user_role
INSERT INTO `ucode_review`.`t_user_role` (`id`, `namespace`, `user`, `role_id`, `status`, `operator`, `create_time`, `update_time`) VALUES (1, 'acl', 'admin', 1, 1, 'admin', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());
-- t_ai_manager
INSERT INTO `ucode_review`.`t_ai_manager` (`type`, `model`, `api_url`, `status`, `create_time`, `update_time`)
VALUES 
  ('UCloud', 'deepseek-ai/DeepSeek-V3-0324', 'https://deepseek.modelverse.cn/v1/chat/completions', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  ('DeepSeek', 'deepseek-chat', 'https://api.deepseek.com/chat/completions', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());


