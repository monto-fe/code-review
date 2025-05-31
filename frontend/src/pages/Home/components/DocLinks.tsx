import React from 'react';
import { LinkOutlined } from '@ant-design/icons';
import ALink from '@/components/ALink';

const DocLinks = () => {
  return (
    <div>
      <ALink to='https://docs.example.com'>
        <LinkOutlined /> 查看使用文档
      </ALink>
    </div>
  );
};

export default DocLinks; 