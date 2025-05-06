import { NextFunction, Request, Response } from 'express';
import { ResourceReq, Resource } from '../interfaces/resource.interface';
import ResourceService from '../services/resource.service';
import { ResponseMap, ResourceCategory } from '../utils/const';
import { pageCompute } from '../utils/pageCompute';
import ResponseHandler from '../utils/responseHandler';

const { ParamsError } = ResponseMap

class ResourceController {
  public ResourceService = new ResourceService();

  public getResources = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, resource, name, category, current, pageSize } = req.query as any;
      if (!namespace) {
        return ResponseHandler.error(res, ParamsError);
      }
      const { offset, limit } = pageCompute(current, pageSize);

      const { rows, count }: any = await this.ResourceService.findWithAllChildren(
				{
          namespace,
          resource,
          name,
          category,
          offset,
          limit
        }
			);
      return ResponseHandler.success(res, {
        data: rows || [],
        total: count
      });
		} catch (error) {
			next(error);
		}
  }

  public createResource = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { namespace, category, resource, properties, name, describe } : ResourceReq = req.body;
      const remoteUser = req.headers['remoteUser'] as string;
      // 参数校验
      if(!namespace || !resource || !name){
        return ResponseHandler.error(res, ParamsError);
      }
      const params: any = {
        namespace,
        resource,
        name,
        describe,
        operator: remoteUser
      }
      if(!ResourceCategory.includes(category)){
        return ResponseHandler.error(res, ParamsError);
      }else{
        params.category = category;
      }
      if(properties){
        params.properties = properties;
      }
      // 校验name是否已存在
      const ResourceResult = await this.ResourceService.findByNamespaceAndName({
        namespace,
        name
      })
      if(ResourceResult){
        return ResponseHandler.error(res, ParamsError, 'resource is exist');
      }
      
			const createData: any = await this.ResourceService.create(params);
      
      if(createData){
        return ResponseHandler.success(res, createData);
      }else{
        return ResponseHandler.error(res, ParamsError);
      }
		} catch (error) {
			next(error);
		}
  }

  public updateResource = async(req: Request, res: Response, next: NextFunction) => {
    try {
      const { id, namespace, category, resource, properties, name, describe }: Resource = req.body;
      // 先查询要更新的数据，判断是否存在
      if(!id || !namespace){
        return ResponseHandler.error(res, ParamsError);
      }
      const body: any = {
        id,
        namespace
      }
      const currentResource = await this.ResourceService.findAllById(id)
      if(!currentResource){
        return ResponseHandler.error(res, ParamsError, 'resource is not exist');
      }
      // 校验name是否已存在
      if(name){
        body.name = name;
      }
      if(!ResourceCategory.includes(category)){
        return ResponseHandler.error(res, ParamsError);
      }else{
        body.category = category;
      }
      if(properties){
        body.properties = properties;
      }
      if(describe){
        body.describe = describe;
      }
      if(resource){
        body.resource = resource;
      }

      const response: number[] = await this.ResourceService.update(body)
      return ResponseHandler.success(res, response);
		} catch (error) {
			next(error);
		}
  }

  public deleteResource = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const body: ResourceReq = req.body;
      const { id } = body;
    
      const getData: any = await this.ResourceService.deleteSelf({
				id
      });
      return ResponseHandler.success(res, getData);
		} catch (error) {
			next(error);
		}
  }
}

export default ResourceController;