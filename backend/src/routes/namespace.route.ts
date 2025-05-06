import { Router } from 'express';
import { Routes } from '../interfaces/routes.interface';
import NamespaceController from '../controllers/namespace.controller';

class Route implements Routes {
  public router = Router();
  public NamespaceController = new NamespaceController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.get('/project', this.NamespaceController.getProjects)
    this.router.post('/project', this.NamespaceController.createProject)
    this.router.put('/project', this.NamespaceController.updateProject)
    this.router.delete('/project', this.NamespaceController.deleteProject)
  }
}

export default Route;