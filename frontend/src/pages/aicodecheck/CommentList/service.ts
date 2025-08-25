import request from '@/utils/request';
import { TableQueryParam } from './data.d';

const config = {
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
}

export async function queryList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/message/list`,
    method: 'post',
    data: params,
  });
}

export async function updateRating(id: number, human_rating: number, remark?: string): Promise<any> {
  return request({
    ...config,
    url: `/ai/message`,
    method: 'put',
    data: {
      id,
      human_rating,
      remark
    },
  });
}

export async function queryProjectNamespaceList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/project-namespaces`,
    method: 'get',
    params,
  });
}

export async function queryProjectList(params?: TableQueryParam): Promise<any> {
  return request({
    ...config,
    url: `/ai/project-namespaces`,
    method: 'get',
    params,
  });
}