export interface RolePermission {
  id: number;
  namespace: string;
  role_id: number;
  resource_id: number;
  describe: string;
  operator: string;
  create_time: number;
  update_time: number;
}

export interface RolePermissionReq {
  namespace: string;
  role_id: number;
  resource_ids: number[];
  describe: string;
  operator: string;
}

export interface RoleSinglePermissionReq {
  namespace: string;
  role_id: number;
  resource_id: number;
  describe: string;
  operator: string;
}

export interface UserRoleReq {
  namespace: string;
  user: number;
  role_id?: number;
  role_ids: number[];
  describe: string;
  operator: string;
  create_time: number;
  update_time: number;
}

export interface CancelRolePermissionReq {
  namespace: string;
  role_id: number;
  resource_id: number;
  operator: string;
}

export interface CancelUserRoleReq {
  namespace: string;
  user: number;
  role_id: number;
  role_ids: number[];
  operator: string;
}

export interface GetRolePermissionsReq {
  namespace: string;
  role_id: number;
}

export interface GetSelfPermissionsReq {
  namespace: string;
  category: string;
}

export interface GetUserRoles {
  namespace: string,
  user: string,
  role_id: 0,
  role_ids: number[],
  userId: 0,
  userIds: number[]
}

export interface PermissionItems {
  resource_id: number;
  describe: string;
}

export interface UpdateRolePermissionsReq {
  namespace: string;
  role_id: number;
  permissions: PermissionItems[];
}

export interface CheckPermissionReq {
  namespace: string;
  user: string;
  resource: string;
  properties: object;
}