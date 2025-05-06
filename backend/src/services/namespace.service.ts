import { Op } from 'sequelize';
import DB from '../databases';
import { NamespaceReq, NamespaceParams, Namespace } from '../interfaces/namespace.interface';
import { getUnixTimestamp } from '../utils';

class NamespaceService {
	public Namespace = DB.Namespace;
	public now:number = getUnixTimestamp();

	public async create(Data: NamespaceReq): Promise<any> {
		// const { namespace, parent, name, describe, operator } = Data;
		const data: NamespaceParams = {
			...Data,
			create_time: this.now,
			update_time: this.now
		}
		const res: any = await this.Namespace.create({ ...data });
		return res;
	}

	public async update(Data : Namespace): Promise<[number]> {
		const { id, describe } = Data;
		const params:Namespace = {
			...Data
		}

		if (describe) {
			params.describe = describe
		}

		// 无更新项目
		if (Object.keys(params).length === 0) {
			return [0]
		}

		params.update_time = this.now
		const res: [number] = await this.Namespace.update({ ...params }, {
			where: { id }
		})
		return res;
	}
	
	public async deleteSelf(Data: any): Promise<any> {
		const { id } = Data;
		const res:any = await this.Namespace.destroy({
			where: { id }
		})
		return res;
	}

	/**
	 * 获取项目组及其子项目组信息
	*/
	public async findWithAllChildren(Data: any): Promise<any> {
		const { namespace, offset, limit } = Data;
		let queryOptions:any = {
			offset,
			limit,
			order: [['id', 'DESC']]
		};
		
		if (namespace) {
			queryOptions.where = {
				namespace
			};
		}
		const result: any = await this.Namespace.findAndCountAll(queryOptions)
		return result;
	}

	/* 校验项目组是否存在
	 */
	public async checkNamespace({
		namespace,
		name
	}:{
		namespace: string,
		name: string
	}): Promise<any> {
		const result: any = await this.Namespace.findOne({
			where: {
				[Op.or]: [{ namespace }, { name }],
			}
		})
		return result ? true : false;
	}
}

export default NamespaceService;