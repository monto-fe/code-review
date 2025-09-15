import React, { useState, useEffect } from 'react';
import { Form, Input, Select, Button, Card, Row, Col, Divider, message, Spin, Checkbox, Collapse, Tag, Cascader } from 'antd';
import { useNavigate } from 'react-router-dom';
import { ArrowLeftOutlined, SendOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';
import { useContext } from 'react';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { startManualCheck, getAIModelOptions } from './service';
import styles from './CreateForm.module.less';

const { Option } = Select;
const { TextArea } = Input;

interface CreateFormProps {
  onSuccess?: () => void;
}

interface BotOption {
  id: string;
  name: string;
  category: string;
  description: string;
}

interface AIModelOption {
  id: number;
  type: string;
  model: string;
  api_url: string;
  status: number;
  create_time: number;
  update_time: number;
}

const CreateForm: React.FC<CreateFormProps> = ({ onSuccess }) => {
  const navigate = useNavigate();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  
  // AI模型选项，从API获取
  const [aiModels, setAiModels] = useState<AIModelOption[]>([]);
  const [aiModelsLoading, setAiModelsLoading] = useState(false);
  const [cascaderOptions, setCascaderOptions] = useState<any[]>([]);
  
  // 在组件挂载时获取AI模型选项
  useEffect(() => {
    fetchAIModelOptions();
  }, []);
  
  const [availableBots, setAvailableBots] = useState<BotOption[]>([
    { id: 'code-style', name: '代码规范检测机器人', category: '代码规范', description: '检测代码风格、命名规范、注释等' },
    { id: 'logic-error', name: '逻辑错误检测机器人', category: '逻辑错误', description: '检测代码逻辑错误、边界条件等' },
    { id: 'security', name: '安全漏洞检测机器人', category: '安全漏洞', description: '检测SQL注入、XSS等安全漏洞' },
    { id: 'performance', name: '性能优化检测机器人', category: '性能优化', description: '检测性能瓶颈、内存泄漏等' },
  ]);

  const [selectedBots, setSelectedBots] = useState<string[]>(availableBots.map(bot => bot.id));

  // 获取AI模型选项
  const fetchAIModelOptions = async () => {
    setAiModelsLoading(true);
    try {
      const response = await getAIModelOptions();
      if (response && response.ret_code === 0 && response.data) {
        setAiModels(response.data);
        
        // 生成联级选择器选项
        const activeModels = response.data.filter((model: AIModelOption) => model.status === 1);
        const typeGroups = activeModels.reduce((groups: any, model: AIModelOption) => {
          if (!groups[model.type]) {
            groups[model.type] = [];
          }
          groups[model.type].push({
            value: model.model,
            label: model.model || '未知模型',
            isLeaf: true,
            modelData: model
          });
          return groups;
        }, {});
        
        const cascaderOpts = Object.keys(typeGroups).map(type => ({
          value: type,
          label: type || '未知类型',
          children: typeGroups[type]
        }));
        
        setCascaderOptions(cascaderOpts);
        
        // 如果有数据，设置第一个活跃的模型为默认值
        if (activeModels.length > 0 && form) {
          try {
            form.setFieldValue('aiModel', activeModels[0].model);
          } catch (error) {
            console.warn('设置默认AI模型失败:', error);
          }
        }
      } else {
        message.error(response?.message || '获取AI模型失败');
      }
    } catch (error) {
      message.error('获取AI模型失败，请稍后重试');
      console.error('获取AI模型失败:', error);
    } finally {
      setAiModelsLoading(false);
    }
  };

  // 选择机器人
  const handleBotSelection = (botIds: string[]) => {
    setSelectedBots(botIds);
  };

  // 提交表单（现在用于验证和准备数据）
  const handleSubmit = async (values: any) => {
    // 验证Merge链接格式
    const mergeUrl = values.mergeUrl;
    if (!mergeUrl) {
      message.error('请输入Merge链接');
      return;
    }

    // 验证Merge链接格式
    const mergeMatch = mergeUrl.match(/(?:https?:\/\/[^\/]+)\/([^\/]+)\/([^\/]+)\/-\/merge_requests\/(\d+)/);
    if (!mergeMatch) {
      message.error('无法从Merge链接解析项目信息，请检查链接格式。支持格式：http://host:port/group/project/-/merge_requests/id');
      return;
    }

    // 验证通过后，开始检测
    startAICheck();
  };

  // 调用真实AI检测接口
  const startAICheck = async () => {
    if (selectedBots.length === 0) {
      message.warning('请先选择检测机器人');
      return;
    }

    const mergeUrl = form.getFieldValue('mergeUrl');
    if (!mergeUrl) {
      message.error('请输入Merge链接');
      return;
    }

    // 解析Merge链接获取项目信息
    // 支持格式：http://165.154.112.72:9980/usms/api/-/merge_requests/3
    const mergeMatch = mergeUrl.match(/(?:https?:\/\/[^\/]+)\/([^\/]+)\/([^\/]+)\/-\/merge_requests\/(\d+)/);
    if (!mergeMatch) {
      message.error('无法从Merge链接解析项目信息，请检查链接格式');
      return;
    }

    const [, projectGroup, projectName, mergeId] = mergeMatch;
    const projectPath = `${projectGroup}/${projectName}`;
    
    setLoading(true);

    try {
      const response = await startManualCheck({
        projectPath,
        mergeId: parseInt(mergeId),
        selectedBots,
        aiModel: form.getFieldValue('aiModel'),
        // 可以根据需要添加其他参数
      });

      if (response && response.ret_code === 0) {
        message.success('检测任务已启动，请稍候查看结果');
        
        // 清空表单
        form.resetFields();
        setSelectedBots(availableBots.map(bot => bot.id));
        
        // 如果有成功回调，执行它
        if (onSuccess) {
          onSuccess();
        }
      } else {
        message.error(response?.message || '启动检测失败');
      }
    } catch (error) {
      message.error('启动检测失败，请稍后重试');
      console.error('检测失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.createForm}>
      <div className={styles.header}>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate('/aicodecheck/merge-review')}
          className={styles.backButton}
        >
          返回列表
        </Button>
        <h1>创建 Merge Review 任务</h1>
      </div>

      <Row gutter={24}>
        {/* 左侧：输入表单 */}
        <Col span={12}>
          <Card title="任务配置" className={styles.leftCard}>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              initialValues={{
                aiModel: ''
              }}
            >
              <Form.Item
                label="Merge 链接"
                name="mergeUrl"
                rules={[{ required: true, message: '请输入Merge链接' }]}
              >
                                                 <Input 
                  placeholder="请输入Merge Request链接，例如：http://165.154.112.72:9980/usms/api/-/merge_requests/3"
                />
              </Form.Item>



              <Form.Item
                label="AI模型"
                name="aiModel"
                rules={[{ required: true, message: '请选择AI模型' }]}
              >
                <Cascader
                  placeholder="选择AI模型类型和具体模型"
                  loading={aiModelsLoading}
                  options={cascaderOptions}
                  expandTrigger="hover"
                  showSearch={{
                    filter: (inputValue, path) => {
                      return path.some(option => 
                        option.label.toLowerCase().indexOf(inputValue.toLowerCase()) > -1
                      );
                    }
                  }}
                  displayRender={(labels, selectedOptions) => {
                    if (selectedOptions && selectedOptions.length > 0) {
                      const lastOption = selectedOptions[selectedOptions.length - 1];
                      return lastOption?.label || '未知模型';
                    }
                    return labels?.join(' / ') || '请选择模型';
                  }}
                />
              </Form.Item>

              {/* 高级配置 */}
              <Form.Item label="高级配置">
                <Collapse defaultActiveKey={[]} ghost>
                  <Collapse.Panel header="提示词配置" key="prompt">
                    <div className={styles.promptContent}>
                      <p className={styles.promptDescription}>
                        可以自定义每个机器人的检测提示词，以获得更精准的检测结果：
                      </p>
                      
                      {availableBots.map(bot => (
                        <div key={bot.id} className={styles.promptItem}>
                          <div className={styles.promptHeader}>
                            <strong>{bot.name}</strong>
                            <span className={styles.botCategory}>{bot.category}</span>
                          </div>
                          <TextArea
                            placeholder={`${bot.name}的检测提示词...`}
                            rows={3}
                            className={styles.promptTextarea}
                          />
                        </div>
                      ))}
                    </div>
                  </Collapse.Panel>
                  
                  <Collapse.Panel header="选择检测机器人" key="bots">
                    <div className={styles.botSelection}>
                      <p className={styles.description}>
                        选择要使用的检测机器人：
                      </p>
                      
                      <div className={styles.botCheckboxButtons}>
                        {availableBots.map(bot => (
                          <Checkbox
                            key={bot.id}
                            checked={selectedBots.includes(bot.id)}
                            onChange={(e) => {
                              if (e.target.checked) {
                                setSelectedBots([...selectedBots, bot.id]);
                              } else {
                                setSelectedBots(selectedBots.filter(id => id !== bot.id));
                              }
                            }}
                            className={styles.botCheckboxButton}
                          >
                            <div className={styles.botButtonContent}>
                              <strong>{bot.name}</strong>
                              <span className={styles.botCategory}>{bot.category}</span>
                            </div>
                          </Checkbox>
                        ))}
                      </div>
                    </div>
                  </Collapse.Panel>
                </Collapse>
              </Form.Item>

              {/* 开始检测按钮 */}
              <Form.Item>
                <Button 
                  type="primary" 
                  onClick={startAICheck}
                  loading={loading}
                  disabled={selectedBots.length === 0}
                  icon={<SendOutlined />}
                  block
                >
                  开始检测
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        {/* 右侧：检测状态 */}
        <Col span={12}>
          <Card title="检测状态" className={styles.rightCard}>
            <Spin spinning={loading}>
              {selectedBots.length > 0 && (
                <div className={styles.outputContent}>
                  {loading ? (
                    <div className={styles.checkingStatus}>
                      <div className={styles.checkingIcon}>
                        <SendOutlined style={{ fontSize: 48, color: '#1890ff' }} />
                      </div>
                      <h3>正在启动检测任务...</h3>
                      <p>请稍候，系统正在处理您的请求</p>
                      <div className={styles.checkingSteps}>
                        <div className={styles.stepItem}>
                          <span className={styles.stepNumber}>1</span>
                          <span className={styles.stepText}>解析Merge链接</span>
                          <span className={styles.stepStatus}>完成</span>
                        </div>
                        <div className={styles.stepItem}>
                          <span className={styles.stepNumber}>2</span>
                          <span className={styles.stepText}>启动AI检测</span>
                          <span className={styles.stepStatus}>进行中</span>
                        </div>
                        <div className={styles.stepItem}>
                          <span className={styles.stepNumber}>3</span>
                          <span className={styles.stepText}>等待检测完成</span>
                          <span className={styles.stepStatus}>等待中</span>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className={styles.readyStatus}>
                      <div className={styles.readyIcon}>
                        <SendOutlined style={{ fontSize: 48, color: '#52c41a' }} />
                      </div>
                      <h3>准备就绪</h3>
                      <p>点击左侧的"开始检测"按钮启动代码审核任务</p>
                      <div className={styles.readyInfo}>
                        <p><strong>已选择 {selectedBots.length} 个检测机器人：</strong></p>
                        <div className={styles.selectedBots}>
                          {selectedBots.map(botId => {
                            const bot = availableBots.find(b => b.id === botId);
                            return (
                              <Tag key={botId} color="blue" className={styles.botTag}>
                                {bot?.name}
                              </Tag>
                            );
                          })}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </Spin>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default observer(CreateForm); 