import React from 'react';
import { Collapse, List } from 'antd';

const { Panel } = Collapse;

const data = [
  {
    summary: 'v0.0.4 支持RAG能力',
    features: ['读取Git项目，通过RAG能力进行代码审核'],
    extra: '2025-04-05'
  },
  {
    summary: 'v0.0.3 优化Git项目读取逻辑',
    features: ['优化服务拉取Git项目缓存'],
    extra: '2025-04-05'
  },
  {
    summary: 'v0.0.2 实现快速部署',
    features: ['实现AI代码检测服务快速部署', '支持WebHook推送检测结果'],
    extra: '2025-03-15'
  },
  {
    summary: 'v0.0.1 初始版本发布',
    features: ['支持针对Git Merge代码进行审核'],
    extra: '2025-03-01'
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