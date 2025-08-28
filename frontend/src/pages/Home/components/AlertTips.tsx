import { useState, useEffect, useContext, ReactNode } from 'react';
import { Alert, Space, Typography } from 'antd';

import aiConfig from '@/config/aiConfig';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import DocDrawer from './DocDrawer';

const { Paragraph } = Typography;
const IP = import.meta.env.VITE_APP_APIHOST;

const DOC_URL = 'https://github.com/monto-fe/code-review/wiki';
const WEBHOOK_URL = `${IP}/webhook/merge`;

interface AlertTipsProps {
  AIConfig: boolean;
  GitlabConfig: boolean;
}

function AlertTips({ AIConfig, GitlabConfig }: AlertTipsProps) {
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [closed, setClosed] = useState(false);

  useEffect(() => {
    if (AIConfig && GitlabConfig) {
      setClosed(localStorage.getItem(aiConfig.home_alert_close_flag) === 'true');
    } else {
      setClosed(false); // 配置未完成时始终显示
    }
  }, [AIConfig, GitlabConfig]);

  if (closed) return null;

  let message: ReactNode = '';
  let closable = false;
  let type: 'success' | 'error' = 'error';

  if (AIConfig && GitlabConfig) {
    message = (
      <Space size='middle'>
        {t('page.home.text-alert-tips-success-gitlab-hook')}
        <Paragraph style={{ display: 'inline-block', marginBottom: 0 }} copyable={{ text: WEBHOOK_URL }}>
          {WEBHOOK_URL}
        </Paragraph>
        。
        <DocDrawer />
      </Space>
    );
    closable = true;
    type = 'success';
  } else if (!GitlabConfig) {
    message = (
      <Space size='middle'>
        {t('page.home.text-alert-tips-failed-gitlab-config')}
        <a href={DOC_URL} target='_blank' rel='noopener noreferrer'>
          {t('app.global.doc')}
        </a>
      </Space>
    );
  } else {
    message = (
      <Space size='middle'>
        {t('page.home.text-alert-tips-failed-all-gitlab-config')}
        <a href={DOC_URL} target='_blank' rel='noopener noreferrer'>
          {t('app.global.doc')}
        </a>
      </Space>
    );
  }

  return (
    <Alert
      message={message}
      type={type}
      showIcon
      closable={closable}
      onClose={() => {
        if (AIConfig && GitlabConfig) {
          localStorage.setItem(aiConfig.home_alert_close_flag, 'true');
          setClosed(true);
        }
      }}
      className='mb-24'
    />
  );
}

export default AlertTips;
