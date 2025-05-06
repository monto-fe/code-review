import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { RolePermission } from '../interfaces/role.interface';

export type RolePermissionCreationAttributes = Optional<RolePermission, 'id'>;

export class RolePermissionModel extends Model<RolePermission, RolePermissionCreationAttributes> implements RolePermission {
    public id: number; 
    public namespace: string;
    public role_id: number;
    public resource_id: number;
    public describe: string;
    public operator: string;
    public readonly create_time: number;
    public readonly update_time: number;
}

export default function (sequelize: Sequelize): typeof RolePermissionModel {
    RolePermissionModel.init(
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
        role_id: {
            type: DataTypes.INTEGER.UNSIGNED,
            allowNull: false,
            comment: '角色ID',
        },
        resource_id: {
            type: DataTypes.INTEGER.UNSIGNED,
            allowNull: false,
            comment: '资源ID',
        },
        describe: {
            type: DataTypes.TEXT,
            allowNull: true,
            comment: '描述',
        },
        operator: {
            type: new DataTypes.STRING(128),
            allowNull: false,
            comment: '创建人',
            defaultValue: '',
        },
        create_time: {
            type: DataTypes.INTEGER.UNSIGNED,
            allowNull: false,
            comment: '创建时间',
            defaultValue: '0',
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
        deletedAt: 'deletedAt',
        tableName: 't_role_permission',
        timestamps: false,
        indexes: [
            {
                name: 'role_id',
                fields: ['role_id'],
            },
            {
                name: 'resource_id',
                fields: ['resource_id'],
            }
        ]
    }
  );

  return RolePermissionModel;
}