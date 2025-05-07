// models/aiConfig.model.ts
import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { AIConfigAttributes } from '../interfaces/AIConfig.interface';


type AIConfigCreationAttributes = Optional<AIConfigAttributes, 'id'>;

export class AIConfigModel extends Model<AIConfigAttributes, AIConfigCreationAttributes> implements AIConfigAttributes {
  public id!: number;
  public name!: string;
  public api_url!: string;
  public api_key!: string;
  public model!: string;
  public is_active!: boolean;
  public readonly create_time!: number;
  public readonly update_time!: number;
}
export default function (sequelize: Sequelize): typeof AIConfigModel {
  AIConfigModel.init(
    {
      id: {
        type: DataTypes.INTEGER.UNSIGNED,
        autoIncrement: true,
        primaryKey: true,
      },
      name: {
        type: DataTypes.STRING(64),
        allowNull: false,
      },
      api_url: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      api_key: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      model: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      is_active: {
        type: DataTypes.BOOLEAN,
        allowNull: false,
        defaultValue: false,
      },
      create_time: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
        allowNull: false
      },
      update_time: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
        allowNull: true
      }
    },
    {
      sequelize,
      tableName: 't_ai_config',
      timestamps: false
    }
  );
return AIConfigModel
}