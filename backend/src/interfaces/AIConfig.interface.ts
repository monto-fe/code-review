export interface AIConfigAttributes {
  id: number;
  name: string;            // 配置名称，方便识别
  api_url: string;         // OpenAI 或其他服务的 URL
  api_key: string;         // 接口密钥
  model: string;           // 模型名称（如 gpt-4、gpt-3.5-turbo）
  is_active: boolean;      // 是否为当前使用的配置
  create_time: number;
  update_time: number;
}

export interface AIConfigCreationAttributes {
  id?: number;
  name: string;
  api_url: string;
  api_key: string;
  model: string;
  is_active: boolean;
}