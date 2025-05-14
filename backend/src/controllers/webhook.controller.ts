import { NextFunction, Request, Response } from 'express';
import axios from 'axios';
import AICheckService from '../services/aiCheck.service';
// import GitlabService from '../services/gitlab.service';
import { GitlabManagerService } from '../services/gitlabManager.service';
import ResponseHandler from '../utils/responseHandler';
import { pageCompute } from '../utils/pageCompute';
import { PushWeChatInfo } from '../utils/util';
import { GitlabCache } from '../interfaces/gitlab.interface';

class WebhookController {
  public AICheckService = new AICheckService();
  // public GitlabService = new GitlabService();

  public AICheck = async (req: Request, res: Response, next: NextFunction) => {
    const { project: { id, path_with_namespace }, object_attributes: { iid, url: merge_url } } = req.body;
    ResponseHandler.success(res, { projectId: id, mergeRequestId: iid }, 'Webhook处理成功，等待AI检测');
    // 获取gitlab信息
    // 读取gitlab缓存信息
    const gitlabService = await GitlabManagerService.init();
    const gitlabInfoCache: GitlabCache = gitlabService.getCache();
    const gitlabInfoResult = this.findTokenByProjectId(gitlabInfoCache, id)
    if(!gitlabInfoResult){
      ResponseHandler.error(res, {}, '请配置gitlab Token');
      return
    }
    const { token: gitlabToken, config: { gitlabAPI, webhook_url } } = gitlabInfoResult;

    // const { api: gitlabAPI, token: gitlabToken, webhook_url } = gitlabInfo[0].dataValues;
    
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

    // 推送webhook
    if(webhook_url){
        const webhookContent = PushWeChatInfo({
            path_with_namespace,
            merge_url,
            result
        });
        await this.sendMarkdownToWechatBot(webhook_url, webhookContent);
        console.log("推送webhook成功", mergeRequest);
    }
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
  
  /**
   * 向企业微信机器人发送 Markdown 消息
   * @param {string} markdownContent - Markdown 格式的内容
   */
  public sendMarkdownToWechatBot = async (webhookURL: string, markdownContent: string) => {
    try {
      const res = await axios.post(webhookURL, {
        msgtype: 'markdown',
        markdown: {
          content: markdownContent
        }
      });

      console.log('发送成功:', res.data);
      return res.data;
    } catch (error:any) {
      console.error('发送失败:', error.response?.data || error.message);
      throw error;
    }
  };

  public findTokenByProjectId(gitlabInfoCache: GitlabCache, projectId: string): {
    token: string;
    config: {
      projectids: string[];
      gitlabAPI: string;
      webhook_url: string;
    };
  } | null {
    for (const [token, config] of Object.entries(gitlabInfoCache)) {
      if (config.projectids.includes(projectId)) {
        return { token, config };
      }
    }
    return null;
  }
}

export default WebhookController;