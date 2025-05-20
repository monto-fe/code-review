import Sequelize from 'sequelize';
import { Op } from 'sequelize';
import DB from '../databases';
import { getUnixTimestamp } from '../utils';
import { CommonRuleCreate } from '../interfaces/common.interface';

class AIRuleService {
	public CommonRule = DB.CommonRule;
    public CustomRule = DB.CustomRule;
	public now:number = getUnixTimestamp();

    public async getCommonRule({ language }: { language: string}): Promise<any> {
        let query: any = {};
        if(language){
            query = {
                where: {
                    [Op.and]: [
                        Sequelize.where(Sequelize.fn('LOWER', Sequelize.col('language')), Op.eq, language.toLowerCase())
                    ]
                }
            }
        }
        const rule = await this.CommonRule.findAll(query)
        console.log('rule', rule);
        /* 循环rule，pick出language === language的规则，language全部转小写对比*/
        if(!rule.length){
            return [];
        }
        const pickRule = rule.filter((item:any) => {
            return item.language.toLowerCase === item.language.toLowerCase();
        })
		return pickRule;
	}

    // 创建通用规则
    public async createCommonRule(Data: CommonRuleCreate): Promise<any> {
        const { name, language, rule, description } = Data;
        const result = await this.CommonRule.create({
            name,
            language,
            rule,
            description,
            create_time: this.now,
            update_time: this.now
        })
        return result;
    }

    // 更新
    public async updateCommonRule(Data: any): Promise<any> {
        const { id, name, language, rule, description } = Data;
        const result = await this.CommonRule.update({
            name,
            language,
            rule,
            description,
            update_time: this.now
        }, {
            where: {
                id
            }
        })
        return result;
    }

    // 创建自定义规则
    public async createCustomRule(Data: any): Promise<any> {
        const { project_name, project_id, rule, status=1, operator } = Data;
        const result = await this.CustomRule.create({
            project_name,
            project_id,
            rule,
            status,
            operator,
            create_time: this.now,
            update_time: this.now
        })
        return result;
    }

    // 更新自定义规则
    public async updateCustomRule(Data: any): Promise<any> {
        const { id, project_name, project_id, rule, status=1, operator } = Data;
        const result = await this.CustomRule.update({
            project_name,
            project_id,
            rule,
            status,
            operator,
            update_time: this.now
        }, {
            where: {
                id
            }
        })
        return result;
    }

    // 获取自定义规则
    public async getCustomRuleByProjectId({ project_id }: { project_id: string}): Promise<any> {
		// const { namespace, parent, name, describe, operator } = Data;
        // 先通过projectid获取对应的规则
        let params:any = {
            status: 1
        }
        if(project_id){
            params.project_id = project_id;
        }
        // 如果没有则通过language获取通用规则
        const result = await this.CustomRule.findAndCountAll({
            where: params
        })
        return result;
    }

    // 更新human_rating和remark字段
    public async updateHumanRatingAndRemark(Data: any): Promise<any> {
        const { id, human_rating, remark } = Data;
        const result = await this.CustomRule.update({
            human_rating, remark
        }, {
            where: {
                id
            }
        })
        return result;
    }
}

export default AIRuleService;