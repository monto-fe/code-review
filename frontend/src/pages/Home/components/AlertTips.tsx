import { useState, useEffect } from 'react';
import { Alert } from 'antd';

const DOC_URL = 'https://example.com/ai-gitlab-config-doc';

export default function AlertTips({ AIConfig, GitlabConfig }: { AIConfig: boolean, GitlabConfig: boolean }) {
  const [closed, setClosed] = useState(false);

  useEffect(() => {
    if (AIConfig && GitlabConfig) {
      setClosed(localStorage.getItem('alertTipsClosed') === 'true');
    } else {
      setClosed(false); // 配置未完成时始终显示
    }
  }, [AIConfig, GitlabConfig]);

  if (closed) return null;

  let message: React.ReactNode = '';
  let closable = false;
  let type: 'success' | 'error' = 'error';

  if (AIConfig && GitlabConfig) {
    message = (
      <span>
        恭喜您配置成功，Gitlab接入地址为：<b>http://ip:9000/webhook</b>，
        <a href={DOC_URL} target="_blank" rel="noopener noreferrer">参考文档</a>
      </span>
    );
    closable = true;
    type = 'success';
  } else if (!GitlabConfig) {
    message = (
      <span>
        请在下面表单配置Gitlab Token后使用，
        <a href={DOC_URL} target="_blank" rel="noopener noreferrer">参考文档</a>
      </span>
    );
  } else {
    message = (
      <span>
        请在下面表单配置AI Model和Gitlab Token后使用，
        <a href={DOC_URL} target="_blank" rel="noopener noreferrer">参考文档</a>
      </span>
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
          localStorage.setItem('alertTipsClosed', 'true');
          setClosed(true);
        }
      }}
      style={{ marginBottom: 24 }}
    />
  );
}