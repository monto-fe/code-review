/**
 * UniversalLayout 路由配置 入口
 * @author duheng1992
 */

import { lazy } from 'react';

import {
  HomeOutlined,
  DashboardOutlined,
  InsuranceOutlined,
  UserOutlined,
  TeamOutlined,
  KeyOutlined,
  ToolOutlined,
  QuestionCircleTwoTone,
  CodeOutlined,
  OrderedListOutlined
} from '@ant-design/icons';

import { IRouter } from '@/@types/router.d';

// 这里是业务路由
// 若要配置权限，请使用 meta.roles
const universalLayoutRotes: IRouter[] = [
  {
    path: '/home',
    meta: {
      icon: HomeOutlined,
      title: 'universal-layout.menu.home',
    },
    redirect: '/home/workplace',
    children: [
      {
        path: 'workplace',
        meta: {
          icon: DashboardOutlined,
          title: 'universal-layout.menu.home.workplace',
        },
        component: lazy(() => import('@/pages/Home')),
      },
    ],
  },
  // {
  //   path: 'https://chromewebstore.google.com/detail/%E5%A4%9A%E5%BD%A9%E4%B9%A6%E7%AD%BE/ilcmekmgeldhckdembghkiohkdffihpe?hl=zh-CN&utm_source=ext_sidebar',
  //   meta: {
  //     icon: ToolOutlined,
  //     title: 'universal-layout.menu.home.bookmark',
  //     selectLeftMenu: '/home',
  //   },
  // },
  // {
  //   path: 'https://chromewebstore.google.com/detail/%E6%97%B6%E9%97%B4%E6%88%B3%E8%BD%AC%E6%8D%A2-%E6%97%B6%E5%8C%BA%E6%97%B6%E9%92%9F/pjcapgdifnhgkkojoggfdlijpelpohcf?hl=zh-CN&utm_source=ext_sidebar',
  //   meta: {
  //     icon: ToolOutlined,
  //     title: 'universal-layout.menu.home.timestamp',
  //     selectLeftMenu: '/home',
  //   },
  // },
  // {
  //   path: 'testMenuPermission',
  //   meta: {
  //     icon: QuestionCircleTwoTone,
  //     title: 'test.menuacl',
  //     roles: ['admin'],
  //   },
  //   component: lazy(() => import('@/pages/test/TestMenuPermission')),
  // },
  // {
  //   path: 'testPagePermission',
  //   meta: {
  //     icon: QuestionCircleTwoTone,
  //     title: 'test.pageacl',
  //   },
  //   component: lazy(() => import('@/pages/test/TestPagePermission')),
  // },
  // {
  //   path: 'testAPIPermission',
  //   meta: {
  //     icon: QuestionCircleTwoTone,
  //     title: 'test.apiacl',
  //   },
  //   component: lazy(() => import('@/pages/test/TestAPIPermission')),
  // },
  {
    path: '/aicodecheck',
    redirect: '/aicodecheck/commentList',
    meta: {
      icon: CodeOutlined,
      title: 'universal-layout.menu.aicodecheck',
      roles: ['admin'],
    },
    children: [
      {
        path: 'commentList',
        meta: {
          icon: OrderedListOutlined,
          title: 'universal-layout.menu.aicodecheck.commentlist',
        },
        component: lazy(() => import('@/pages/aicodecheck/CommentList')),
      },
      {
        path: 'commonRule',
        meta: {
          icon: OrderedListOutlined,
          title: 'universal-layout.menu.aicodecheck.commonrule',
        },
        component: lazy(() => import('@/pages/aicodecheck/CommonRuleList')),
      },
      {
        path: 'customRule',
        meta: {
          icon: OrderedListOutlined,
          title: 'universal-layout.menu.aicodecheck.customrule',
        },
        component: lazy(() => import('@/pages/aicodecheck/CustomRuleList')),
      },
      {
        path: 'AIModel',
        meta: {
          icon: OrderedListOutlined,
          title: 'universal-layout.menu.aicodecheck.AIModel',
        },
        component: lazy(() => import('@/pages/aicodecheck/AIModelManager')),
      },
      {
        path: 'GitlabToken',
        meta: {
          icon: OrderedListOutlined,
          title: 'universal-layout.menu.aicodecheck.GitlabToken',
        },
        component: lazy(() => import('@/pages/aicodecheck/GitlabToken')),
      },
    ],
  },
  {
    path: '/acl',
    redirect: '/acl/user',
    meta: {
      icon: InsuranceOutlined,
      title: 'universal-layout.menu.roles',
      roles: ['admin'],
    },
    children: [
      {
        path: 'user',
        meta: {
          icon: UserOutlined,
          title: 'universal-layout.menu.roles.user',
        },
        component: lazy(() => import('@/pages/acl/user/List')),
      },
      {
        path: 'role',
        meta: {
          icon: TeamOutlined,
          title: 'universal-layout.menu.roles.role',
        },
        component: lazy(() => import('@/pages/acl/role/List')),
      },
      {
        path: 'resource',
        meta: {
          icon: KeyOutlined,
          title: 'universal-layout.menu.roles.resource',
        },
        component: lazy(() => import('@/pages/acl/resource/List')),
      },
    ],
  },
];

export default universalLayoutRotes;
