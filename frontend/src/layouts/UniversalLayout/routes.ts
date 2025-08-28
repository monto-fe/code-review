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
  CodeOutlined,
  OrderedListOutlined,
} from '@ant-design/icons';

import { IRouter } from '@/@types/router.d';

// 这里是业务路由
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
  {
    path: '/aicodecheck',
    redirect: '/aicodecheck/commentList',
    meta: {
      icon: CodeOutlined,
      title: 'universal-layout.menu.aicodecheck',
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
      // {
      //   path: 'commonRule',
      //   meta: {
      //     icon: OrderedListOutlined,
      //     title: 'universal-layout.menu.aicodecheck.commonrule',
      //   },
      //   component: lazy(() => import('@/pages/aicodecheck/CommonRuleList')),
      // },
      // {
      //   path: 'customRule',
      //   meta: {
      //     icon: OrderedListOutlined,
      //     title: 'universal-layout.menu.aicodecheck.customrule',
      //   },
      //   component: lazy(() => import('@/pages/aicodecheck/CustomRuleList')),
      // },
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
      {
        path: 'GitlabConfig',
        meta: {
          icon: OrderedListOutlined,
          hidden: true,
          title: 'universal-layout.menu.aicodecheck.GitlabConfig',
        },
        component: lazy(() => import('@/pages/aicodecheck/GitlabToken/CreateToken')),
      },
      {
        path: 'GitlabTokenDetail',
        meta: {
          icon: OrderedListOutlined,
          hidden: true,
          title: 'universal-layout.menu.aicodecheck.GitlabTokenDetail',
        },
        component: lazy(() => import('@/pages/aicodecheck/GitlabToken/detail')),
      },
    ],
  },
  {
    path: '/acl',
    redirect: '/acl/user',
    meta: {
      icon: InsuranceOutlined,
      title: 'universal-layout.menu.user.title',
    },
    children: [
      {
        path: 'user',
        meta: {
          icon: UserOutlined,
          title: 'universal-layout.menu.user',
        },
        component: lazy(() => import('@/pages/acl/user/List')),
      },
    ],
  },
];

export default universalLayoutRotes;
