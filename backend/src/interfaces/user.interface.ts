import { PageModel } from './common.interface';
export enum UserStatus {
    Init,
    Enable,
    Disable
}

export interface User {
    id: number; // Note that the `null assertion` `!` is required in strict mode.
    o_id: number;
    namespace: string;
    user: string;
    name: string;
    job: string;
    password: string;
    phone_number: string;
    email: string;
    create_time: number;
    update_time: number;
}

export interface UserAndRoles {
    User: User;
    Roles: any[];
}

export interface UserResponse {
    data: User[];
    count: number;
}

export interface UserLogin {
    namespace: string;
    user: string;
    password: string;
}

export interface LoginParams {
    id?: number;
    user: string, 
    namespace: string
}

export interface UserQuery {
    id: number;
    username: string;
    password: string;
}

// 注册code码，保存邮箱、code、时间、是否已失效、type(注册、修改密码)
export interface Token {
    id: number;
    token: string;
    user: string;
    expired_at: number;
    create_time: number;
    update_time: number;
}


export interface UserListReq extends PageModel {
    id?: number;
    user?: string[];
    namespace?: string;
    name?: string;
    roleName?: string;
}

export interface UserReq {
    id?: number; // Note that the `null assertion` `!` is required in strict mode.
    user: string;
    namespace: string;
    name: string;
    job: string;
    password: string;
    phone_number?: string;
    email?: string;
    role_ids: number[];
    operator: string;
}

export interface UpdateUserReq {
    id: number;
    user: string;
    namespace: string;
    name: string;
    job: string;
    password: string;
    phone_number?: string;
    email?: string;
    role_ids: number[];
    operator: string;
}

export interface UserListQuery {
    id: number;
    user: string;
    userName: string;
    roleName: string;
    current: number;
    pageSize: number;
    namespace: string;
}

export interface UserRole {
    id: number;
    namespace: string;
    user: string;
    name: string;
    job: string;
    phone_number: string;
    email: string;
    role_id: number;
}