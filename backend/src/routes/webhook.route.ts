import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import WebhookController from '../controllers/webhook.controller';
import AIMessageController from '../controllers/aiMessage.controller';

class Route implements Routes {
  public router = Router();
  public WebhookController = new WebhookController();
  public AIMessageController = new AIMessageController();
  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/webhook/merge', this.WebhookController.AICheck)
    this.router.get('/ai/message', this.WebhookController.GetAIMessage)
  }
}

export default Route;