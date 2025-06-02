import { Card, List } from 'antd';

const docLinks = [
  { title: '视频指引', url: 'https://example.com/video-guide' },
  { title: 'Gitlab接入手册', url: 'https://example.com/gitlab-guide' },
  { title: 'AI Model配置说明', url: 'https://example.com/ai-model-guide' },
  { title: 'Gitlab Token配置说明', url: 'https://example.com/gitlab-token-guide' },
];

export default function DocLinksCard() {
  return (
      <List
        dataSource={docLinks}
        renderItem={item => (
          <List.Item>
            <a href={item.url} target="_blank" rel="noopener noreferrer">{item.title}</a>
          </List.Item>
        )}
      />
  );
}