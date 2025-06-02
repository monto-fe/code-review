import request from '@/utils/request';
import { TableQueryParam, TableListItem, AIModelCreateItem } from './data';

const config = {
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
}

export const MODEL_TYPE_OPTIONS = [
  { label: 'UCloud', value: 'ucai' },
  { label: 'DeepSeek', value: 'deepseek' },
  { label: 'OpenAI', value: 'openai' },
];

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/config`,
    method: 'get',
    params,
  })
}

export async function createData(params: AIModelCreateItem): Promise<any> {
  return request({
    ...config,
    url: '/ai/config',
    method: 'post',
    data: {
      ...params
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
