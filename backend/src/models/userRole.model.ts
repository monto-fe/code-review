import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { UserRole } from '../interfaces/role.interface';

export type UserRoleCreationAttributes = Optional<UserRole, 'id'>;

export class UserRoleModel extends Model<UserRole, UserRoleCreationAttributes> implements UserRole {
    public id: number;
    public namespace: string;
    public user: string;
    public role_id: number;
    public status: number;
    public operator: string;
    public create_time: number;
    public update_time: number;
}

export default function (sequelize: Sequelize): typeof UserRoleModel {
    UserRoleModel.init(
    {
        id: {
            type: DataTypes.INTEGER.UNSIGNED, // you can omit the `new` but this is discouraged
            autoIncrement: true,
            primaryKey: true,
        },
        namespace: {
            type: new DataTypes.STRING(128),
            allowNull: false,
            comment: '所属项目组',
            defaultValue: '',
        },
        user: {
            type: new DataTypes.STRING(255),
            allowNull: false,
            comment: '用户英文名',
            defaultValue: '',
        },
        role_id: {
            type: DataTypes.INTEGER.UNSIGNED,
            allowNull: false,
            comment: '角色ID',
        },
        status: {
            type: DataTypes.TINYINT,
            allowNull: false,
            defaultValue: 1,
            comment: '当前状态(-1: 禁用；1: 可用；)',
        },
        operator: {
            type: new DataTypes.STRING(128),
            allowNull: false,
            comment: '创建人',
            defaultValue: '',
        },
        create_time: {
            type: DataTypes.INTEGER.UNSIGNED,
            comment: '创建时间',
            allowNull: false,
            defaultValue: "0",
        },
        update_time: {
            type: DataTypes.INTEGER.UNSIGNED,
            allowNull: false,
            comment: '更新时间',
            defaultValue: '0',
        }
    },
    {
        sequelize,
        tableName: 't_user_role',
        timestamps: false,
        indexes: [
            {
                name: 'namespace',
                fields: ['namespace'],
            },
            {
                name: 'user',
                fields: ['user']
            },
            {
                name: 'role_id',
                fields: ['role_id']
            }
        ]
    }
  );

  return UserRoleModel;
}