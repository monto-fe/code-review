export interface TableQueryParam {
  id?: number;
  current?: number;
  pageSize?: number;
}

export interface TableListItem {
  id: number;
  api: string;
  token: string;
  expired: any;
}
