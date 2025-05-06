/**
 * 自定义 token 操作
 * @author duheng1992
 */
import settings from '@/config/settings';

/**
 * 获取本地
 */
export const getToken = () => localStorage.getItem(settings.siteTokenKey);

/**
 * 设置存储本地
 */
export const setToken = (token: string) => {
  // 若使用cookie，https 下应设置 { secure: true, sameSite: 'strict' }
  localStorage.setItem(settings.siteTokenKey, token);
};

/**
 * 移除本地Token
 */
export const removeToken = () => {
  localStorage.removeItem(settings.siteTokenKey);
};
