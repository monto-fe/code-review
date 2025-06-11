import { Card, List } from 'antd';

const docLinks = [
  { title: '视频指引', url: 'https://example.com/video-guide' },
  { title: 'Gitlab接入手册', url: 'https://github.com/monto-fe/code-review/wiki/Gitlab%E6%8E%A5%E5%85%A5%E6%89%8B%E5%86%8C' },
  { title: 'AI Model配置说明', url: '' },
  { title: 'Gitlab Token配置说明', url: 'https://github.com/monto-fe/code-review/wiki/Gitlab-Token%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E' },
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