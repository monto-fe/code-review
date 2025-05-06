import { NextFunction, Request, Response } from 'express';
import AICheckService from '../services/aiCheck.service';
import GitlabService from '../services/gitlab.service';
import ResponseHandler from '../utils/responseHandler';
import { pageCompute } from '../utils/pageCompute';

class WebhookController {
  public AICheckService = new AICheckService();
  public GitlabService = new GitlabService();

  // private requestQueue: Promise<void>[] = [];

  public AICheck = async (req: Request, res: Response, next: NextFunction) => {
    const { project: { id }, object_attributes: { iid } } = req.body;
    ResponseHandler.success(res, { projectId: id, mergeRequestId: iid }, 'Webhook处理成功，等待AI检测');
    // 获取gitlab信息
    const gitlabInfo = await this.GitlabService.getGitlabInfo();
    if(!gitlabInfo.length){
        ResponseHandler.error(res, {}, '请配置gitlab Token');
        return
    }
    const { api: gitlabAPI, token: gitlabToken } = gitlabInfo[0].dataValues;
    
    // 获取merge信息
    const mergeRequest = await this.AICheckService.getMergeRequestInfo({
      gitlabAPI: gitlabAPI, projectId: id, gitlabToken
    });

    let diff:any = [];
    // 读取diff信息
    if(gitlabAPI && id && iid){
        diff = await this.AICheckService.getMergeRequestDiff({
            gitlabAPI: gitlabAPI,
            projectId: id,
            gitlabToken,
            mergeRequestId: iid
        });
    }

    const result:any = await this.AICheckService.checkMergeRequestWithAI({
      mergeRequest, 
      diff,
      gitlabAPI: gitlabAPI,
      gitlabToken: gitlabToken
    });
    console.log("AI检查结果:", result);

    // 将 AI 检查结果作为评论写入 GitLab
    let commentResponse = ''
    if(result && iid && id && gitlabToken && gitlabAPI){
        commentResponse = await this.AICheckService.postCommentToGitLab({
            comment: result,
            mergeRequestId: iid,
            projectId: id,
            gitlabToken,
            gitlabAPI
        });
    }
    console.log("评论已成功提交:", commentResponse);
  }

  public GetAIMessage = async (req: Request, res: Response, next: NextFunction) => {
    const { projectId, current, pageSize } = req.query as any;
    const pageParams = pageCompute(current, pageSize)
        const params: any = {
          ...pageParams,
          projectId: projectId
        }
    const { rows, count } = await this.AICheckService.getAIMessage(params);
    ResponseHandler.success(res, { data: rows, total: count }, 'success');
  }
}

export default WebhookController;