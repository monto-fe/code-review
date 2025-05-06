import { Router } from 'express';

import { Routes } from '../interfaces/routes.interface';
import UserController from '../controllers/user.controller';

class Route implements Routes {
  public router = Router();
  public UserController = new UserController();

  constructor() {
    this.initializeRoutes();
  }

  private initializeRoutes() {
    this.router.post('/login', this.UserController.login)
    this.router.get('/user', this.UserController.getUserList)
    this.router.post('/user', this.UserController.createInnerUser)
    this.router.put('/user', this.UserController.updateInnerUser)
    this.router.delete('/user', this.UserController.deleteUser)
    this.router.get('/userInfo', this.UserController.getUserInfo)
  }
}

export default Route;