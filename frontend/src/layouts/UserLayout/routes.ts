import { lazy } from 'react';
import { IRouter } from '@/@types/router';

const pathPre = '/user';

const UserLayoutRoutes: IRouter[] = [
  {
    path: `${pathPre}/login`,
    meta: {
      title: 'user-layout.menu.login',
    },
    component: lazy(() => import('@/pages/acl/user/Login')),
  },
  {
    path: `${pathPre}/register`,
    meta: {
      title: 'user-layout.menu.register',
    },
    component: lazy(() => import('@/pages/acl/user/Register')),
  },
];

export default UserLayoutRoutes;
