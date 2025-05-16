// src/services/gitlabManager.service.ts
import axios from 'axios';
import GitlabService from './gitlab.service';
import { GitlabCache, GitlabCacheItem } from '../interfaces/gitlab.interface';
//import DB from '../databases';


export class GitlabManagerService {
  private static instance: GitlabManagerService;
  private cache: GitlabCache = {};
  public GitlabInfo = new GitlabService();

  private constructor() {}

  // 单例初始化
  static async init(): Promise<GitlabManagerService> {
    if (!this.instance) {
      this.instance = new GitlabManagerService();
      await this.instance.loadAllTokens();
    }
    return this.instance;
  }

  // 加载所有 Token 对应的项目 ID
  private async loadAllTokens() {
    const tokens = await this.GitlabInfo.getGitlabData();

    for (const t of tokens) {
      await this.loadToken({
        token: t.token,
        gitlabAPI: t.api,
        webhook_url: t.webhook_url,
      });
    }
  }

  // 加载单个 Token 项目ID并缓存
  private async loadToken(config: {
    token: string;
    gitlabAPI: string;
    webhook_url: string;
  }): Promise<void> {
    const { token, gitlabAPI, webhook_url } = config;

    try {
      const projects: string[] = await this.fetchProjectIds(token, gitlabAPI);

      this.cache[token] = {
        projectids: projects,
        gitlabAPI,
        webhook_url,
      };
      console.log(`GitLab Token 加载成功, 项目数: ${projects.length}`);
    } catch (err) {
      console.error(`加载 GitLab Token 失败: `, err);
    }
  }

  // 通过 GitLab API 获取项目 ID 列表
  private async fetchProjectIds(token: string, gitlabAPI: string): Promise<string[]> {
    const projectIds: string[] = [];
    let page = 1;
    const perPage = 100;

    while (true) {
      const res = await axios.get(`${gitlabAPI}/v4/projects`, {
        headers: { 'PRIVATE-TOKEN': token },
        params: { per_page: perPage, page },
      });
      if (Array.isArray(res.data) && res.data.length) {
        const ids = res.data.map((proj: any) => String(proj.id));
        projectIds.push(...ids);
        if (res.data.length < perPage) break;
        page += 1;
      } else {
        break;
      }
    }

    return projectIds;
  }

  // 提供全局缓存获取方法
  getCache(): GitlabCache {
    return this.cache;
  }

}