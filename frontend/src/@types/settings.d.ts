/**
 * 站点配置 ts定义
 * @author duheng1992
 */

export type Theme = 'dark' | 'light' | 'all-dark';

export type NavMode = 'inline' | 'horizontal';

export interface SettingsType {
  /**
   * 站点名称 & 名称缩写
   */
  siteTitle: string;
  siteAbbreviationTitle: string;

  /**
   * 站点本地存储Token 的 Key值
   */
  siteTokenKey: string;

  /**
   * Ajax请求头发送Token 的 Key值
   */
  ajaxHeadersTokenKey: string;

  /**
   * Ajax返回值不参加统一验证的api地址
   */
  ajaxResponseNoVerifyUrl: string[];

  /**
   * Layout 头部固定开启
   */
  headFixed: boolean;

  /**
   * Layout 模板主题
   */
  theme: Theme;

  /**
   * UniversalLayout tab菜单开启
   */
  tabNavEnable: boolean;

  /**
   * UniversalLayout 菜单导航模式
   */
  navMode: NavMode;
}
