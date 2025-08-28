export interface TableQueryParam {
  id?: number;
  user?: string;
  userName?: string;
  current?: number;
  pageSize?: number;
  order?: SortOrder | undefined;
  field?: React.Key | readonly React.Key[] | undefined;
  password?: string;
}

export interface TableListItem {
  id: number;
  o_id: number;
  namespace: 'string';
  user: 'string';
  name: 'string';
  job: 'string';
  phone_number: number;
  email: 'string';
  password?: string;
  create_time: number;
  update_time: number;
}
