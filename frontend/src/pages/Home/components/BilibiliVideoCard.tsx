import { Card } from 'antd';

export default function BilibiliVideoCard() {
  return (
      <div style={{ position: 'relative', paddingBottom: '56.25%', height: 0, overflow: 'hidden' }}>
        <iframe
          src=""
          scrolling="no"
          width="100%"
          height={575}
          frameBorder="0"
          allowFullScreen
          title="配置视频教程"
        />
      </div>
  );
}