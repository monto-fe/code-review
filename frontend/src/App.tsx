import { memo, useContext, useEffect } from 'react';
import { ConfigProvider, theme } from 'antd';
import dayjs from 'dayjs';
import { observer } from 'mobx-react-lite';

import { setHtmlLang } from '@/utils/i18n';
import Routes from '@/config/routes';
import { BasicContext } from '@/store/context';
import { getAntdI18nMessage } from '@/store/i18n';

export default memo(
  observer(() => {
    const context = useContext(BasicContext) as any;
    const { i18nLocale, globalConfig } = context.storeContext;
    const antdMessage = getAntdI18nMessage(i18nLocale);

    dayjs.locale(i18nLocale.toLocaleLowerCase());

    useEffect(() => {
      setHtmlLang(i18nLocale);
    }, []);

    return (
      <ConfigProvider
        locale={antdMessage}
        theme={{
          algorithm: globalConfig.theme === 'all-dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
        }}
      >
        <Routes />
      </ConfigProvider>
    );
  }),
);
