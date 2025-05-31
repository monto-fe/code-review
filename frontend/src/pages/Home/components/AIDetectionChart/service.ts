import request from '@/utils/request';

export async function getAIDetectionStats(): Promise<any> {
  return request({
    url: '/api/ai/detection/stats',
    method: 'GET',
  });
} 