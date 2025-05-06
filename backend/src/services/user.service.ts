import md5 from 'md5';
import { Op } from 'sequelize';

import DB from '../databases';
import { HttpException } from '../exceptions/HttpException';
import { User, UpdateUserReq, UserReq, LoginParams, UserRole } from '../interfaces/user.interface';
import { Role } from '../interfaces/role.interface';
import { getUnixTimestamp, DefaultPassword } from '../utils';
import { generateSignToken, TokenSecretKey, TokenExpired } from '../utils/util'

class UserService {
  public User = DB.Users;
  public Token = DB.Token;
  public UserRole = DB.UserRole;
  public Role = DB.Role;
  public now = getUnixTimestamp()

  public async findUserById({id}:{id: number}): Promise<User> {
    const findUser: User | null = await this.User.findOne({ where: { id } });
    if (!findUser) throw new HttpException(409, "You're not user");
    return findUser;
  }

  public async findUserAndRoleById({id}:{id: number}): Promise<any> {
    const findUser: User | null = await this.User.findOne({ where: { id } });
    // 查user信息
    if (!findUser) throw new HttpException(409, "You're not user");
    // 查user关联的roleId信息
    const user = findUser.user;
    const findUserRole: UserRole[] = await this.UserRole.findAll({ where: { user } });
    // 查roleList信息
    const roleIds: number[] = findUserRole.map((item: UserRole) => item.role_id); 
    const roleList: Role[] = await this.Role.findAll({ where: { id: roleIds } });
    return {
      userInfo: findUser,
      roleList
    }
  }

  public async findUserByUsername({user, namespace}:LoginParams): Promise<User | null> {
    const findUser: User | null = await this.User.findOne({
      where: { user, namespace },
      raw: true
    });
    return findUser
  }

  public async checkUsernameExists(namespace: string, user: string): Promise<boolean | null> {
    const result = await this.User.findOne({
      where: { user, namespace }
    })
    return result !== null
  }

  public async checkEmailExists(email: string): Promise<boolean | null> {
    const user = await this.User.findOne({
      where: { email }
    })
    return user !== null
  }

  public async login(Data: LoginParams): Promise<any> {
    const { id, user, namespace } = Data;
    // TODO: 先检查当前的有效期是否过期，如果不过期则取最新的值
    const existingToken = await this.Token.findOne({ 
      where: { user },
      order: [['create_time', 'DESC']] 
    });
    console.log("this.now", this.now)
    const now = getUnixTimestamp()
    if (existingToken && existingToken.dataValues.expired_at > now) {
      return { jwtToken: existingToken.dataValues.token };
    }else {
      const jwtToken = generateSignToken({id, user, namespace, create_time: now }, TokenSecretKey, TokenExpired)
      const expired = now + 8 * 60 * 60
      if(jwtToken) {
        try{
          await this.Token.create({
            token: jwtToken,
            user,
            expired_at: expired,
            create_time: now,
            update_time: now
          })
          return { jwtToken }
        }catch(err) {
          console.log("err", err)
          return { jwtToken: '' }
        }
      }
    }
  }

  public async createUser(Data: User): Promise<User> {
    const { job } = Data;
    const createUser: User = await this.User.create({ ...Data, job });
    return createUser;
  }

  public async updateUser(Data: UpdateUserReq): Promise<any> {
    const { user, name, email, phone_number, job } = Data;
    const { id } = Data;
    const res: any = await this.User.update({
      user, name, email, phone_number, job
    }, {
      where: {
        id
      }
    })
    return res
  }

  // getUserList
  public async getUserList(query: any): Promise<any> {
    const { offset, limit, user, id, ...rest } = query
    let params:any = {
      where: {
        ...rest
      },
      offset,
      limit,
      order: [
        ['id', 'DESC']
      ],
      attributes: {
        exclude: ['password']
      }
    }
    if(user && user.length){
      params.where.user = {
        [Op.in]: user
      }
    }
    if(id){
      params.where.id = id
    }
    const { count, rows }  = await this.User.findAndCountAll(params);
    return { data: rows, count: count };
  }
  // createUser
  public async createInnerUser(Data: UserReq): Promise<any> {
    const { namespace, user, name, job, password=DefaultPassword, email, phone_number, role_ids, operator } = Data;
    const t = await DB.sequelize.transaction();
    try {
      // insert user info
      const userInfo: any = await this.User.create({
        namespace,
        user,
        name,
        job,
        password: md5(password),
        email,
        phone_number
      }, { transaction: t });
      if(Array.isArray(role_ids) && role_ids.length > 0){
        const batchRoles = role_ids.map((item: any) => {
          return {
            namespace,
            user,
            role_id: item,
            status: 1,
            operator,
            create_time: this.now,
            update_time: this.now
          }
        })
        await this.UserRole.bulkCreate(batchRoles, { transaction: t });
      }
      t.commit();
      return userInfo;
    }catch(err){
      console.log("err", err)
      t.rollback();
      return err
    }
  }
  // updateInnerUser
  public async updateInnerUser(Data: UpdateUserReq): Promise<any> {
    const { namespace, id, role_ids, operator, password, name, phone_number, job, email, user } = Data;
    const t = await DB.sequelize.transaction();
    try {
      await this.User.update({
        password, name, phone_number, job, email
      }, { where: { id, namespace } }, { transaction: t });
      // 更新角色信息
      if(Array.isArray(role_ids) && role_ids.length > 0){
        await this.UserRole.destroy({where: {
          namespace,
          user
        }}, { transaction: t })
        await this.UserRole.bulkCreate(role_ids.map((item: number) => {
          return {
            namespace,
            user,
            role_id: item,
            status: 1,
            operator,
            create_time: this.now,
            update_time: this.now
          }
        }), { transaction: t })
      }
      t.commit();
      return true
    }catch(err){
      console.log("err", err)
      t.rollback();
      return err;
    }
  }
  // deleteUser
  public async deleteUser({id, namespace, user}:{ id: number, namespace: string, user: string}): Promise<any> {
    const res: any = await this.User.destroy({where: {
      id, namespace, user
    }})
    return res
  }

  // 通过用户名查询user[]
	public async findUserByUserName(userName: string) {
		const result: any = await DB.sequelize.query(`SELECT
		t_user.user FROM t_user WHERE
		t_user.name LIKE ?`, { replacements: [`%${userName}%`], type: DB.sequelize.QueryTypes.SELECT })
		return result
	}
}

export default UserService;
