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
}
