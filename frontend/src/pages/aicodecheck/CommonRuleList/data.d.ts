export interface TableQueryParam {
  id?: number;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  language: string;
  name: string;
  rule: string;
  description: string;
  create_time?: number;
  update_time?: number;
}
