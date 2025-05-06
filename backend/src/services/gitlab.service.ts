// import { Op } from 'sequelize';
import OpenAI from "openai";
import DB from '../databases';
import { getUnixTimestamp } from '../utils';
import { GitlabInfoCreate } from '../interfaces/gitlab.interface';

const ENABLE = 1

class gitlabService {
    public GitlabInfo = DB.GitlabInfo;
    
    public now:number = getUnixTimestamp();
    public openai = new OpenAI({
        baseURL: "",
        apiKey: "",
    });

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

export default gitlabService;