import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data';

const config = {
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
}

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/gitlab`,
    method: 'get',
    params,
  });
}

export async function createData(params: any): Promise<any> {
  return request({
    ...config,
    url: '/gitlab',
    method: 'post',
    data: {
      ...params,
      gitlab_version: '1',
      gitlab_url: 'xxx'
    },
  });
}

export async function updateData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: `/gitlab`,
    method: 'put',
    data: {
      ...params,
      gitlab_version: '1',
      gitlab_url: 'xxx'
    },
  });
}

export async function removeData(id: number): Promise<any> {
  return request({
    ...config,
    url: `/gitlab`,
    method: 'delete',
    data: {
      id,
      namespace,
    },
  });
}

export async function getDetail(id: number): Promise<any> {
  const res = await request({
    ...config,
    url: `/gitlab/token/${id}`,
    method: 'get',
  });
  return res.data as unknown as any;
}
