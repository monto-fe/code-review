import { memo, useCallback, useContext, useEffect, useMemo } from 'react';
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

    const isLogin = useMemo(() => user.id > 0, [user]);
    const getUser = useCallback(async () => {
      try {
        const response: ResponseData<CurrentUser> = await queryCurrent();
        const { data: { userInfo, roleList } } = response;
        storeContext.updateUserInfo({
          ...userInfo,
          roleList: roleList || []
        });
      } catch (error: any) {
        if (error.message && error.message === 'CustomError') {
          const { response } = error;
          if (response) {
            const redirect = window.location.pathname + window.location.search;
            navigate(`/user/login?redirect=${encodeURIComponent(redirect)}`, { replace: true });
          }
        }
      }
    }, []);

    useEffect(() => {
      getUser();
    }, []);

    return <>{isLogin ? children : <PageLoading />}</>;
  }),
);
