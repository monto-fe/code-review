import settings from '@/config/settings';

import { Theme, NavMode } from '@/@types/settings.d';
import { TabNavItem } from '@/@types/router.d';

export interface StateType {
  /* 以下是针对所有 Layout 扩展字段 */
  // 左侧展开收起
  collapsed: boolean;
  // 头部固定开启
  headFixed: boolean;
  // 模板主题
  theme: Theme;

  /* 以下是针对 UniversalLayout 扩展字段 */
  // tab菜单开启
  tabNavEnable: boolean;
  // 头部tab导航列表
  headTabNavList: TabNavItem[];
  // 菜单栏模式
  navMode: NavMode;
}

export const initialGlobalState: StateType = {
  collapsed: false,
  headFixed: settings.headFixed,
  theme: settings.theme,
  tabNavEnable: settings.tabNavEnable,
  headTabNavList: [],
  navMode: settings.navMode,
};
