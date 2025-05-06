import { NextFunction, Request, Response } from 'express';

import { RolePermissionReq, UserRoleReq } from '../interfaces/permission.interface';
import PermissionService from '../services/permission.service';
import { ResponseMap } from '../utils/const';
import { getUnixTimestamp } from '../utils';
import { pageCompute } from '../utils/pageCompute';
import ResponseHandler from '../utils/responseHandler';

const { ParamsError } = ResponseMap

class PermissionController {
  public PermissionService = new PermissionService();

  public getRolePermissions = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, role_id, current, pageSize } = req.body;
      const params: any = {
        namespace,
        role_id,
        ...pageCompute(current, pageSize)
      }
      const { rows, count }: any = await this.PermissionService.getRolePermissions(params);
      return ResponseHandler.success(res, {
        data: rows || [],
        total: count
      });
		} catch (error) {
			next(error);
		}
  }

  public getSelfPermissions = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, category } = req.body;
      const remoteUser = req.headers['remoteUser'] as string;
      const params = {
        namespace,
        user: remoteUser,
        category
      }
      const result: any = await this.PermissionService.getSelfPermissions(params);
      return ResponseHandler.success(res, result);
    }catch(error) {
      next(error);
    }
  }

  // 角色分配权限
  public assetRolePermission = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, role_id, resource_ids, describe }: RolePermissionReq = req.body;
      const remoteUser = req.headers['remoteUser']
      
      if (!namespace || !role_id || !resource_ids || !describe) {
        return ResponseHandler.error(res, ParamsError);
      }
      const params: any = {
        namespace,
        role_id,
        describe,
        create_time: getUnixTimestamp(),
        update_time: getUnixTimestamp(),
        operator: remoteUser
      }
      if(!resource_ids || !resource_ids.length){
        return ResponseHandler.error(res, ParamsError);
      }
      
      // 检查当前角色是否已经有了该权限，根据角色id查询当前role_permission是否有resource_id
      const result = await this.PermissionService.getRoleResourceIds({ namespace, role_id })
      const resource_id_list = result.map((item:any) => {
        return item.dataValues.resource_id
      })
      
      // 数组取差值，对比resource_ids和resource_id_list的值，取出resource_ids中值不属于resource_id_list的值组成的数组
      params.resource_ids = resource_ids.filter((item: any) => {
        return !resource_id_list.includes(item)
      })
			const createData: any = await this.PermissionService.assertBulkRolePermission(
				params
			);
      return ResponseHandler.success(res, createData);
		} catch (error) {
			next(error);
		}
  }

  // 给用户分配角色
  public assertUserRole = async(req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, user, role_id, role_ids, describe }: UserRoleReq = req.body;
      const remoteUser = req.headers['remoteUser'] as string;
      const now = getUnixTimestamp();
      if(role_id){
        role_ids.push(role_id);
      }
      const params = {
          namespace,
          user,
          role_ids,
          operator: remoteUser,
          create_time: now,
          update_time: now,
          describe
      }
      const response: any = await this.PermissionService.AssertUserRole(params)
      return ResponseHandler.success(res, response); 
		} catch (error) {
			next(error);
		}
  }

  public cancelRolePermission = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, role_id, resource_id } = req.body;
      const response: any = await this.PermissionService.cancelRolePermissionReq({
				namespace, role_id, resource_id
      });
      return ResponseHandler.success(res, response);
		}catch  (error) {
			next(error);
		}
  }

  public cancelUserRole = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, user, role_id, role_ids=[] } = req.body;
      // TODO: 在controller校验参数，service不处理参数
      const result: any = await this.PermissionService.cancelUserRole({
				namespace, user, role_id, role_ids
      });
      if (result) {
        return ResponseHandler.success(res);
      }else{
        return ResponseHandler.error(res, Error);
      }
		} catch (error) {
			next(error);
		}
  }
}

export default PermissionController;