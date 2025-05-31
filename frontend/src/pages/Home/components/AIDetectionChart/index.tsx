import React, { useRef, useState, useEffect } from 'react';
import { Spin } from 'antd';
import useEcharts, { EChartsOption } from '@/hooks/useEcharts';
import { getAIDetectionStats } from './service';

const pieChartOption: EChartsOption = {
  tooltip: {
    trigger: 'item',
    formatter: '{a} <br/>{b}: {c} ({d}%)'
  },
  legend: {
    orient: 'vertical',
    left: 'left',
    data: ['Level 1', 'Level 2', 'Level 3', 'Level 4', 'Level 5']
  },
  series: [
    {
      name: 'AI Detection Level',
      type: 'pie',
      radius: ['50%', '70%'],
      avoidLabelOverlap: false,
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '20',
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: [
        { value: 0, name: 'Level 1' },
        { value: 0, name: 'Level 2' },
        { value: 0, name: 'Level 3' },
        { value: 0, name: 'Level 4' },
        { value: 0, name: 'Level 5' }
      ]
    }
  ]
};

const AIDetectionChart = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const chartRef = useRef<HTMLDivElement>(null);

  useEcharts(chartRef, pieChartOption, async (chart) => {
    setLoading(true);
    try {
      const response = await getAIDetectionStats();
      const data = response.data || [];
      
      const option: EChartsOption = {
        series: [{
          data: data.map((item: any) => ({
            value: item.value,
            name: `Level ${item.level}`
          }))
        }]
      };
      
      chart.setOption(option);
    } catch (error) {
      console.error('Failed to fetch AI detection stats:', error);
    }
    setLoading(false);
  });

  return (
    <Spin spinning={loading}>
      <div ref={chartRef} style={{ height: '300px' }} />
    </Spin>
  );
};

export default AIDetectionChart; 