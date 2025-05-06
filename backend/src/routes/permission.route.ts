import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import PermissionController from '../controllers/permission.controller';

class Route implements Routes {
  public router = Router();
  public PermissionController = new PermissionController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/permission/assertRolePermission', this.PermissionController.assetRolePermission)
    this.router.post('/permission/assertUserRole', this.PermissionController.assertUserRole)
    this.router.post('/permission/cancelRolePermission', this.PermissionController.cancelRolePermission)
    this.router.post('/permission/cancelUserRole', this.PermissionController.cancelUserRole)
    
    this.router.post('/permission/getRolePermissions', this.PermissionController.getRolePermissions)
    this.router.post('/permission/getSelfPermissions', this.PermissionController.getSelfPermissions)
  }
}

export default Route;
