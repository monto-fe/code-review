import React, { useEffect, useState } from 'react';
import { Card, Descriptions, Spin, message } from 'antd';
import { useLocation } from 'react-router-dom';
import { getDetail } from './service';
import { TableListItem } from './data';
import dayjs from 'dayjs';

const GitlabTokenDetail: React.FC = () => {
  const location = useLocation();
  const searchParams = new URLSearchParams(location.search);
  const id = searchParams.get('id');
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<TableListItem | null>(null);

  console.log("id", id);
  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getDetail(Number(id))
      .then((res) => {
        const { data } = res;
        setData(data);
      })
      .catch(() => {
        message.error('获取详情失败');
      })
      .finally(() => setLoading(false));
  }, [id]);

  return (
    <Card title="Gitlab Token 详情" style={{ maxWidth: 800, margin: '0 auto' }}>
      <Spin spinning={loading}>
        {data && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="ID">{data.id}</Descriptions.Item>
            <Descriptions.Item label="API">{data.api}</Descriptions.Item>
            <Descriptions.Item label="Token">{data.token}</Descriptions.Item>
            <Descriptions.Item label="有效期">--</Descriptions.Item>
            <Descriptions.Item label="Webhook名称">{data.webhook_name}</Descriptions.Item>
            <Descriptions.Item label="Webhook地址">{data.webhook_url}</Descriptions.Item>
            <Descriptions.Item label="源分支">{data.source_branch}</Descriptions.Item>
            <Descriptions.Item label="目标分支">{data.target_branch}</Descriptions.Item>
            <Descriptions.Item label="Webhook状态">{data.webhook_status === 1 ? '开启' : '关闭'}</Descriptions.Item>
            <Descriptions.Item label="检测规则状态">{data.rule_check_status === 1 ? '开启' : '关闭'}</Descriptions.Item>
            <Descriptions.Item label="状态">{data.status === 1 ? '启用' : '禁用'}</Descriptions.Item>
          </Descriptions>
        )}
      </Spin>
    </Card>
  );
};

export default GitlabTokenDetail;
