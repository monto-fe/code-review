import { memo, Suspense, useContext } from 'react';
import { Link } from 'react-router-dom';
import { Flex, theme } from 'antd';

import BreadCrumbs from '@/components/BreadCrumbs';
import SelectLang from '@/components/SelectLang';

import logoDark from '@/assets/images/logo-dark.svg';
import logoWhite from '@/assets/images/logo-white.svg';
import RightTopUser from './RightTopUser';
import RightTopMessage from './RightTopMessage';

import Settings from './Settings';

import { IRouter, IPathKeyRouter, BreadcrumbType } from '@/@types/router';
import { BasicContext } from '@/store/context';
import { IRoleInfo } from '@/@types/permission';
import LeftSider from './LeftSider';
// import { useI18n } from '@/store/i18n';

export interface RightTopProps {
  menuData: IRouter[];
  jsonMenuData?: IPathKeyRouter;
  routeItem: IRouter;
  userRoles?: IRoleInfo[];
  breadCrumbs?: BreadcrumbType[];
}

export default memo(({ menuData, routeItem, userRoles = [], breadCrumbs = [] }: RightTopProps) => {
  const { storeContext } = useContext(BasicContext) as any;
  const { globalConfig } = storeContext;
  const {
    token: { colorTextLightSolid, colorTextBase },
  } = theme.useToken();

  return (
    <div className='universallayout-right-top'>
      <Flex className='universallayout-right-top-header'>
        <Flex
          className='universallayout-right-top-top'
          align='center'
          justify='space-between'
          wrap="nowrap"
        >
          {globalConfig.navMode === 'inline' ? (
            <>
              <div className='universallayout-top-menu'>
                <BreadCrumbs className='breadcrumb' list={breadCrumbs} />
              </div>
              <Flex gap={12} className='universallayout-top-menu-right'>
                <Suspense fallback={<>loading...</>}>
                  <RightTopMessage />
                </Suspense>
                <RightTopUser />
                <SelectLang className='universallayout-top-selectlang' />
                <Settings />
              </Flex>
            </>
          ) : (
            <>
              <Link
                to='/'
                className='logo-url'
                style={{ color: globalConfig.theme === 'light' ? colorTextBase : colorTextLightSolid }}
              >
                <Flex align='center' gap={8}>
                  <img alt='' src={globalConfig.theme === 'light' ? logoDark : logoWhite} width='130' />
                </Flex>
              </Link>
              <LeftSider
                collapsed={globalConfig.collapsed}
                userRoles={userRoles}
                menuData={menuData}
                routeItem={routeItem}
                theme={globalConfig.theme}
                mode='horizontal'
              />
              <Flex
                gap={12}
                className='universallayout-top-menu-right'
                style={{ color: globalConfig.theme === 'light' ? colorTextBase : colorTextLightSolid }}
              >
                <Suspense fallback={<>loading...</>}>
                  <RightTopMessage />
                </Suspense>
                <RightTopUser />
                <SelectLang className='universallayout-top-selectlang' />
                <Settings />
              </Flex>
            </>
          )}
        </Flex>
      </Flex>
    </div>
  );
});
