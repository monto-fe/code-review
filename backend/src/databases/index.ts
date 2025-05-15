import Sequelize from 'sequelize';
import { NODE_ENV, DB_TYPE, DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_DATABASE } from '../config';
import UserModel from '../models/user.model';
import NamespaceModel from '../models/namespace.model';
import ResourceModel from '../models/resource.model';
import RoleModel from '../models/role.model';
import TokenModel from '../models/token.model';
import RolePermissionModel from '../models/rolePermission.model';
import UserRoleModel from '../models/userRole.model';
import AIMessageModel from '../models/aiMessage.model';
import CommonRuleModel from '../models/commonRule.model';
import CustomRuleModel from '../models/customRule.model';
import GitlabInfoModel from '../models/gitlab.model';
import AIConfigModel from '../models/AIConfig.model';

console.log('数据库地址：' + DB_HOST + ': ' + DB_PORT)
console.log('NODE_ENV', NODE_ENV)
const sequelize = new Sequelize.Sequelize(DB_DATABASE, DB_USER, DB_PASSWORD, {
  dialect: DB_TYPE,
  host: DB_HOST,
  port: DB_PORT,
  timezone: '+08:00',
  timestamps: false,
  define: {
    charset: 'utf8mb4',
    collate: 'utf8mb4_general_ci',
    freezeTableName: true,
  },
  createdAt: 'create_time',
  updatedAt: 'update_time',
  deletedAt: 'delete_time',
  paranoid: true,
  pool: {
    min: 0,
    max: 5,
  },
  logQueryParameters: true,
  logging: console.log,
  benchmark: true,
} as any);

sequelize.authenticate();

const DB:any = {
  Users: UserModel(sequelize), // 用户管理
  Token: TokenModel(sequelize), // 登录Token
  Namespace: NamespaceModel(sequelize),
  Resource: ResourceModel(sequelize),
  UserRole: UserRoleModel(sequelize),
  Role: RoleModel(sequelize),
  RolePermission: RolePermissionModel(sequelize),
  GitlabInfo: GitlabInfoModel(sequelize),
  AIMessage: AIMessageModel(sequelize),
  CommonRule: CommonRuleModel(sequelize),
  CustomRule: CustomRuleModel(sequelize),
  AIConfig: AIConfigModel(sequelize),
  sequelize, // connection instance (RAW queries)
  Sequelize, // library
};

export default DB;