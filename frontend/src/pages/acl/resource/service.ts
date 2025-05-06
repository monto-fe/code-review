import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data.d';

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    url: `/resource?namespace=${namespace}`,
    method: 'get',
    params,
  });
}

export async function createData(params: TableListItem): Promise<any> {
  return request({
    url: '/resource',
    method: 'post',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function updateData(params: TableListItem): Promise<any> {
  return request({
    url: `/resource`,
    method: 'put',
    data: {
      ...params,
      namespace,
    },
  });
}

export async function removeData(id: number): Promise<any> {
  return request({
    url: `/resource`,
    method: 'delete',
    data: {
      id,
      namespace,
    },
  });
}
