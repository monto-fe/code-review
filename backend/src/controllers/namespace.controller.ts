import { NextFunction, Request, Response } from 'express';
import { NamespaceReq, Namespace, NamespaceReqList } from '../interfaces/namespace.interface';
import NamespaceService from '../services/namespace.service';
import { ResponseMap } from '../utils/const';
import { pageCompute } from '../utils/pageCompute';
import ResponseHandler from '../utils/responseHandler';

const { ParamsError } = ResponseMap

interface NamespaceResult {
  rows: any[]; 
  count: number; 
}

class NamespaceController {
  public NamespaceService = new NamespaceService();

  public getProjects = async (req: Request, res: Response, next: NextFunction) => {
    const { namespace, current, pageSize } = req.query as any;
    const pageParams = pageCompute(current, pageSize)
    const params: NamespaceReqList = {
      ...pageParams
    }
    if(namespace){
      params.namespace = namespace;
    }
    
    try {
      const result: NamespaceResult = await this.NamespaceService.findWithAllChildren(params);
      const { rows, count } = result;
      return ResponseHandler.success(res, { data: rows||[], total: count });
		} catch (error: any) {
      return ResponseHandler.error(res, error);
		}
  }

  public createProject = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const query: NamespaceReq = req.body;
      const remoteUser = req.headers['remoteUser'] as string;
      const { namespace, parent, name, describe } = query;
      
      if (!namespace || !name ) {
        return ResponseHandler.error(res, ParamsError);
      }

      const result = await this.NamespaceService.checkNamespace({
        namespace,
        name
      });

      if(result){
        return ResponseHandler.error(res, ParamsError, 'namespace or name already exists');
      }

      const namespaceParams: any = {
        namespace,
        name,
        operator: remoteUser
      }
      if (parent) {
        namespaceParams.parent = parent
      }
      if (describe) {
        namespaceParams.describe = describe
      }

			const createData: any = await this.NamespaceService.create(namespaceParams);
      return ResponseHandler.success(res, createData);
		} catch (error) {
			next(error);
		}
  }

  public updateProject = async(req: Request, res: Response, next: NextFunction) => {
    try {
      const body: Namespace = req.body;
      const { id } = body;
      if(!id){
        return ResponseHandler.error(res, ParamsError, 'id is required');
      }
      const response: any = await this.NamespaceService.update(body)
      if(response){
        return ResponseHandler.success(res, response);
      }else{
        return ResponseHandler.error(res, ParamsError);
      }
		} catch (error) {
			next(error);
		}
  }

  public deleteProject = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id } = req.body as any;
      const result: any = await this.NamespaceService.deleteSelf({ id });
      if(result){
        return ResponseHandler.success(res);
      }else{
        return ResponseHandler.error(res, ParamsError);
      }
		} catch (error) {
			next(error);
		}
  }
}

export default NamespaceController;