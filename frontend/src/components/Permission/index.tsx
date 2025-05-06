import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { Result, Button } from 'antd';

import { BasicContext } from '@/store/context';
import { hasPermissionRoles } from '@/utils/router';
import { observer } from 'mobx-react-lite';

const Forbidden = (
  <Result
    status={403}
    title='403'
    subTitle={
      <div>
        <p>抱歉，您没有当前页面的权限.</p>
        <p>Sorry, you are not authorized to access this page.</p>
      </div>
    }
    extra={
      <Button type='primary'>
        <Link to='/'>Home</Link>
      </Button>
    }
  />
);

export interface ALinkProps {
  children: React.ReactNode;
  role?: string | string[];
  noNode?: React.ReactNode;
}

const Permission: React.FC<ALinkProps> = observer(({ role, noNode = Forbidden, children }) => {
  const context = useContext(BasicContext) as any;
  const user = context.storeContext.userInfo;
  return hasPermissionRoles(user.roleList, role) ? <>{children}</> : <>{noNode}</>;
});

export default Permission;
