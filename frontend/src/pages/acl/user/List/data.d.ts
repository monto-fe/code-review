import { TableListItem as RoleTableListItem } from '@/pages/role/data.d';

export interface TableQueryParam {
  id?: number;
  user?: string;
  userName?: string;
  roleName?: string;
  current?: number;
  pageSize?: number;
  order?: SortOrder | undefined;
  role: RoleTableListItem[];
  field?: React.Key | readonly React.Key[] | undefined;
  role_ids?: number[];
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
  role: RoleTableListItem[];
  roleName?: string;
  email: 'string';
  password?: string;
  role_ids?: number[];
  create_time: number;
  update_time: number;
}
