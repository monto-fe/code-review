import { NextFunction, Request, Response } from 'express';
import GitlabService from '../services/gitlab.service';
import { ResponseMap } from '../utils/const';
import ResponseHandler from '../utils/responseHandler';

const { ParamsError } = ResponseMap

interface NamespaceResult {
  rows: any[]; 
  count: number; 
}

class GitlabController {
  public GitlabService = new GitlabService();


  public getGitlabList = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const result: NamespaceResult = await this.GitlabService.getGitlabInfo();
      const { rows, count } = result;
      return ResponseHandler.success(res, { data: rows||[], total: count });
	} catch (error: any) {
      return ResponseHandler.error(res, error);
	}
  }

  public createGitlabToken = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { api, token, status, gitlab_version, gitlab_url, expired } = req.body;
      const data = {
        api,
        token,
        status,
        gitlab_version,
        expired,
        gitlab_url
      }
    
      if (!api || !token ) {
        return ResponseHandler.error(res, ParamsError);
      }

      const result = await this.GitlabService.addGitlabToken({
        ...data
      });

      return ResponseHandler.success(res, result);
    } catch (error: any) {
      next(error);
    }
  }

  public updateGitlabToken = async(req: Request, res: Response, next: NextFunction) => {
    try {
      const body: any = req.body;
      const { id } = body;
      if(!id){
        return ResponseHandler.error(res, ParamsError, 'id is required');
      }
      const response: any = await this.GitlabService.updateGitlabInfo(body)
      if(response){
        return ResponseHandler.success(res, response);
      }else{
        return ResponseHandler.error(res, ParamsError);
      }
    } catch (error) {
        next(error);
    }
  }

//   public deleteGitlabToken = async (req: Request, res: Response, next: NextFunction) => {
//     try {
//       const { id } = req.body as any;
//       const result: any = await this.GitlabService.deleteSelf({ id });
//       if(result){
//         return ResponseHandler.success(res);
//       }else{
//         return ResponseHandler.error(res, ParamsError);
//       }
// 		} catch (error) {
// 			next(error);
// 		}
//   }
}

export default GitlabController;