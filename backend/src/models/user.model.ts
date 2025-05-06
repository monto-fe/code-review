import dayjs from 'dayjs';
import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { User } from '../interfaces/user.interface';

export type UserCreationAttributes = Optional<User, 'id'>;

export class UserModel extends Model<User, UserCreationAttributes> implements User {
  public id: number;
  public o_id: number;
  public namespace: string;
  public user: string;
  public name: string;
  public job: string;
  public password: string;
  public phone_number: string;
  public email: string;
  public create_time: number;
  public update_time: number;
  public delete_time: number;
}

export default function (sequelize: Sequelize): typeof UserModel {
    UserModel.init(
    {
      id: {
        type: DataTypes.INTEGER.UNSIGNED, // you can omit the `new` but this is discouraged
        autoIncrement: true,
        primaryKey: true,
        comment: '自有id, 应该是不会使用的',
      },
      o_id: {
          type: DataTypes.INTEGER, // you can omit the `new` but this is discouraged
          unique: true,
          allowNull: true,
          comment: 'oa那边的Id, 负数代表虚拟用户',
      },
      namespace: {
          type: new DataTypes.STRING(128),
          allowNull: false,
          comment: '项目名称',
      },
      user: {
          type: new DataTypes.STRING(255),
          allowNull: false,
          comment: '用户英文名',
      },
      name: {
          type: new DataTypes.STRING(255),
          allowNull: false,
          comment: '用户中文名',
      },
      job: {
          type: new DataTypes.STRING(255),
          allowNull: true,
          defaultValue: '',
          comment: '职位名称',
      },
      password: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        defaultValue: '',
        comment: '密码',
      },
      phone_number: {
          type: new DataTypes.STRING(255),
          allowNull: true,
          defaultValue: '',
          comment: '手机号码',
      },
      email: {
          type: new DataTypes.STRING(255),
          allowNull: true,
          defaultValue: '',
          comment: '邮箱',
      },
      create_time: {
          type: DataTypes.INTEGER.UNSIGNED,
          allowNull: false,
          defaultValue: '0',
          comment: '创建时间',
      },
      update_time: {
          type: DataTypes.INTEGER.UNSIGNED,
          allowNull: false,
          defaultValue: '0',
          comment: '更新时间',
      }
    },
    {
      tableName: 't_user',
      sequelize,
      timestamps: false,
      hooks: {
          beforeCreate: (user) => {
              const now = dayjs().unix()
              user.create_time = now
              user.update_time = now
          }
      }
    },
  );

  return UserModel;
}