// import { Op } from 'sequelize';
import DB from '../databases';
import { getUnixTimestamp } from '../utils';
import { GitlabInfoCreate } from '../interfaces/gitlab.interface';

const ENABLE = 1

class GitlabService {
    public GitlabInfo = DB.GitlabInfo;
    
    public now:number = getUnixTimestamp();

    // 添加gitlab的信息
    public async addGitlabToken(Data: GitlabInfoCreate): Promise<any> {
        const { api, token, status=1, gitlab_version, expired, gitlab_url } = Data;
        const data = {
            api,
            token,
            status,
            gitlab_version,
            expired,
            gitlab_url,
            create_time: this.now,
            update_time: this.now
        }
        const res: any = await this.GitlabInfo.create({ ...data });
        return res;
    }
    // 更新gitlab的信息
    public async updateGitlabInfo(Data: any): Promise<any> {
        const { id, api, token, status, gitlab_version, expired, gitlab_url } = Data;
        const data = {
            api,
            token,
            status,
            gitlab_version,
            expired,
            gitlab_url,
            update_time: this.now
        }
        const response: any = await this.GitlabInfo.update({ ...data }, {
            where: {
                id
            }
        });
        return response;
    }
    // 获取gitlab的信息
    public async getGitlabInfo() {
        const response = await this.GitlabInfo.findAll({
            where: {
                status: ENABLE
            }
        })
        return response
    }
}

export default GitlabService;