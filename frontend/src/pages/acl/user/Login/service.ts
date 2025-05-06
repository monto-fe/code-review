import request from '@/utils/request';
import { LoginParamsType } from './data.d';

const namespace = 'acl';

export async function accountLogin(params: LoginParamsType): Promise<any> {
  return request({
    url: '/login',
    method: 'post',
    data: { ...params, namespace },
  });
}
