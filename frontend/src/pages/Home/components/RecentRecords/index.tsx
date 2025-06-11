import React, { useState, useEffect } from 'react';
import { Card, Tag, Spin, Typography, message, Table } from 'antd';
import dayjs from 'dayjs';
import { getRecentRecords } from './service';
import ResultParagraph from './ResultParagraph';

const { Text, Paragraph } = Typography;

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

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '项目',
      dataIndex: 'merge_url',
      key: 'merge_url',
      width: 120,
      render: (text: string, record: any) => (
        <a href={record.url} target="_blank" rel="noopener noreferrer">
          {record.project_name}
        </a>
      ),
    },
    {
      title: 'result',
      dataIndex: 'result',
      key: 'result',
      render: (text: string) => {
        return <ResultParagraph text={text} />;
      },
    },
    {
      title: '时间',
      dataIndex: 'create_time',
      key: 'create_time',
      width: 120,
      render: (item: any) => {
       return dayjs(item*1000).format('YYYY-MM-DD HH:mm:ss')
      },
    },
  ];

  useEffect(() => {
    const fetchRecords = async () => {
      setLoading(true);
      try {
        const response = await getRecentRecords();
        console.log("recent records response", response);
        const { data, ret_code, message: errorMessage } = response;
        if (ret_code === 0) {
          setRecords(data.data || []);
        } else {
          message.error(errorMessage);
        }
      } catch (error) {
        setRecords([]);
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
        <Table
          columns={columns}
          dataSource={records}
          rowKey="id"
          pagination={{ pageSize: 10 }}
          scroll={{ x: true }}
        />
      </Spin>
    </Card>
  );
};

export default RecentRecords; 