import { useContext, useState } from 'react';
import { Drawer, Button } from 'antd';
import ReactMarkdown from 'react-markdown';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { gitlabConfigDoc } from './docs';

const IP = import.meta.env.VITE_APP_APIHOST;

function DocDrawer() {
  const [open, setOpen] = useState(false);
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  return (
    <>
      <Button type='link' onClick={() => setOpen(true)}>
        {t('app.global.doc')}
      </Button>
      <Drawer
        title='📡 Config Webhook（GitLab）'
        placement='right'
        width={1080}
        onClose={() => setOpen(false)}
        open={open}
      >
        <ReactMarkdown>{gitlabConfigDoc(t, IP)}</ReactMarkdown>
      </Drawer>
    </>
  );
}

export default DocDrawer;
