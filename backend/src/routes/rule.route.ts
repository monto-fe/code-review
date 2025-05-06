import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import AIMessageController from '../controllers/aiMessage.controller';

class Route implements Routes {
  public router = Router();
  public AIMessageController = new AIMessageController();
  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/common-rule', this.AIMessageController.CreateCommonRule)
    this.router.get('/common-rule', this.AIMessageController.GetCommonRule)
    this.router.put('/common-rule', this.AIMessageController.UpdateCommonRule)
    this.router.get('/custom-rule', this.AIMessageController.GetCustomRule)
    this.router.post('/custom-rule', this.AIMessageController.CreateCustomRule)
    this.router.put('/custom-rule', this.AIMessageController.UpdateCustomRule)
  }
}

export default Route;