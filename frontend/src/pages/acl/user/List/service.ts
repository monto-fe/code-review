import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data.d';

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    url: `/user?namespace=${namespace}`,
    method: 'get',
    params,
  });
}

export async function createData(params: TableListItem): Promise<any> {
  return request({
    url: '/user',
    method: 'post',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function updateData(params: TableListItem): Promise<any> {
  return request({
    url: `/user`,
    method: 'put',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function removeData(id: number, user: string): Promise<any> {
  return request({
    url: `/user`,
    method: 'delete',
    data: {
      id,
      user,
      namespace,
    },
  });
}
