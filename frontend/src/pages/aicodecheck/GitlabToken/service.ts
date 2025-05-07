import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data';

const config = {
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
}

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/gitlab-info`,
    method: 'get',
    params,
  });
}

export async function createData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: '/gitlab-info',
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
    url: `/gitlab-info`,
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
    url: `/gitlab-info`,
    method: 'delete',
    data: {
      id,
      namespace,
    },
  });
}
