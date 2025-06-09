import React, { useContext, useState, useEffect } from 'react';
import { Form, Input, DatePicker, Button, Collapse, Layout, Card, message } from 'antd';
import type { FormInstance } from 'antd';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { queryList, removeData, updateData as updateDataService, createData } from './service';
import { TableQueryParam, TableListItem } from './data';
import PromptDrawer from './promptDrawer';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';

const { TextArea } = Input;
const { Sider, Content } = Layout;
const { Panel } = Collapse;

const MAX_PROMPT_LENGTH = 1000;

// 禁用3天内的日期
const disabledDate = (current: dayjs.Dayjs) => {
  return current && current < dayjs().add(3, 'day').startOf('day');
};

interface AIConfigPageProps {
  initialValues?: Partial<TableListItem>;
  visible?: boolean;
  setVisible?: (visible: boolean) => void;
  onSubmit?: (values: TableListItem, form: FormInstance) => void;
  onSubmitLoading?: boolean;
}

const AIConfigPage: React.FC<AIConfigPageProps> = ({ 
  initialValues = {}, 
  visible = false, 
  setVisible,
  onSubmit: externalOnSubmit,
  onSubmitLoading = false 
}) => {
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [prompt, setPrompt] = useState('');
  const [promptDrawerVisible, setPromptDrawerVisible] = useState(false);
  const [createSubmitLoading, setCreateSubmitLoading] = useState(false);
  const [updateData, setUpdateData] = useState<Partial<TableListItem>>(initialValues);

  const [prompts, setPrompts] = useState([{
    title: '推荐提示词',
    content: '请输入自定义提示词，支持Markdown，最多1000字',
  },{
    title: 'AI提示词',
    content: '请输入自定义提示词，支持Markdown，最多1000字',
  }]);

  // useEffect(() => {
  //   if (initialValues) {
  //     setUpdateData(initialValues);
  //     form.setFieldsValue({
  //       ...initialValues,
  //       expired: initialValues.expired ? dayjs(initialValues.expired * 1000) : undefined
  //     });
  //   }
  // }, [initialValues, form]);

  const handleFinish = (values: any) => {
    console.log("values", values)
    if (prompt.length > MAX_PROMPT_LENGTH) {
      message.error('自定义提示词最多1000字');
      return;
    }
    onSubmit({ 
      ...values, 
      prompt, 
      expired: values.expired?.unix(),
      status: 1
    }, form);
  };

  const handleUsePrompt = (content: string) => {
    setPrompt(content);
    setPromptDrawerVisible(false);
  };

  const onSubmit = (values: TableListItem, form: FormInstance) => {
    if (externalOnSubmit) {
      externalOnSubmit(values, form);
      return;
    }

    setCreateSubmitLoading(true);
    const request = updateData.id ? updateDataService : createData;
    request({ ...values, id: updateData.id as number })
      .then(() => {
        // form.resetFields();
        message.success(values.id ? t('app.global.tip.update.success') : t('app.global.tip.create.success'));
        // 跳转到GitlabToken列表页
        navigate('/aicodecheck/GitlabToken');
        setCreateSubmitLoading(false);
      })
      .catch(() => {
        setCreateSubmitLoading(false);
      });
  };

  return (
    <Layout style={{ background: '#fff', minHeight: 600 }}>
      <Content style={{ padding: 24, flex: 1 }}>
        <Form
          form={form}
          layout="horizontal"
          labelCol={{ span: 6 }}
          wrapperCol={{ span: 16 }}
          initialValues={updateData}
          onFinish={handleFinish}
        >
          <Card title="基础配置" bordered={false} style={{ marginBottom: 16 }}>
            <Form.Item
              label="名称"
              name="name"
              rules={[{ required: true, message: '名称为必填项' }]}
            >
              <Input placeholder="请输入Token说明" />
            </Form.Item>
            <Form.Item
              label="API"
              name="api"
              rules={[{ required: true, message: 'API为必填项' }]}
            >
              <Input placeholder="请输入API地址" />
            </Form.Item>
            <Form.Item
              label="Token"
              name="token"
              rules={[{ required: true, message: 'Token为必填项' }]}
            >
              <Input placeholder="请输入Token" />
            </Form.Item>
            <Form.Item
              label="有效期"
              name="expired"
              rules={[
                { required: true, message: '有效期为必填项' },
                {
                  validator: async (_, value) => {
                    if (value && value < dayjs().add(3, 'day').startOf('day')) {
                      throw new Error('有效期必须大于3天');
                    }
                  }
                }
              ]}
            >
              <DatePicker 
                style={{ width: '100%' }} 
                placeholder="请选择有效期" 
                disabledDate={disabledDate}
              />
            </Form.Item>
          </Card>

          <Collapse bordered={false} style={{ marginBottom: 16 }}>
            <Panel header="通知配置（可选）" key="notify">
              <Form.Item
                label="Webhook名称"
                name="webhook_name"
                tooltip="企业微信机器人名称或webhook名称"
              >
                <Input placeholder="请输入Webhook名称" />
              </Form.Item>
              <Form.Item
                label="Webhook地址"
                name="webhook_url"
                tooltip="企业微信机器人地址或webhook地址"
              >
                <Input placeholder="请输入Webhook地址" />
              </Form.Item>
            </Panel>
            <Panel header="检测配置（可选）" key="detect">
              <Form.Item
                label="源分支"
                name="source_branch"
                tooltip="触发源分支，不填写则匹配全部"
              >
                <Input placeholder="请输入源分支" />
              </Form.Item>
              <Form.Item
                label="目标分支"
                name="target_branch"
                tooltip="触发目标分支，不填写则匹配全部"
              >
                <Input placeholder="请输入目标分支" />
              </Form.Item>
            </Panel>
          </Collapse>

          <Form.Item wrapperCol={{ offset: 6 }}>
            <Button type="primary" htmlType="submit" loading={onSubmitLoading || createSubmitLoading}>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Content>
      <Sider width={600} style={{ background: '#fafbfc', padding: 24, borderLeft: '1px solid #eee' }}>
        <Card
          title="自定义提示词"
          style={{ height: '100%' }}
          extra={<Button type="primary" onClick={() => setPromptDrawerVisible(true)}>推荐提示词</Button>}
        >
          <TextArea
            value={prompt}
            onChange={e => setPrompt(e.target.value)}
            placeholder="请输入自定义提示词，支持Markdown，最多1000字"
            autoSize={{ minRows: 10, maxRows: 18 }}
            maxLength={MAX_PROMPT_LENGTH}
            showCount
          />
          <div style={{ color: '#aaa', marginTop: 8 }}>
            支持Markdown语法，最多1000字
          </div>
        </Card>
      </Sider>
      <PromptDrawer
        visible={promptDrawerVisible}
        onClose={() => setPromptDrawerVisible(false)}
        prompts={prompts}
        onUse={handleUsePrompt}
      />
    </Layout>
  );
};

export default AIConfigPage;
