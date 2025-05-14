export interface GitlabInfo {
    id: number;
    api: string;
    token: string;
    webhook_url: string;
    webhook_name: string;
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
    webhook_url: string;
    webhook_name: string;
    status: 1 | -1;
    gitlab_version: string;
    expired: number | undefined;
    gitlab_url: string;
    create_time?: number;
    update_time?: number;
}

export interface GitlabCacheItem {
  projectids: string[];
  gitlabAPI: string;
  webhook_url: string;
}

export type GitlabCache = Record<string, GitlabCacheItem>;