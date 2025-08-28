import request from '@/utils/request';

export async function queryInnerAIModel(): Promise<any> {
  return request({
    url: `/ai/manager`,
    method: 'get',
  });
}
