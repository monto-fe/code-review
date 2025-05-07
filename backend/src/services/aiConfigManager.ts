// services/AIConfigManager.ts
import DB from '../databases';
import { getUnixTimestamp } from '../utils';
import { AIConfigCreationAttributes } from '../interfaces/AIConfig.interface';

class AIConfigManager {
    private static instance: AIConfigManager;
    private config: {
        name: string;
        apiKey: string;
        apiUrl: string;
        model: string;
    } | null = null;
    public AIconfig = DB.AIConfig;
	public now:number = getUnixTimestamp();

    private refreshInterval = 10 * 1000; // 10秒刷新一次

    private constructor() {
        this.loadConfig();               // 初始加载
        this.startAutoRefresh();         // 启动定时刷新
    }

    public static getInstance(): AIConfigManager {
        if (!AIConfigManager.instance) {
            AIConfigManager.instance = new AIConfigManager();
        }
        return AIConfigManager.instance;
    }

    public async loadConfig() {
        const configRow = await DB.AIConfig.findOne({
        where: { is_active: true },
        order: [['update_time', 'DESC']]
        });
        console.log("configRow", configRow)
        if (!configRow) throw new Error('No active AI config found');
        this.config = {
        name: configRow.name,
        apiKey: configRow.api_key,
        apiUrl: configRow.api_url,
        model: configRow.model
        };
    }

    private startAutoRefresh() {
        setInterval(() => {
        this.loadConfig();
        }, this.refreshInterval);
    }
  
    public async addAIConfig(Data: AIConfigCreationAttributes): Promise<any> {
        const { name, api_url, api_key, model, is_active } = Data;
        return await this.AIconfig.create({
            name,
            api_url,
            api_key,
            model,
            is_active,
            create_time: this.now,
            update_time: this.now
        });
    }

    public async updateAIConfig(Data: AIConfigCreationAttributes): Promise<any> {
        const { id, name, api_url, api_key, model, is_active } = Data;
        return await this.AIconfig.update({
            name,
            api_url,
            api_key,
            model,
            is_active,
            update_time: this.now
        }, {
            where: {
                id
            }
        });
    }

    public async getConfigList(): Promise<any> {
        const list  = await this.AIconfig.findAll({
            order: [
                ['id', 'DESC']
            ]
        });
        return list;
    }

    public getConfig() {
        if (!this.config) {
        throw new Error('AI config not loaded yet');
        }
        return this.config;
    }
}

export default AIConfigManager.getInstance();