import { Sequelize, DataTypes, Model, Optional } from 'sequelize';
import { GitlabInfo } from '../interfaces/gitlab.interface';

export type GitlabInfoCreationAttributes = Optional<GitlabInfo, 'id'>;

export class GitlabModel extends Model<GitlabInfo, GitlabInfoCreationAttributes> implements GitlabInfo {
    public id!: number;
    public api!: string;
    public token!: string;
    public status!: 1 | -1;
    public gitlab_version!: string;
    public expired!: number;
    public gitlab_url!: string;
    public create_time!: number;
    public update_time!: number;
}

export default function (sequelize: Sequelize): typeof GitlabModel {
    GitlabModel.init(
    {
      id: {
        type: DataTypes.INTEGER,
        autoIncrement: true,
        primaryKey: true,
      },
      api: {
          type: DataTypes.STRING(100),
          allowNull: false,
          comment: 'GitLab API 标识',
      },
      token: {
          type: DataTypes.STRING(255),
          allowNull: false,
          comment: 'GitLab API Token',
      },
      status: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '状态（启用/禁用）',
      },
      gitlab_version: {
          type: DataTypes.STRING(20),
          allowNull: false,
          comment: 'GitLab 版本',
      },
      expired: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '过期时间',
      },
      gitlab_url: {
          type: DataTypes.STRING(255),
          allowNull: false,
          comment: 'GitLab URL',
      },
      create_time: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '创建时间',
      },
      update_time: {
          type: DataTypes.INTEGER,
          allowNull: false,
          comment: '更新时间'
      },
    },
    {
      sequelize,
      tableName: 't_gitlab_info',
      timestamps: false
    },
  );

  return GitlabModel;
}