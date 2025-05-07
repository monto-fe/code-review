import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import WebhookController from '../controllers/webhook.controller';
import AIMessageController from '../controllers/aiMessage.controller';
import AIConfigController from '../controllers/aiConfig.controller';

class Route implements Routes {
  public router = Router();
  public WebhookController = new WebhookController();
  public AIMessageController = new AIMessageController();
  public AIConfigController = new AIConfigController();
  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/webhook/merge', this.WebhookController.AICheck)
    this.router.get('/ai/message', this.WebhookController.GetAIMessage)
    this.router.get('/ai/config', this.AIConfigController.GetAIConfig)
    this.router.post('/ai/config', this.AIConfigController.CreateAIConfig)
    this.router.put('/ai/config', this.AIConfigController.UpdateAIConfig)
  }
}

export default Route;