import request from '@/utils/request';
import { TableQueryParam, TableListItem } from './data';

export async function queryList(params: TableQueryParam): Promise<any> {
  return request({
    url: '/ai/check/history',
    method: 'get',
    params,
  });
}

export async function removeData(id: string): Promise<any> {
  return request({
    url: `/ai/check/history/${id}`,
    method: 'delete',
  });
}

export async function getTaskDetail(id: string): Promise<any> {
  return request({
    url: `/ai/check/result/${id}`,
    method: 'get',
  });
}

export async function createTask(data: any): Promise<any> {
  return request({
    url: '/ai/check/create',
    method: 'post',
    data,
  });
}

export async function startManualCheck(data: any): Promise<any> {
  return request({
    url: '/ai/check/manual',
    method: 'post',
    data,
  });
}

export async function getAIModelOptions(): Promise<any> {
  return request({
    url: '/ai/manager',
    method: 'get',
  });
} 