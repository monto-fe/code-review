import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data.d';

const config = {
  baseURL: import.meta.env.VITE_APP_AI_APIHOST || '',
}

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/common-rule?namespace=${namespace}`,
    method: 'get',
    params,
  });
}

export async function createData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: '/common-rule',
    method: 'post',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function updateData(params: TableListItem): Promise<any> {
  return request({
    ...config,
    url: `/common-rule`,
    method: 'put',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function removeData(id: number): Promise<any> {
  return request({
    ...config,
    url: `/common-rule`,
    method: 'delete',
    data: {
      id,
      namespace,
    },
  });
}
