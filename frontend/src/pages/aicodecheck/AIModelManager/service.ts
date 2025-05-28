import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data';

const config = {
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
}

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/config`,
    method: 'get',
    params,
  });
}

export async function createData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: '/ai/config',
    method: 'post',
    data: {
      ...params,
      is_active: 1
    },
  });
}

export async function updateData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: `/ai/config`,
    method: 'put',
    data: {
      ...params,
      is_active: 1,
    },
  });
}

export async function removeData(id: number): Promise<any> {
  return request({
    ...config,
    url: `/ai/config`,
    method: 'delete',
    data: {
      id,
      namespace,
    },
  });
}
