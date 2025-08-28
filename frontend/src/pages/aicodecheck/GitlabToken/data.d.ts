export interface TableQueryParam {
  id?: number;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  api?: string;
  token?: string;
  expired?: any;
  webhook_url?: string;
  webhook_name?: string;
  source_branch?: string;
  target_branch?: string;
  webhook_status?: number;
  rule_check_status?: number;
  status?: number;
  comment_type?: number;
  project_ids_synced?: number; // 项目ID同步状态：1-失败，2-同步中，3-成功
  create_time?: number;
  update_time?: number;
}
