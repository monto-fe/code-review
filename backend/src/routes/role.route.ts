import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import RoleController from '../controllers/role.controller';

class Route implements Routes {
  public router = Router();
  public RoleController = new RoleController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.get('/role', this.RoleController.getRoles)
    this.router.post('/role', this.RoleController.createRole)
    this.router.put('/role', this.RoleController.updateRole)
    this.router.delete('/role', this.RoleController.deleteRole)
    this.router.post('/role/permission', this.RoleController.createRolePermission)
  }
}

export default Route;