import { memo } from 'react';
import { Spin } from 'antd';

export default memo(() => (
  <div style={{
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '100%',
    height: '100%',
    minHeight: '100vh'
  }}>
    <Spin size='large' />
  </div>
));
