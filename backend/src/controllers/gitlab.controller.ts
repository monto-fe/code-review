import { NextFunction, Request, Response } from 'express';
import GitlabService from '../services/gitlab.service';
import { GitlabManagerService } from '../services/gitlabManager.service';
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
      const result = await this.GitlabService.getGitlabInfo();
      return ResponseHandler.success(res, { data: result });
	} catch (error: any) {
      return ResponseHandler.error(res, error);
	}
  }

  public createGitlabToken = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { api, token, status, gitlab_version, webhook_name, webhook_url, gitlab_url, expired, source_branch, target_branch } = req.body;
      const data = {
        api,
        token,
        webhook_name,
        webhook_url,
        status,
        gitlab_version,
        expired,
        gitlab_url,
        source_branch,
        target_branch
      }
    
      if (!api || !token ) {
        return ResponseHandler.error(res, ParamsError);
      }

      const result = await this.GitlabService.addGitlabToken({
        ...data
      });

      // 添加成功后，刷新缓存
      await GitlabManagerService.init();

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
      // 更新成功后，刷新缓存
      await GitlabManagerService.init();
      if(response){
        return ResponseHandler.success(res, response);
      }else{
        return ResponseHandler.error(res, ParamsError);
      }
    } catch (error) {
        next(error);
    }
  }

  public refreshGitlabToken = async (req: Request, res: Response, next: NextFunction) => { 
    try {
      await GitlabManagerService.init();
      return ResponseHandler.success(res);
    } catch (error) {
        next(error);
    }
  }

  // public deleteGitlabToken = async (req: Request, res: Response, next: NextFunction) => {
  //   try {
  //     const { id } = req.body as any;
  //     const result: any = await this.GitlabService.deleteSelf({ id });
  //     if(result){
  //       return ResponseHandler.success(res);
  //     }else{
  //       return ResponseHandler.error(res, ParamsError);
  //     }
	// 	} catch (error) {
	// 		next(error);
	// 	}
  // }
  public deleteGitlabToken = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { id } = req.body as any;
      const result: any = await this.GitlabService.deleteGitlabToken({ id });
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

export default GitlabController;