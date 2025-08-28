import { memo, useContext } from 'react';
import { Link } from 'react-router-dom';
import { Result, Button } from 'antd';
import { useI18n } from '@/store/i18n';
import { BasicContext } from '@/store/context';

export default memo(() => {
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  return (
    <Result
      status='404'
      title={t('app.404')}
      subTitle={t('app.404.description')}
      extra={
        <Link to='/'>
          <Button type='primary'>{t('app.404.back')}</Button>
        </Link>
      }
    />
  );
});
