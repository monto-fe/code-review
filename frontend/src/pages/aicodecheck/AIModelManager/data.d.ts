export interface TableQueryParam {
  id?: number;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  project_name: string;
  project_id: string;
  rule: string;
  create_time?: number;
  update_time?: number;
}

export interface AIModelCreateItem {
  api_key: string;
  api_url: string;
  type: string;
  model: string;
  name?: string;
  is_active: number;
}