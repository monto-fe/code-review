import React, { useState, useEffect } from 'react';
import { List, Card, Tag, Spin, Typography } from 'antd';
import { getRecentRecords } from './service';

const { Text } = Typography;

interface RecordItem {
  id: string;
  type: 'commit' | 'review';
  title: string;
  author: string;
  status: 'pending' | 'approved' | 'rejected';
  createdAt: string;
  url: string;
}

const RecentRecords = () => {
  const [loading, setLoading] = useState(false);
  const [records, setRecords] = useState<RecordItem[]>([]);

  useEffect(() => {
    const fetchRecords = async () => {
      setLoading(true);
      try {
        const response = await getRecentRecords();
        setRecords(response.data || []);
      } catch (error) {
        console.error('Failed to fetch recent records:', error);
      }
      setLoading(false);
    };

    fetchRecords();
  }, []);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved':
        return 'success';
      case 'rejected':
        return 'error';
      default:
        return 'warning';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'approved':
        return '已通过';
      case 'rejected':
        return '已拒绝';
      default:
        return '待审核';
    }
  };

  return (
    <Card title="最近提交和审查记录" className="recent-records">
      <Spin spinning={loading}>
        <List
          dataSource={records}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                title={
                  <a href={item.url} target="_blank" rel="noopener noreferrer">
                    {item.title}
                  </a>
                }
                description={
                  <>
                    <Text type="secondary">
                      {item.author}
                    </Text>
                    <Tag color={getStatusColor(item.status)} style={{ marginLeft: 8 }}>
                      {getStatusText(item.status)}
                    </Tag>
                  </>
                }
              />
            </List.Item>
          )}
        />
      </Spin>
    </Card>
  );
};

export default RecentRecords; 