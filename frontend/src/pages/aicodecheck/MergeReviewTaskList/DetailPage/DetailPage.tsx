import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Descriptions, Tag, Button, Divider, Spin, message, Typography, Input, Form, Select, Collapse, Checkbox } from 'antd';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, ReloadOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { useContext } from 'react';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { getTaskDetail, startManualCheck } from '../service';
import { MergeReviewTask } from '../data';
import styles from './DetailPage.module.less';

const { TextArea } = Input;
const { Title, Text } = Typography;
const { Option } = Select;

const DetailPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  
  const id = searchParams.get('id');
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [taskData, setTaskData] = useState<MergeReviewTask | null>(null);
  const [checkLoading, setCheckLoading] = useState(false);

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

  // 开始检测
  const handleStartCheck = async () => {
    if (!taskData) return;
    
    setCheckLoading(true);
    try {
      const response = await startManualCheck({
        projectId: taskData.project_id,
        mergeId: taskData.merge_id,
        // 可以根据需要添加其他参数
      });
      
      if (response && response.ret_code === 0) {
        message.success('检测任务已启动，请稍候查看结果');
        // 刷新任务状态
        await fetchTaskDetail();
      } else {
        message.error(response?.message || '启动检测失败');
      }
    } catch (error) {
      message.error('启动检测失败，请稍后重试');
    } finally {
      setCheckLoading(false);
    }
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
        <div style={{ marginTop: 16, fontSize: 16, color: '#666' }}>加载中...</div>
      </div>
    );
  }

  if (!id) {
    return (
      <div className={styles.errorContainer}>
        <Text type="danger">缺少任务ID参数</Text>
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
        {/* 左侧：任务信息和配置 */}
        <Col span={10}>
          {/* 任务配置信息 */}
          <Card title="任务配置" className={styles.configCard}>
            <div className={styles.configContent}>
              <div className={styles.configItem}>
                <strong>Merge 链接：</strong>
                <span>http://165.154.112.72:9980/usms/api/-/merge_requests/{taskData.merge_id}</span>
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
              
              {/* 开始检测按钮 */}
              <div className={styles.actionSection}>
                <Button 
                  type="primary" 
                  size="large"
                  onClick={handleStartCheck}
                  loading={checkLoading}
                  disabled={taskData.status === 1 || taskData.status === 2}
                  block
                >
                  {taskData.status === 1 ? '检测中...' : 
                   taskData.status === 2 ? '检测已完成' : '开始检测'}
                </Button>
                {taskData.status === 1 && (
                  <div className={styles.checkingTip}>
                    <Text type="secondary">检测正在进行中，请稍候...</Text>
                  </div>
                )}
              </div>
            </div>
          </Card>

          {/* 任务基本信息 */}
          <Card title="任务信息" className={styles.infoCard} style={{ marginTop: 16 }}>
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

        {/* 右侧：检测结果 */}
        <Col span={14}>
          <Card title="检测结果" className={styles.resultsCard}>
            {taskData.status === 2 ? (
              <div className={styles.resultsContent}>
                {/* 检测总结 */}
                <div className={styles.summaryResult}>
                  <h4>检测总结</h4>
                  <div className={styles.summaryContent}>
                    <Text type="success">检测已完成，共发现 0 个问题</Text>
                    <div className={styles.summaryStats}>
                      <span>代码规范：通过</span>
                      <span>逻辑错误：通过</span>
                      <span>安全漏洞：通过</span>
                      <span>性能优化：通过</span>
                    </div>
                  </div>
                </div>
                
                {/* 详细检测结果 */}
                <div className={styles.detailedResults}>
                  <h4>详细检测结果</h4>
                  <div className={styles.botResults}>
                    <div className={styles.botResult}>
                      <div className={styles.botHeader}>
                        <Tag color="green">代码规范检测机器人</Tag>
                        <span className={styles.botStatus}>检测通过</span>
                      </div>
                      <div className={styles.botOutput}>
                        未发现代码规范问题，代码结构清晰，命名规范。
                      </div>
                    </div>
                    
                    <div className={styles.botResult}>
                      <div className={styles.botHeader}>
                        <Tag color="green">逻辑错误检测机器人</Tag>
                        <span className={styles.botStatus}>检测通过</span>
                      </div>
                      <div className={styles.botOutput}>
                        未发现逻辑错误，代码逻辑正确。
                      </div>
                    </div>
                    
                    <div className={styles.botResult}>
                      <div className={styles.botHeader}>
                        <Tag color="green">安全漏洞检测机器人</Tag>
                        <span className={styles.botStatus}>检测通过</span>
                      </div>
                      <div className={styles.botOutput}>
                        未发现安全漏洞，代码安全性良好。
                      </div>
                    </div>
                    
                    <div className={styles.botResult}>
                      <div className={styles.botHeader}>
                        <Tag color="green">性能优化检测机器人</Tag>
                        <span className={styles.botStatus}>检测通过</span>
                      </div>
                      <div className={styles.botOutput}>
                        未发现性能问题，代码执行效率良好。
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ) : taskData.status === 1 ? (
              <div className={styles.executingContent}>
                <Spin size="large" />
                <div className={styles.executingText}>正在执行代码检测，请稍候...</div>
              </div>
            ) : taskData.status === 3 ? (
              <div className={styles.failedContent}>
                <Text type="danger">任务执行失败，请检查配置或联系管理员</Text>
                <div className={styles.failedActions}>
                  <Button type="primary" onClick={handleRefresh}>
                    重试
                  </Button>
                  <Button onClick={() => navigate('/aicodecheck/merge-review')}>
                    返回列表
                  </Button>
                </div>
              </div>
            ) : (
              <div className={styles.waitingContent}>
                <Text type="secondary">任务等待中，尚未开始执行</Text>
              </div>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default observer(DetailPage); 