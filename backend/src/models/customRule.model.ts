import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { CustomRule } from '../interfaces/custom.interface';

export type CustomRuleCreationAttributes = Optional<CustomRule, 'id'>;

export class CustomRuleModel extends Model<CustomRule, CustomRuleCreationAttributes> implements CustomRule {
  public id!: number;
  public project_name!: string;
  public project_id!: string;
  public rule!: string;
  public status!: number;
  public operator!: string;

  public create_time!: number;
  public update_time!: number;
}

export default function (sequelize: Sequelize): typeof CustomRuleModel {
  CustomRuleModel.init(
    {
      id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
      },
      project_name: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: '项目名称',
      },
      project_id: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: '规则对应的项目Id',
      },
      rule: {
        type: DataTypes.TEXT,
        allowNull: false,
        comment: '规则内容',
      },
      status: {
        type: DataTypes.INTEGER,
        defaultValue: 1,
        allowNull: true,
        comment: '状态 -1:禁用 1:启用',
      },
      operator: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: '操作人',
      },
      create_time: {
        type: DataTypes.INTEGER,
        allowNull: false,
        comment: '创建时间',
        defaultValue: 0,
      },
      update_time: {
        type: DataTypes.INTEGER,
        allowNull: false,
        comment: '更新时间',
        defaultValue: 0
      }
    },
    {
      sequelize,
      tableName: 't_custom_rule',
      timestamps: false
    }
  );

  return CustomRuleModel;
}