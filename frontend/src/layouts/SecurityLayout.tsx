import { memo, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';

import PageLoading from '@/components/PageLoading';

import { ResponseData } from '@/utils/request';
import { queryCurrent } from '@/services/user';
import { CurrentUser } from '@/store/user';
import { BasicContext } from '@/store/context';

export interface SecurityLayoutProps {
  children: React.ReactNode;
}

export default memo(
  observer(({ children }: SecurityLayoutProps) => {
    const navigate = useNavigate();
    const context = useContext(BasicContext) as any;
    const { storeContext } = context;
    const user = storeContext.userInfo;
    const [isLoading, setIsLoading] = useState(true);

    const isLogin = useMemo(() => user.id > 0, [user]);
    const getUser = useCallback(async () => {
      try {
        setIsLoading(true);
        const response: ResponseData<CurrentUser> = await queryCurrent();
        const { data: { userInfo, roleList } } = response;
        storeContext.updateUserInfo({
          ...userInfo,
          roleList: roleList || []
        });
      } catch (error: any) {
        // 无论什么错误，都跳转到登录页面
        const redirect = window.location.pathname + window.location.search;
        navigate(`/user/login?redirect=${encodeURIComponent(redirect)}`, { replace: true });
      } finally {
        setIsLoading(false);
      }
    }, [navigate, storeContext]);

    useEffect(() => {
      // 如果用户已经登录，不需要重新获取用户信息
      if (isLogin) {
        setIsLoading(false);
        return;
      }
      
      // 检查是否有token，如果没有token直接跳转登录
      const token = localStorage.getItem('monto_acl_react_token');
      if (!token) {
        const redirect = window.location.pathname + window.location.search;
        navigate(`/user/login?redirect=${encodeURIComponent(redirect)}`, { replace: true });
        return;
      }
      
      getUser();
    }, [isLogin, getUser, navigate]);

    // 如果正在加载，显示加载页面
    if (isLoading) {
      return <PageLoading />;
    }

    // 如果未登录，不渲染任何内容（因为已经跳转到登录页面）
    if (!isLogin) {
      return null;
    }

    // 已登录，渲染子组件
    return <>{children}</>;
  }),
);
