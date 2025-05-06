import { memo, useContext, useMemo } from 'react';
import { /* Outlet, */ useLocation } from 'react-router-dom';

import { useI18n } from '@/store/i18n';

import useTitle from '@/hooks/useTitle';

import SelectLang from '@/components/SelectLang';

import { formatRoutes } from '@/utils/router';

import layoutRotes from './routes';

import './css/index.less';
import { BasicContext } from '@/store/context';

export interface UserLayoutProps {
  children: React.ReactNode;
}

export default memo(({ children }: UserLayoutProps) => {
  const location = useLocation();

  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  // 框架所有菜单路由 与 patch key格式的所有菜单路由
  const routerPathKeyRouter = useMemo(() => formatRoutes(layoutRotes), []);

  // 当前路由item
  const routeItem = useMemo(() => routerPathKeyRouter.pathKeyRouter[location.pathname], [location]);

  // 设置title
  useTitle(t(routeItem?.meta?.title || ''));

  return (
    <div className='user-layout'>
      <div className='lang'>
        <SelectLang />
      </div>
      {/* <Outlet /> */}
      {children}
    </div>
  );
});
