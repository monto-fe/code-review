import request from '@/utils/request';

export async function getMergeStats(data: any): Promise<any> {
  return request({
    url: '/ai/merge/check-count',
    method: 'GET',
    params: data,
  });
} 