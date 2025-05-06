import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { Namespace } from '../interfaces/namespace.interface';

export type NamespaceAttributes = Optional<Namespace, 'id'>;

export class NamespaceModel extends Model<Namespace, NamespaceAttributes> implements Namespace {
  public id: number;
  public namespace: string;
  public parent: string;
  public name: string;
  public describe: string;
  public operator: string;
  public readonly create_time: number;
  public readonly update_time: number;
}

export default function (sequelize: Sequelize): typeof NamespaceModel {
  NamespaceModel.init(
    {
      id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
      },
      namespace: {
          type: new DataTypes.STRING(128),
          allowNull: false,
          comment: '项目组标识',
          defaultValue: '',
      },
      parent: {
          type: new DataTypes.STRING(128),
          allowNull: true,
          comment: '所属父项目组',
          defaultValue: '',
      },
      name: {
          type: new DataTypes.STRING(128),
          allowNull: false,
          comment: '项目组中文名',
          defaultValue: '',
      },
      describe: {
          type: new DataTypes.STRING(128),
          allowNull: true,
          comment: '描述',
          defaultValue: '',
      },
      operator: {
          type: new DataTypes.STRING(128),
          allowNull: false,
          comment: '创建人',
          defaultValue: '',
      },
      create_time: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '创建时间',
          defaultValue: '0',
      },
      update_time: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '更新时间',
          defaultValue: '0',
      }
    },
    {
      tableName: 't_namespace',
      sequelize,
      timestamps: false
    },
  );

  return NamespaceModel;
}
