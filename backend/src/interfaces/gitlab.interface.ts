export interface GitlabInfo {
    id: number;
    api: string;
    token: string;
    status: 1 | -1;
    gitlab_version: string;
    expired: number;
    gitlab_url: string;
    create_time: number;
    update_time: number;
}

export interface GitlabInfoCreate {
    api: string;
    token: string;
    status: 1 | -1;
    gitlab_version: string;
    expired: number | undefined;
    gitlab_url: string;
    create_time?: number;
    update_time?: number;
}