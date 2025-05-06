export interface TableQueryParam {
  id?: number;
  role?: string;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  namespace: string;
  role: string;
  name: string;
  describe: string;
  operator: string;
  resource?: any[];
  create_time: number;
  update_time: number;
}

export type Permission = {
  resource_id: number;
};
