import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { CommonRule } from '../interfaces/common.interface';

export type CommonRuleCreationAttributes = Optional<CommonRule, 'id'>;

export class CommonRuleModel extends Model<CommonRule, CommonRuleCreationAttributes> implements CommonRule {
    public id: number; // Note that the `null assertion` `!` is required in strict mode.
    public name: string;
    public language: 'Python' | 'Java' | 'JavaScript' | 'Golang' | 'Ruby' | 'C++' | 'other';
    public rule: string;
    public description: string;
    public readonly create_time: number;
    public readonly update_time: number;
}

export default function (sequelize: Sequelize): typeof CommonRuleModel {
    CommonRuleModel.init(
    {
      id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
      },
      name: {
          type: new DataTypes.STRING(128),
          allowNull: false,
          comment: '规则名称',
          defaultValue: '',
      },
      language: {
          type: new DataTypes.STRING(255),
          allowNull: false,
          comment: '语言',
          defaultValue: '',
      },
      rule: {
          type: new DataTypes.TEXT,
          allowNull: false,
          comment: '规则',
          defaultValue: '',
      },
      description: {
          type: new DataTypes.STRING(255),
          allowNull: true,
          comment: '描述',
          defaultValue: '',
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
        comment: '创建时间',
        defaultValue: 0,
      }
    },
    {
      sequelize,
      tableName: 't_common_rule',
      timestamps: false
    }
  );

  return CommonRuleModel;
}