import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { Token } from '../interfaces/user.interface';

export type UserCreationAttributes = Optional<Token, 'id'>;

export class TokenModel extends Model<Token, UserCreationAttributes> implements Token {
  public id: number;
  public token: string;
  public user: string;
  public expired_at: number;

  public readonly create_time!: number;
  public readonly update_time!: number;
}

export default function (sequelize: Sequelize): typeof TokenModel {
    TokenModel.init(
    {
      id: {
        autoIncrement: true,
        primaryKey: true,
        type: DataTypes.INTEGER
      },
      token: {
        allowNull: false,
        type: DataTypes.STRING
      },
      user: {
        allowNull: false,
        type: DataTypes.STRING
      },
      expired_at: {
        allowNull: false,
        type: DataTypes.INTEGER
      },
      create_time: {
        type: DataTypes.INTEGER
      },
      update_time: {
        type: DataTypes.INTEGER
      }
    },
    {
      tableName: 't_login_token',
      sequelize,
      timestamps: false
    },
  );

  return TokenModel;
}