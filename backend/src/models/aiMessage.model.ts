import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { AImessage, HumanRating } from '../interfaces/aiMessage.interface';

export type AImessageCreationAttributes = Optional<AImessage, 'id'>;

export class AImessageModel extends Model<AImessage, AImessageCreationAttributes> implements AImessage {
  public id!: number;
  public project_id!: number;
  public merge_id!: string;
  public merge_url!: string;
  public ai_model!: "DeepSeek" | "ChatGPT" ;
  public rule!: 1 | 2;
  public rule_id!: number;
  public result!: string;
  public human_rating!: HumanRating;
  public remark!: string;
  public passed?: boolean;
  public checked_by?: string | null;
  public create_time!: number;
}

export default function (sequelize: Sequelize): typeof AImessageModel {
  AImessageModel.init(
    {
      id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
      },
      project_id: {
        type: DataTypes.INTEGER.UNSIGNED,
        allowNull: false,
        comment: '项目ID',
      },
      merge_url: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: '合并url',
      },
      merge_id: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: '合并mergeId',
      },
      ai_model: {
        type: new DataTypes.STRING(255),
        allowNull: false,
        comment: 'AI模型',
      },
      rule: {
        type: DataTypes.INTEGER.UNSIGNED,
        allowNull: false,
        comment: '规则类型',
      },
      rule_id: {
        type: DataTypes.INTEGER.UNSIGNED,
        allowNull: false,
        comment: '规则ID',
      },
      result: {
        type: DataTypes.TEXT,
        allowNull: false,
        comment: '结果',
      },
      human_rating: {
        type: DataTypes.INTEGER.UNSIGNED,
        allowNull: false,
        comment: '人工评分',
      },
      remark: {
        type: DataTypes.TEXT,
        allowNull: false,
        comment: '备注',
      },
      passed: {
        type: DataTypes.INTEGER,
        defaultValue: 0,
        allowNull: true,
        comment: '是否通过',
      },
      checked_by: {
        type: new DataTypes.STRING(255),
        allowNull: true,
        comment: '检查人',
      },
      create_time: {
        type: DataTypes.INTEGER,
        allowNull: false,
        comment: '创建时间',
        defaultValue: 0,
      }
    },
    {
      sequelize,
      tableName: 't_ai_message',
      timestamps: false
    }
  );

  // Foreign key constraints:
  // AImessageModel.belongsTo(sequelize.models.TCommonRule, { foreignKey: 'rule_id' });
  // AImessageModel.belongsTo(sequelize.models.TCustomRule, { foreignKey: 'project_id' });

  return AImessageModel;
}