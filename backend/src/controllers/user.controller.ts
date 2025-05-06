import { NextFunction, Request, Response } from 'express';
import md5 from 'md5';

import { User, UserListReq, UserReq, UserLogin, UserListQuery } from '../interfaces/user.interface';
import UserService from '../services/user.service';
import RoleService from '../services/role.service';
import { ResponseMap } from '../utils/const';
import ResponseHandler from '../utils/responseHandler';
import { Administrator } from '../utils';
import { pageCompute } from '../utils/pageCompute';

const { SystemError, UserError, LoginError,} = ResponseMap

class UsersController {
  public UserService = new UserService();
  public RoleService = new RoleService();
  
  public login = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { user, password, namespace } : UserLogin = req.body;
      console.log('login:', user, namespace, password)
      const findData: User | null = await this.UserService.findUserByUsername({user, namespace});
      console.log('login:', findData)
      if (!findData || findData.password !== md5(password)) {
        return ResponseHandler.error(res, UserError);
      }
      // 密码校验，并写入token
      const { jwtToken } = await this.UserService.login({
        id: findData.id, user, namespace
      })
      if (!jwtToken) {
        return ResponseHandler.error(res, LoginError);
      }
      return ResponseHandler.success(res, { jwt_token: jwtToken, user });
    } catch (error) {
      next(error);
    }
  };

  public getUserInfo = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const userId = req.headers['userId'] as string;
      const { userInfo, roleList }: any = await this.UserService.findUserAndRoleById({id: Number(userId)});
      const { id, namespace, user, name, job, phone_number, email } = userInfo;
      return ResponseHandler.success(res, { id, namespace, user, name, job, phone_number, email, roleList });
    } catch (error) {
      next(error);
    }
  };

  // getUserList
  public getUserList = async (req: Request, res: Response) => {
    const { id, user, userName, roleName, current, pageSize, namespace } = req.query as unknown as UserListQuery;
    const query: UserListReq = {
      namespace,
      user: user ? [user] : [],
      ...pageCompute(current, pageSize),
      ...(id && { id: Number(id) })
    };

    try{
      // 如果roleName存在，则查询user，number[]
      let roleUserList: string[] = [];
      let userNameList: string[] = [];
      if(roleName){
        const userListLike = await this.RoleService.findUserByRoleName(roleName)
        const users = userListLike.length > 0 ? userListLike.map((item:any) => item.user) : [];
        query.user = users.concat(query.user)
        roleUserList = users;
      }
      if(userName){
        const userListLike = await this.UserService.findUserByUserName(userName)
        const users = userListLike.length > 0 ? userListLike.map((item:any) => item.user) : [];
        query.user = users.concat(query.user)
        userNameList = users;
      }
      // 如果roleName和userName都存在，就取两次查询出来的交集，否则就没有数据
      if(roleName && userName){
        const intersection: string[] = roleUserList.filter((item:string) => userNameList.includes(item))
        query.user = intersection.concat(query.user ? query.user : []);
      }
      // 如果name存在，则查询对应的user，User[]
      // 合并user,进行查询
      const result = await this.UserService.getUserList(query)
      const { data, count } = result;

      if(count > 0){
        // 给每个user，查询对应角色
        const userList = data.map((user:any) => {
          return user.dataValues.user;
        })
        const uniqueArray:any = [...new Set(userList)];
        const roleList = await this.RoleService.findRoleByUser(uniqueArray);
        data.forEach((user:any) => {
          user.dataValues.role = roleList.filter((role:any) => {
            return role.user === user.dataValues.user
          });
        })
      }
      return ResponseHandler.success(res, { data, count });
    }catch(err:any){
      return ResponseHandler.error(res, err);
    }
  }
 
  // createInnerUser
  public createInnerUser = async (req: Request, res: Response) => {
    const { namespace, user, name, job, password, email, phone_number, role_ids }: UserReq = req.body;
    const remoteUser = req.headers['remoteUser'] as string
    if(!namespace || !user){
      return ResponseHandler.error(res, SystemError, 'Namespace and user is required!');
    }
    // 1、检查新增用户是否已存在
    try{
      const checkUser = await this.UserService.checkUsernameExists(namespace, user)
      if(checkUser){
        return ResponseHandler.error(res, SystemError, 'User already exist!');
      }
      // 2、判断新增用户是否选择角色，如果选择角色
      if(!name||!password){
        return ResponseHandler.error(res, SystemError, 'Name and password is required!');
      }
      const result = await this.UserService.createInnerUser({
        namespace, user, name, job, password: md5(password), email, phone_number, role_ids, operator: remoteUser
      })
      return ResponseHandler.success(res, result);
    }catch(err:any){
      return ResponseHandler.error(res, err);
    }
  }

  public updateInnerUser = async (req: Request, res: Response) => {
    const { namespace, id=0, name, job, user, password, email, phone_number, role_ids }: UserReq = req.body;
    const operator = req.headers['remoteUser'] as string
    if(!namespace || id === 0){
      return ResponseHandler.error(res, SystemError, 'Namespace and id is required!');
    }
    // 查询当前用户信息
    const userInfo:any = await this.UserService.findUserById({ id })
    if(!userInfo){
      return ResponseHandler.error(res, SystemError, 'User not found!');
    }
    if(userInfo.namespace !== namespace){
      return ResponseHandler.error(res, SystemError, 'Namespace is not match!');
    }
    const updateInfo:any = {
      namespace, id, user: userInfo.user, operator
    }
    if(name){
      updateInfo.name = name
    }
    if(job){
      updateInfo.job = job
    }
    if(password){
      updateInfo.password = md5(password)
    }
    if(email){
      updateInfo.email = email
    }
    if(phone_number){
      updateInfo.phone_number = phone_number
    }
    if(Array.isArray(role_ids) && role_ids.length > 0){
      updateInfo.role_ids = role_ids
    }

    try{
      const result = await this.UserService.updateInnerUser(updateInfo)
      return ResponseHandler.success(res, result);
    }catch(err:any){
      return ResponseHandler.error(res, err);
    }
  }
  // deleteUser
  public deleteUser = async (req: Request, res: Response) => {
    const { namespace, id, user } = req.body;
    // params check
    if(!namespace || !id || !user){
      return ResponseHandler.error(res, SystemError, 'params is required!');
    }
    // Super Admin not allow to delete
    if(user === Administrator){
      return ResponseHandler.error(res, SystemError, 'Super Admin not allow to delete');
    }
    try{
      const result = await this.UserService.deleteUser({ namespace, id, user })
      if(result){
        return ResponseHandler.success(res, result);
      }else{
        return ResponseHandler.error(res, SystemError, 'please check your params');
      }
    }catch(err:any){
      return ResponseHandler.error(res, err);
    }
  }
}

export default UsersController;