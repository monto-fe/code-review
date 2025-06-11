import React, { useState, useEffect } from 'react';
import { Spin, Statistic, message } from 'antd';
import { ResponseData } from '@/utils/request';
import { getMergeStats } from './service';

const MergeProblemStats = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [stats, setStats] = useState<number>(0);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response: ResponseData<{ count: number }> = await getMergeStats({passed: 1});
        const { data, ret_code, message: errorMessage } = response;
        if (ret_code === 0) {
          setStats(data?.count || 0);
        } else {
          message.error(errorMessage);
        }
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

export default MergeProblemStats; 