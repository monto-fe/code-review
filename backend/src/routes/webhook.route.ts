import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import WebhookController from '../controllers/webhook.controller';
import AIMessageController from '../controllers/aiMessage.controller';
import AIConfigController from '../controllers/aiConfig.controller';
import GitlabController from '../controllers/gitlab.controller';

class Route implements Routes {
  public router = Router();
  public WebhookController = new WebhookController();
  public AIMessageController = new AIMessageController();
  public AIConfigController = new AIConfigController();
  public GitlabController = new GitlabController();
  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/webhook/merge', this.WebhookController.AICheck)
    this.router.get('/ai/message', this.WebhookController.GetAIMessage)
    // 更新human_rating和remark字段
    this.router.put('/ai/message', this.AIMessageController.UpdateHumanRatingAndRemark)


    this.router.get('/ai-manager', this.AIConfigController.GetAIConfig)
    this.router.post('/ai-manager', this.AIConfigController.CreateAIConfig)
    this.router.put('/ai-manager', this.AIConfigController.UpdateAIConfig)

    this.router.get('/gitlab-info', this.GitlabController.getGitlabList)
    this.router.post('/gitlab-info', this.GitlabController.createGitlabToken)
    this.router.put('/gitlab-info', this.GitlabController.updateGitlabToken)
    this.router.delete('/gitlab-info', this.GitlabController.deleteGitlabToken)

    // 加载所有gitlab token
    this.router.post('/gitlab/token/refresh', this.GitlabController.refreshGitlabToken)

  }
}

export default Route;