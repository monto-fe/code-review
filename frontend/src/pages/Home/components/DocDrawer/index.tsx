import React, { useState } from 'react';
import { Drawer, Button } from 'antd';
import ReactMarkdown from 'react-markdown';
import { gitlabConfigDoc } from './docs';
import './index.less';

const IP = import.meta.env.VITE_APP_APIHOST
const domain = import.meta.env.VITE_APP_DOMAIN || window.location.origin

const DocDrawer: React.FC = () => {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button type="link" onClick={() => setOpen(true)}>
        参考文档
      </Button>
      <Drawer
        title="📡 配置 Webhook（GitLab）"
        placement="right"
        width={1080}
        onClose={() => setOpen(false)}
        open={open}
      >
        <div className="doc-drawer-content">
          <ReactMarkdown>{gitlabConfigDoc(IP, domain)}</ReactMarkdown>
        </div>
      </Drawer>
    </>
  );
};

export default DocDrawer;