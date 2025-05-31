import React from 'react';
import { Collapse, List } from 'antd';

const { Panel } = Collapse;

const data = [
  {
    summary: 'v1.0.0 初始版本发布',
    features: ['功能A', '功能B', '功能C'],
    extra: <div>2023-10-01</div>
  },
  {
    summary: 'v1.1.0 新增功能A',
    features: ['功能D', '功能E'],
    extra: '2023-10-15'
  },
  {
    summary: 'v1.2.0 修复Bug',
    features: ['修复X问题', '优化Y性能'],
    extra: '2023-11-01'
  },
];

const UpdateLog = () => (
  <Collapse bordered={false} ghost expandIconPosition={'start'}>
    {data.map((item, idx) => (
      <Panel
        header={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>{item.summary}</span>
            <span>{item.extra}</span>
          </div>
        }
        key={idx}
      >
        <List
          size="small"
          dataSource={item.features}
          renderItem={feature => <List.Item>{feature}</List.Item>}
        />
      </Panel>
    ))}
  </Collapse>
);

export default UpdateLog; 