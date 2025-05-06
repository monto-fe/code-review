import { memo, useContext } from 'react';
import { Popover, Divider, Switch } from 'antd';
import { CheckOutlined, SettingOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import classnames from 'classnames';

import { Theme, NavMode } from '@/@types/settings';

import style from './index.module.less';
import { BasicContext } from '@/store/context';
import { i18nLocaleDefault, useI18n } from '@/store/i18n';

export default memo(
  observer(() => {
    const t = useI18n(i18nLocaleDefault);

    const { storeContext } = useContext(BasicContext) as any;
    const { globalConfig } = storeContext;

    // 模板主题
    const setTheme = (theme: Theme) => {
      storeContext.updateGlobalConfig({ ...globalConfig, theme });
    };

    // 导航模式
    const setNavMode = (navMode: NavMode) => {
      storeContext.updateGlobalConfig({ ...globalConfig, navMode });
    };

    // 固定头部
    const onChangeHeadFixed = () => {
      storeContext.updateGlobalConfig({ ...globalConfig, headFixed: !globalConfig.headFixed });
    };

    // tabNavEnable
    const onChangeTabNavEnable = () => {
      storeContext.updateGlobalConfig({ ...globalConfig, tabNavEnable: !globalConfig.tabNavEnable });
    };

    return (
      <Popover
        content={
          <div className={style.setting}>
            <div className={style['setting-title']}>{t('setting.pagestyle')}</div>

            <div className={style['setting-radio']}>
              <div
                className={classnames(style['setting-radio-item'], style['style-dark'])}
                title='dark'
                onClick={() => {
                  setTheme('dark');
                }}
              >
                {globalConfig.theme === 'dark' && (
                  <span className={style['choose-icon']}>
                    <CheckOutlined />
                  </span>
                )}
              </div>
              <div
                className={classnames(style['setting-radio-item'], style['style-light'])}
                title='light'
                onClick={() => {
                  setTheme('light');
                }}
              >
                {globalConfig.theme === 'light' && (
                  <span className={style['choose-icon']}>
                    <CheckOutlined />
                  </span>
                )}
              </div>
              <div
                className={classnames(style['setting-radio-item'], style['style-all-dark'])}
                title='all-dark'
                onClick={() => {
                  setTheme('all-dark');
                }}
              >
                {globalConfig.theme === 'all-dark' && (
                  <span className={style['choose-icon']}>
                    <CheckOutlined />
                  </span>
                )}
              </div>
            </div>

            <Divider />

            <div className={style['setting-title']}>{t('setting.navigationmode')}</div>
            <div className={style['setting-radio']}>
              <div
                className={classnames(style['setting-radio-item'], style['nav-inline'])}
                title='inline'
                onClick={() => {
                  setNavMode('inline');
                }}
              >
                {globalConfig.navMode === 'inline' && (
                  <span className={style['choose-icon']}>
                    <CheckOutlined />
                  </span>
                )}
              </div>
              <div
                className={classnames(style['setting-radio-item'], style['nav-horizontal'])}
                title='horizontal'
                onClick={() => {
                  setNavMode('horizontal');
                }}
              >
                {globalConfig.navMode === 'horizontal' && (
                  <span className={style['choose-icon']}>
                    <CheckOutlined />
                  </span>
                )}
              </div>
            </div>

            <Divider />

            <div className={style['setting-list']}>
              <div className={style['setting-list-item']}>
                <span>{t('setting.headfixed')}</span>
                <span className={style['setting-list-item-action']}>
                  <Switch
                    checkedChildren='√'
                    unCheckedChildren=''
                    checked={globalConfig.headFixed}
                    onChange={onChangeHeadFixed}
                  />
                </span>
              </div>
              <div className={style['setting-list-item']}>
                <span>{t('setting.navtabs')}</span>
                <span className={style['setting-list-item-action']}>
                  <Switch
                    checkedChildren='√'
                    unCheckedChildren=''
                    checked={globalConfig.tabNavEnable}
                    onChange={onChangeTabNavEnable}
                  />
                </span>
              </div>
            </div>
          </div>
        }
        trigger='hover'
        placement='bottomRight'
      >
        <span className='universallayout-top-settings cursor'>
          <SettingOutlined />
        </span>
      </Popover>
    );
  }),
);
