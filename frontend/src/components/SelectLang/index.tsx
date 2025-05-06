import { memo, useCallback, useContext, useMemo } from 'react';
import { Dropdown } from 'antd';
import { TranslationOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { setLocale } from '@/utils/i18n';

import { ItemType } from 'antd/lib/menu/interface';
import { I18nKey } from '@/@types/i18n.d';
import { BasicContext } from '@/store/context';

export interface SelectLangProps {
  className?: string;
}

export default memo(
  observer(({ className }: SelectLangProps) => {
    const { storeContext } = useContext(BasicContext) as any;
    const { i18nLocale, updateI18n } = storeContext;

    const menuItems = useMemo<ItemType[]>(
      () => [
        {
          key: 'zh-CN',
          label: <> 简体中文</>,
          icon: <>🇨🇳 </>,
          disabled: i18nLocale === 'zh-CN',
        },
        {
          key: 'zh-TW',
          label: <> 繁体中文</>,
          icon: <>🇭🇰 </>,
          disabled: i18nLocale === 'zh-TW',
        },
        {
          key: 'en-US',
          label: <> English</>,
          icon: <>🇺🇸 </>,
          disabled: i18nLocale === 'en-US',
        },
      ],
      [i18nLocale],
    );

    const onMenuClick = useCallback(
      ({ key }: { key: string }) => {
        const lang = key as I18nKey;
        storeContext.updateI18n(lang);
        setLocale(lang);
      },
      [i18nLocale, updateI18n],
    );
    return (
      <Dropdown className={className} menu={{ items: menuItems, onClick: onMenuClick }} arrow>
        <span className='cursor'>
          <TranslationOutlined />
        </span>
      </Dropdown>
    );
  }),
);
