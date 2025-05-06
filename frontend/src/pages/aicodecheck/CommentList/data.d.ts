export interface TableQueryParam {
  id?: number;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  project_id: string;
  merge_id: string;
  ai_model: string;
  rule: string;
  rule_id: string;
  result: string;
  passed: boolean;
  description: string;
  create_time?: number;
  update_time?: number;
}
