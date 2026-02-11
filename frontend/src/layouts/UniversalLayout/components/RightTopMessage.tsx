import React, { memo } from 'react';
import { Popover } from 'antd';
import wechatImg from '@/assets/images/wechat.png';

function RightTopMessage() {
  return (
    <Popover
    content={
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
        <img
          src={wechatImg}
          alt='微信群二维码'
          style={{ display: 'block', width: 100, height: 100, objectFit: 'contain' }}
        />
      </div>
    }
  >
    <div className='universallayout-top-notocemenu ant-dropdown-link cursor' onClick={(e: React.MouseEvent<HTMLDivElement>) => e.preventDefault()}>
      问题反馈
    </div>
    </Popover>
  );
}

export default memo(RightTopMessage);
