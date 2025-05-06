import { NextFunction, Request, Response } from 'express';

import { RoleReq, UpdateRoleReq } from 'interfaces/role.interface';
import RoleService from '../services/role.service';
import { ResponseMap } from '../utils/const';
import { Administrator } from '../utils/index';
import { pageCompute } from '../utils/pageCompute';
import ResourceService from '../services/resource.service';
import ResponseHandler from '../utils/responseHandler';

const { ParamsError, RoleExisted, SystemEmptyError } = ResponseMap

class RoleController {
  public RoleService = new RoleService();
  public ResourceService = new ResourceService();

  public getRoles = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, role, name, resource, current, pageSize } = req.query as any;

      // 检查必填参数
      if (!namespace) {
        return ResponseHandler.error(res, ParamsError);
      }

      // 计算分页参数
      const params: any = {
        namespace,
        ...pageCompute(current, pageSize),
      };

      // 可选参数添加到 params 中
      if (role) params.role = role;
      if (name) params.name = name;
      if (resource) params.resource = resource.split(",");

      const getData: any = await this.RoleService.findWithAllChildren(params);
      const { rows, count } = getData;
      // 聚合资源id
      const resourceIds = rows.map((item: any) => item.dataValues.resources).flat();
      if(resourceIds.length === 0){
        return ResponseHandler.success(res, { data: rows||[], total: count});
      }
      // 根据resource资源列表
      const resourceList: any = await this.ResourceService.findAllByIds(resourceIds);
      rows.forEach((item: any) => {
        item.dataValues.resource = resourceList.filter((resource: any) => item.dataValues.resources.includes(resource.id));
      });
      rows.forEach((item: any) => {
        delete item.dataValues.resources
      });
      return ResponseHandler.success(res, { data: rows||[], total: count } );
		} catch (error) {
			next(error);
		}
  }

  public createRole = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const params: RoleReq = req.body;
      const remoteUser = req.headers['remoteUser'] as string;
      const { namespace, role, name } = params;
      
      if (!params || !namespace || !role) {
        return ResponseHandler.error(res, ParamsError);
      }

      // 查询当前roleName是否存在
      const checkRole = await this.RoleService.findRoleByRoleName({namespace, roleName: role, name})
      if(checkRole){
        return ResponseHandler.error(res, RoleExisted);
      }else{
        const createData: any = await this.RoleService.create(
          {
            ...params,
            operator: remoteUser
          }
        );
        return ResponseHandler.success(res, createData);
      }
		} catch (error) {
			next(error);
		}
  }

  public updateRole = async(req: Request, res: Response, next: NextFunction) => {
    try {
      const body: UpdateRoleReq = req.body;
      const remoteUser = req.headers['remoteUser']
      const { id, namespace, role: newRole, permissions, name, describe } = body;
      
      if(!id || !namespace){
        return ResponseHandler.error(res, ParamsError);
      }
      // 查询如果是超级管理员，则不允许更改role
      const role = await this.RoleService.findRoleById({namespace, id})
      if(!role){
        return ResponseHandler.error(res, SystemEmptyError, 'Role not found');
      }
      const replaceRole = role.dataValues.role === Administrator ? Administrator : newRole;
      const response: any = await this.RoleService.update({
        id,
        namespace,
        role: replaceRole,
        permissions,
        name,
        describe,
        operator: remoteUser
      })
      return ResponseHandler.success(res, response);
		} catch (error) {
			next(error);
		}
  }

  public deleteRole = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id, namespace } = req.body;

      // 检查是否为管理员，管理员不允许删除
      const role = await this.RoleService.findRoleById({namespace, id})
      
      if(!role){
        return ResponseHandler.error(res, SystemEmptyError, 'Role not found');
      }
      
      if(role.dataValues.role === 'admin'){
        return ResponseHandler.error(res, SystemEmptyError, 'admin not allow to delete');
      }

      const getData: any = await this.RoleService.deleteSelf({
        id,
				namespace
      });
      return ResponseHandler.success(res, getData);
		} catch (error) {
			next(error);
		}
  }

  public createRolePermission = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const query: RoleReq = req.body;
      const remoteUser = req.headers['remoteUser']
      const params  = query as any;
      
      if (!params) {
        return ResponseHandler.error(res, ParamsError);
      }

			const createData: any = await this.RoleService.createRolePermission(
				{
          ...params,
          operator: remoteUser
        }
			);
      return ResponseHandler.success(res, createData);
		} catch (error) {
			next(error);
		}
  }
  
}

export default RoleController;