/**
 * 自定义 request 网络请求工具,基于axios
 * @author duheng1992
 */
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { notification } from 'antd';
import settings from '@/config/settings';
import { getToken } from '@/utils/localToken';

export interface ResponseData<T = unknown> {
  ret_code: number;
  data: {
    data: T;
    [key: string]: any;
  };
  message?: string;
  total?: number;
}

const customCodeMessage: { [key: number]: string } = {
  0: '请求成功',
  10000: '用户名已存在',
  10001: '系统异常，请稍后重试',
  10002: '用户名或密码错误',
  10003: '请求参数错误',
  10004: '邮箱已被使用',
  10005: '该用户还未登陆',
  10006: '验证码错误', // 注册时
  10008: '未知错误',
  10009: '登录失效，请重新登录',
};

const serverCodeMessage: { [key: number]: string } = {
  200: 'Request Successfully',
  400: 'Bad Request',
  401: 'Unauthorized',
  403: 'Forbidden',
  404: 'Not Found',
  500: 'Internal Server Error',
  502: 'Bad Gateway',
  503: 'Service Unavailable',
  504: 'Gateway Timeout',
};

const errorHandler = (error: any) => {
  const { response, message, code } = error;
  if (message === 'CustomError') {
    // 自定义错误
    const { config, data } = response;
    const { url, baseURL } = config;
    const { ret_code, message } = data;
    const reqUrl = url.split('?')[0].replace(baseURL, '');
    const noVerifyBool = settings.ajaxResponseNoVerifyUrl.includes(reqUrl);
    if (!noVerifyBool) {
      notification.error({
        message: `提示 (Tips)`,
        description: customCodeMessage[ret_code] || message || 'Error',
      });

      console.log(`当前code: ${ret_code}`);
      if (ret_code === 10005 && ret_code === 10009) {
        setTimeout(() => {
          window.location.href = '/user/login';
        }, 500);
      }
    }
  } else if (message === 'CancelToken') {
    // TODO 取消请求 Token
    // eslint-disable-next-line no-console
  } else if (code === 'ERR_NETWORK') {
    notification.error({
      message: `请求错误 ${message}`,
    });
  } else if (response && response.status) {
    const errorText = serverCodeMessage[response.status] || response.statusText;
    const { status, request } = response;
    notification.error({
      message: `请求错误 ${status}: ${request.responseURL}`,
      description: errorText,
    });
  } else if (!response) {
    notification.error({
      description: '您的网络发生异常，无法连接服务器',
      message: '网络异常',
    });
  }

  return Promise.reject(error);
};

/**
 * 配置request请求时的默认参数
 */
const request = axios.create({
  baseURL: import.meta.env.VITE_APP_APIHOST || '',
  withCredentials: false, // 当跨域请求时不发送cookie
  timeout: 5000, // 请求超时时间,5000(单位毫秒) / 0 不做限制
});

// 全局设置 - post请求头
// request.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8';

/**
 * 请求前
 * 请求拦截器
 */
request.interceptors.request.use(
  (config: AxiosRequestConfig & { cType?: boolean }) => {
    // 如果设置了cType 说明是自定义 添加 Content-Type类型 为自定义post 做铺垫
    if (config.cType) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers['Content-Type'] = 'application/json;charset=UTF-8';
    }

    // 自定义添加token header
    const headerToken = getToken();
    if (headerToken) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers[settings.ajaxHeadersTokenKey] = headerToken;
    }

    return config;
  },
  /* ,error=> {} */ // 已在 export default catch
);

/**
 * 请求后
 * 响应拦截器
 */
request.interceptors.response.use(
  (response: AxiosResponse<ResponseData>) => {
    const res = response.data || {};
    const { ret_code } = res;

    if (ret_code !== 0) {
      // eslint-disable-next-line prefer-promise-reject-errors
      return Promise.reject({
        response,
        message: 'CustomError',
      });
    }

    return response;
  },
  /* , error => {} */ // 已在 export default catch
);

/**
 * ajax 导出
 *
 * Method: get
 *     Request Headers
 *         无 - Content-Type
 *     Query String Parameters
 *         name: name
 *         age: age
 *
 * Method: post
 *     Request Headers
 *         Content-Type:application/json;charset=UTF-8
 *     Request Payload
 *         { name: name, age: age }
 *         Custom config parameters
 *             { cType: true }  Mandatory Settings Content-Type:application/json;charset=UTF-8
 * ......
 */
export default async function ajax<T = any, R = AxiosResponse<T>>(
  config: AxiosRequestConfig & { cType?: boolean },
): Promise<AxiosResponse<R, any>> {
  return request(config)
    .then((response: AxiosResponse) => response.data)
    .catch((error) => errorHandler(error));
}
