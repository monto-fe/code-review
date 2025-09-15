import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Descriptions, Tag, Button, Divider, Spin, message, Typography, Input, Form, Select, Collapse, Checkbox, Breadcrumb } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, ReloadOutlined, HomeOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { useContext } from 'react';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { getTaskDetail } from './service';
import { MergeReviewTask } from './data';
import styles from './DetailPage.module.less';

const { TextArea } = Input;
const { Title, Text } = Typography;
const { Option } = Select;

const DetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [taskData, setTaskData] = useState<MergeReviewTask | null>(null);

  // 获取任务详情
  const fetchTaskDetail = async () => {
    if (!id) return;
    
    try {
      const response = await getTaskDetail(id);
      if (response && response.ret_code === 0 && response.data) {
        setTaskData(response.data);
      } else {
        message.error(response?.message || '获取任务详情失败');
      }
    } catch (error) {
      message.error('获取任务详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 刷新数据
  const handleRefresh = async () => {
    setRefreshing(true);
    await fetchTaskDetail();
    setRefreshing(false);
  };

  useEffect(() => {
    fetchTaskDetail();
  }, [id]);

  // 渲染状态标签
  const renderStatusTag = (status: number) => {
    const statusConfig: Record<number, { color: string; text: string }> = {
      0: { color: 'orange', text: '待处理' },
      1: { color: 'blue', text: '处理中' },
      2: { color: 'green', text: '完成' },
      3: { color: 'red', text: '失败' },
    };
    const config = statusConfig[status] || { color: 'default', text: `状态${status}` };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <Spin size="large" />
        <div style={{ marginTop: 16 }}>加载中...</div>
      </div>
    );
  }

  if (!taskData) {
    return (
      <div className={styles.errorContainer}>
        <Text type="danger">任务不存在或已被删除</Text>
        <Button 
          type="primary" 
          onClick={() => navigate('/aicodecheck/merge-review')}
          style={{ marginTop: 16 }}
        >
          返回列表
        </Button>
      </div>
    );
  }

  return (
    <div className={styles.detailPage}>
      {/* 面包屑导航 */}
      <div className={styles.breadcrumb}>
        <Breadcrumb
          items={[
            {
              title: <HomeOutlined />,
              onClick: () => navigate('/home'),
            },
            {
              title: 'AI代码检测',
              onClick: () => navigate('/aicodecheck'),
            },
            {
              title: 'Merge Review',
              onClick: () => navigate('/aicodecheck/merge-review'),
            },
            {
              title: `任务详情 #${taskData.id}`,
            },
          ]}
        />
      </div>

      <div className={styles.header}>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate('/aicodecheck/merge-review')}
          className={styles.backButton}
        >
          返回列表
        </Button>
        
        <div className={styles.headerContent}>
          <Title level={2} className={styles.pageTitle}>
            Merge Review 任务详情
          </Title>
          <Button 
            icon={<ReloadOutlined />} 
            onClick={handleRefresh}
            loading={refreshing}
            className={styles.refreshButton}
          >
            刷新
          </Button>
        </div>
      </div>

      <Row gutter={24}>
        {/* 左侧：任务基本信息 */}
        <Col span={12}>
          <Card title="任务信息" className={styles.infoCard}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="任务ID">
                <Text code>{taskData.id}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="项目名称">
                <Text>{taskData.project_name}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Merge ID">
                <Text code>{taskData.merge_id}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Merge 标题">
                <Text strong>{taskData.merge_title || '未设置'}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="任务状态">
                {renderStatusTag(taskData.status)}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {new Date(taskData.create_time * 1000).toLocaleString('zh-CN')}
              </Descriptions.Item>
              <Descriptions.Item label="更新时间">
                {new Date(taskData.update_time * 1000).toLocaleString('zh-CN')}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>

        {/* 右侧：任务配置信息 */}
        <Col span={12}>
          <Card title="任务配置信息" className={styles.configCard}>
            <div className={styles.configContent}>
              <p className={styles.configDescription}>
                此任务已创建，配置信息如下：
              </p>
              
              <div className={styles.configItem}>
                <strong>Merge 链接：</strong>
                <span>https://gitlab.com/project-{taskData.project_id}/merge_requests/{taskData.merge_id}</span>
              </div>
              
              <div className={styles.configItem}>
                <strong>AI模型：</strong>
                <span>GPT-4 (默认)</span>
              </div>
              
              <div className={styles.configItem}>
                <strong>检测机器人：</strong>
                <div className={styles.botTags}>
                  <Tag color="blue">代码规范检测机器人</Tag>
                  <Tag color="blue">逻辑错误检测机器人</Tag>
                  <Tag color="blue">安全漏洞检测机器人</Tag>
                  <Tag color="blue">性能优化检测机器人</Tag>
                </div>
              </div>
              
              <div className={styles.configItem}>
                <strong>提示词配置：</strong>
                <span>使用默认提示词</span>
              </div>
              
              {taskData.status === 2 && (
                <div className={styles.completionInfo}>
                  <Divider />
                  <Text type="success">任务已完成，检测结果已生成</Text>
                  <div style={{ marginTop: 8 }}>
                    <Button 
                      type="link" 
                      size="small"
                      onClick={() => {
                        // 这里可以添加查看详细检测结果的逻辑
                        message.info('检测结果功能开发中...');
                      }}
                    >
                      查看检测结果
                    </Button>
                  </div>
                </div>
              )}
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default observer(DetailPage); 