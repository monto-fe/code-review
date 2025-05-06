export interface TableQueryParam {
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  namespace: string;
  category?: string;
  resource: string;
  properties: string;
  name: string;
  describe: string;
  operator: string;
  create_time: number;
  update_time: number;
}
