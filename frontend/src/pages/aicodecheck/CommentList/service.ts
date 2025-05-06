import request from '@/utils/request';
import { TableQueryParam } from './data.d';

const config = {
  baseURL: import.meta.env.VITE_APP_AI_APIHOST || '',
}

const namespace = 'acl';

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/message?namespace=${namespace}`,
    method: 'get',
    params,
  });
}

