import request from '@/utils/request';

export async function getMergeStats(): Promise<any> {
  return request({
    url: '/api/merge/stats',
    method: 'GET',
  });
} 