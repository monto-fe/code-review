import { memo, useContext, useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { Flex, Menu, theme as antdTheme } from 'antd';
import { ItemType } from 'antd/lib/menu/interface';
import { MenuTheme } from 'antd/lib';
import { observer } from 'mobx-react-lite';

import settings from '@/config/settings';
import logoDark from '@/assets/images/logo-dark.svg';
import logoWhite from '@/assets/images/logo-white.svg';
import { IRouter } from '@/@types/router';
import { IRoleInfo } from '@/@types/permission';
import ALink from '@/components/ALink';
import { isExternal } from '@/utils/validate';
import { hasPermissionRoles } from '@/utils/router';
import { useI18n } from '@/store/i18n';
import { BasicContext } from '@/store/context';

export interface LeftSiderProps {
  menuData: IRouter[];
  routeItem: IRouter;
  userRoles?: IRoleInfo[];
  collapsed?: boolean;
  theme?: MenuTheme;
  mode?: 'inline' | 'horizontal';
}

const createMenuItems = (
  t: (key: string) => string,
  userRoles: IRoleInfo[],
  routes: IRouter[],
  parentPath = '/',
): ItemType[] => {
  const items: ItemType[] = [];

  for (let index = 0; index < routes.length; index++) {
    const item = routes[index];

    // 验证权限
    if (!hasPermissionRoles(userRoles, item.meta?.roles)) {
      // eslint-disable-next-line no-continue
      continue;
    }

    const Icon = item.meta?.icon || undefined;
    const hidden = item.meta?.hidden || false;
    if (hidden === true) {
      // eslint-disable-next-line no-continue
      continue;
    }

    // 设置路径
    let path = item.path || '';
    if (!isExternal(item.path)) {
      path = item.path.startsWith('/')
        ? item.path
        : `${parentPath.endsWith('/') ? parentPath : `${parentPath}/`}${item.path}`;
    }

    const title = item.meta?.title || '--';

    if (item.children) {
      items.push({
        key: path,
        label: (
          <Flex align='center'>
            {Icon && <Icon />}
            <span>{t(title)}</span>
          </Flex>
        ),
        children: createMenuItems(t, userRoles, item.children, path),
      });
    } else {
      items.push({
        key: path,
        label: (
          <ALink to={path}>
            <Flex align='center'>
              {Icon && <Icon />}
              <span>{t(title)}</span>
            </Flex>
          </ALink>
        ),
      });
    }
  }

  return items;
};

export default memo(
  observer(
    ({ menuData, routeItem, userRoles = [], collapsed = false, theme = 'dark', mode = 'inline' }: LeftSiderProps) => {
      const context = useContext(BasicContext) as any;
      const { storeContext } = context;
      const { i18nLocale } = storeContext;

      const t = useI18n(i18nLocale);
      const {
        token: { colorTextLightSolid, colorTextBase },
      } = antdTheme.useToken();

      const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
      const [openKeys, setOpenKeys] = useState<string[]>([]);

      const onOpenChange = (opens: string[]) => {
        setOpenKeys(opens);
      };

      useEffect(() => {
        if (routeItem) {
          setSelectedKeys([routeItem.path]);
          if (collapsed) {
            setOpenKeys([]);
          } else if (routeItem.meta?.parentPath && mode === 'inline') {
            setOpenKeys([routeItem.meta?.parentPath[0]]);
          }
        }
      }, [routeItem, collapsed, mode]);

      if (mode !== 'inline') {
        return (
          <Menu
            theme={theme === 'light' ? 'light' : 'dark'}
            mode={mode}
            selectedKeys={selectedKeys}
            className='menu'
            items={createMenuItems(t, userRoles, menuData)}
          />
        )
      }

      return (
        <div className='universallayout-left'>
          <div className='universallayout-left-sider'>
            <Flex align='center' justify='center' className='logo'>
              <Link to='/' style={{ color: theme === 'light' ? colorTextBase : colorTextLightSolid }}>
                {collapsed ? (
                  settings.siteAbbreviationTitle
                ) : (
                  <img
                    alt=''
                    src={theme === 'light' ? logoDark : logoWhite}
                    height='100'
                    className='mt-12'
                  />
                )}
              </Link>
            </Flex>
            <div className='universallayout-left-menu'>
              <Menu
                theme={theme === 'light' ? 'light' : 'dark'}
                mode={mode}
                selectedKeys={selectedKeys}
                openKeys={openKeys}
                onOpenChange={onOpenChange}
                items={createMenuItems(t, userRoles, menuData)}
              />
            </div>
          </div>
        </div>
      );
    },
  ),
);
