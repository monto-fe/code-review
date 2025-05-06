import request from '@/utils/request';

export async function queryCurrent(): Promise<any> {
  return request({
    url: '/userInfo',
    method: 'get',
  });
}
