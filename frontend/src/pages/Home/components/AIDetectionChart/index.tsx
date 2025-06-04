import React, { useRef, useState, useEffect } from 'react';
import { Spin, message } from 'antd';
import useEcharts, { EChartsOption } from '@/hooks/useEcharts';
import { getAIDetectionStats } from './service';

const levelMap = [
  '', // level 0 占位
  '完全误导性建议',
  '发现BUG但无实用价值',
  '找到Bug',
  '精准定位并提供建议',
  '超出预期的指导'
];

const pieColors = [
  '#5470C6', // 蓝
  '#91CC75', // 绿
  '#FAC858', // 黄
  '#EE6666', // 红
  '#73C0DE'  // 青
];

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
      // radius: ['50%', '70%'],
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
      let result = [];
      const response = await getAIDetectionStats();
      const { data, ret_code, message: errorMessage } = response;
      if (ret_code === 0) {
        result = data || [];
      } else {
        message.error(errorMessage);
      }
      console.log("result", result);
      const option: EChartsOption = {
        color: pieColors,
        legend: {
          orient: 'vertical',
          left: 'left',
          data: levelMap.slice(1)
        },
        series: [{
          name: 'AI Detection Level',
          type: 'pie',
          avoidLabelOverlap: false,
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: true, fontSize: 20, fontWeight: 'bold' }
          },
          labelLine: { show: false },
          data: result.map((item: any) => ({
            value: item.count,
            name: levelMap[item.level]
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