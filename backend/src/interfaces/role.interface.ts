// 命名空间
export interface Role {
    id: number;
    namespace: string;
    role: string;
    name: string;
    describe: string;
    operator: string;
    create_time: number;
    update_time: number;
}


export interface PermissionItems {
    resource_id: number;
    describe: string;
}

export interface UpdateRoleReq {
    id: number;
    namespace: string;
    role: string;
    name: string;
    describe: string;
    permissions: PermissionItems[]
}

export interface RoleReq {
    namespace: string;
    role: string;
    name: string;
    describe: string;
    operator: string;
    create_time: number;
    update_time: number;
}

export interface RoleParams {
    id: number;
    namespace: string;
    role: string;
    name: string;
    describe: string;
    operator: string;
    create_time: number;
    update_time: number;
}

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

export interface UserRole {
    id: number;
    namespace: string;
    user: string;
    role_id: number;
    status: number;
    operator: string;
    create_time: number;
    update_time: number;
}