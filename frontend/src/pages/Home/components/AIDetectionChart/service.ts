import request from '@/utils/request';

export async function getAIDetectionStats(): Promise<any> {
  return request({
    url: '/ai/merge/problem-chart',
    method: 'GET',
  });
} 