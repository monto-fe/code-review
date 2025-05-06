import Sequelize from 'sequelize';
import { Op } from 'sequelize';
import DB from '../databases';
import { RoleReq, Role, UpdateRoleReq } from '../interfaces/role.interface';
import { getUnixTimestamp } from '../utils';

class RoleService {
	public Role = DB.Role;
	public RolePermission = DB.RolePermission;
	public Resource = DB.Resource;
	public now:number = getUnixTimestamp();

	public async create(Data: RoleReq): Promise<any> {
		const data: RoleReq = {
			...Data,
			create_time: this.now,
			update_time: this.now
		}
		const res: any = await this.Role.create({ ...data });
		return res;
	}

	public async createRolePermission(Data: any): Promise<any> {
		const { namespace, role, name, describe, operator, permissions } = Data;
		const t = await DB.sequelize.transaction();
		try {
			// 新增角色
			const roleParams:any = {
				namespace,
				role,
				name,
    			describe,
    			operator,
				create_time: this.now,
    			update_time: this.now
			}
			const roleInfo: any = await this.Role.create({ ...roleParams }, { transaction: t });
			// 角色关联资源
			const permissionParams:any = permissions
			await this.RolePermission.bulkCreate(permissionParams.map((item: any) => {
				return {
					namespace,
					role_id: roleInfo.id,
					resource_id: item.resource_id,
					describe: item.describe,
					operator,
					create_time: this.now,
					update_time: this.now
				}
			}), { transaction: t });
			t.commit();
			return roleInfo;
		} catch (err) {
			await t.rollback();
			throw err;
		}
	}

	public async update(Data : any): Promise<[number]> {
		const { id, namespace, role, permissions, name, describe, operator } = Data;
		const params:any = {
			update_time: this.now
		}
		if (role) {
			params.namespace = namespace
		}
		if (describe) {
			params.describe = describe
		}
		if (name) {
			params.name = name
		}
		const t = await DB.sequelize.transaction();
		let response = [];
		try {
			// 更新角色
			response = await this.Role.update({ ...params }, {
				where: { id },
				transaction: t
			})
			// 当前角色有的资源
			const result = await this.RolePermission.findAll({
				where: { role_id: id, namespace },
				transaction: t
			})
			const currentResourceIds = result.map((item:any) => {
				return item.dataValues.resource_id
			})
			// 新的资源
			const newResourceIds = permissions.map((item:any) => item.resource_id)
			// 两者取交集，即为要更新的资源
			const updateResourceIds = newResourceIds.filter((item:any) => currentResourceIds.includes(item))
			const rolePermission = updateResourceIds.map((item:any) => {
				const describe = permissions.find((resource:any) => resource.resource_id === item).describe
				return {
					namespace,
					role_id: id,
					resource_id: item,
					describe
				}
			})
			await this.RolePermission.bulkCreate(rolePermission, {
				updateOnDuplicate: ['resource_id', 'describe'],
				transaction: t
			})
			// currentResourceIds不在newResourceIds中的，要删除的资源
			const deleteResources = currentResourceIds.filter((item:number) => !newResourceIds.includes(item))
			await this.RolePermission.destroy({
				where: {
					namespace,
					role_id: id,
					resource_id: {
						[Op.in]: deleteResources
					}
				},
				transaction: t
			})
			// newResourceIds不在currentResourceIds中的，要添加的资源
			const newResources = newResourceIds.filter((item:any) => !currentResourceIds.includes(item)).map((item:any) => {
				return {
					namespace,
					role_id: id,
					resource_id: item,
					describe: permissions.find((resource:any) => resource.resource_id === item).describe,
					operator: operator || '',
					created_at: getUnixTimestamp(),
					updated_at: getUnixTimestamp()
				}
			})
			await this.RolePermission.bulkCreate(newResources, {
				transaction: t
			})
			await t.commit();
		} catch (error) {
			await t.rollback();
			throw error;
		}
		return response;
	}
	
	public async deleteSelf(Data: any): Promise<any> {
		const { namespace, id } = Data;
		const res:any = await this.Role.destroy({
			where: { id, namespace }
		})
		return res;
	}

	/**
	 * 获取全部角色信息
	*/
	public async findWithAllChildren(Data: any): Promise<any> {
		const { namespace, role, name, resource, limit, offset } = Data;
		const params:any = {
			namespace
		}
		if (role) {
			params.role = {
				[Op.like]: `%${role}%`
			};
		}
		if (name) {
			params.name = {
				[Op.like]: `%${name}%`
			}
		}
		// 如果resource存在，则根据resource和namespace获取资源id列表
		const resource_ids = resource ? await this.Resource.findAll({
			where: {
				namespace,
				resource: {
					[Op.in]: resource
				}
			}
		}).then((res: any) => res.map((item: any) => item.id)) : [];

		const result: any = await this.Role.findAndCountAll({ 
			where: params,
			limit,
			offset,
			order: [['id', 'DESC']]
		})
		let query: any = {
			role_id: {
				[Op.in]: result.rows.map((item: any) => item.id)
			}
		}
		if (resource_ids.length){
			query.resource_id = {
				[Op.in]: resource_ids
			}
		}
		// 查询role_id对应的role_permission
		const rolePermission: any = await this.RolePermission.findAll({
			where: query
		})
		result.rows.forEach((item: any) => {
			const rolePermissionItem = rolePermission.filter((rolePermissionItem: any) => rolePermissionItem.dataValues.role_id === item.dataValues.id)
			item.dataValues.resources = rolePermissionItem.map((rolePermissionItem: any) => rolePermissionItem.dataValues.resource_id)
		})
		return result;
	}

	/*
		现在有Resource表、RoleResource表，通过RoleResource表中的resourceId关联Resource表， RoleResource表中有role_id
		现在有参数role_id数组，我要查询所有role_id关联的Resource表中的数据
		请写出语句
	*/
	public async findResourceByRoleIds(roleIds: number[]) {
		const roleIdsStr = roleIds.join(',');
		const result: any = await DB.sequelize.query(`SELECT
		resource.* FROM t_resource AS resource JOIN t_role_permission WHERE
		resource.id = t_role_permission.resource_id
		and t_role_permission.role_id in (${roleIdsStr})`)
		return result;
	}
	// 通过user查询对应的角色信息
	public async findRoleByUser(user: string[]) {
		const result: any = await DB.sequelize.query(`SELECT
		role.*, t_user_role.user FROM t_role AS role JOIN t_user_role WHERE
		role.id = t_user_role.role_id
		and t_user_role.user IN (?)`, { replacements: [user], type: DB.sequelize.QueryTypes.SELECT })
		return result;
	}
	// 通过角色名称查询user
	public async findUserByRoleName(roleName: string) {
		const result: any = await DB.sequelize.query(`SELECT t_user_role.user FROM t_user_role JOIN t_role ON t_user_role.role_id = t_role.id WHERE t_role.name LIKE :roleName`, { replacements: { roleName: `%${roleName}%` }, type: DB.sequelize.QueryTypes.SELECT })
		return result;
	}
	// 查询当前角色是否已存在
	public async findRoleByRoleName({namespace, roleName, name}:{namespace: string, roleName: string, name: string}) {
		const result: any = await this.Role.findOne({
			where: Sequelize.or(
				{ namespace, role: roleName },
				{ namespace, name }
			)
		})
		return result;
	}

	// 通过id和namespace获取角色信息
	public async findRoleById(Data: any): Promise<any> {
		const { namespace, id } = Data;
		const result: any = await this.Role.findOne({
			where: { namespace, id }
		})
		return result;
	}
}

export default RoleService;