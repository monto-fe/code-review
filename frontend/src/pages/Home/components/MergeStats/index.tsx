import React, { useState, useEffect } from 'react';
import { Spin, Statistic } from 'antd';
import { ResponseData } from '@/utils/request';
import { getMergeStats } from './service';

const MergeStats = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [stats, setStats] = useState<number>(0);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response: ResponseData<{ total: number }> = await getMergeStats();
        setStats(response.data?.total || 0);
      } catch (error) {
        console.error('Failed to fetch merge stats:', error);
      }
      setLoading(false);
    };

    fetchData();
  }, []);

  return (
    <Spin spinning={loading}>
      <Statistic
        value={stats}
        suffix="次"
        valueStyle={{ fontSize: '24px', fontWeight: 'bold' }}
      />
    </Spin>
  );
};

export default MergeStats; 