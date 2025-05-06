import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import ResourceController from '../controllers/resource.controller';

class Route implements Routes {
  public router = Router();
  public ResourceController = new ResourceController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.get('/resource', this.ResourceController.getResources)
    this.router.post('/resource', this.ResourceController.createResource)
    this.router.put('/resource', this.ResourceController.updateResource)
    this.router.delete('/resource', this.ResourceController.deleteResource)
  }
}

export default Route;