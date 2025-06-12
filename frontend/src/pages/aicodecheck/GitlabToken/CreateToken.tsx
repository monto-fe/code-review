import React, { useContext, useState, useEffect } from 'react';
import { Form, Input, DatePicker, Button, Collapse, Layout, Card, message } from 'antd';
import type { FormInstance } from 'antd';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { queryList, removeData, updateData as updateDataService, createData, getDetail } from './service';
import { TableQueryParam, TableListItem } from './data';
import PromptDrawer from './promptDrawer';
import dayjs from 'dayjs';
import { useNavigate, useLocation } from 'react-router-dom';

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
  const location = useLocation();
  const [form] = Form.useForm();
  const [prompt, setPrompt] = useState('');
  const [promptDrawerVisible, setPromptDrawerVisible] = useState(false);
  const [createSubmitLoading, setCreateSubmitLoading] = useState(false);
  const [updateData, setUpdateData] = useState<Partial<TableListItem>>(initialValues);

  const [prompts, setPrompts] = useState([{
    title: 'Golang 魔法数字检查规则',
    content: `Golang 魔法数字检查规则：
1. 禁止直接使用数值（仅允许 0/1），例如：for i := 0; i < 10; i++ 应改为 const MaxItems = 10，循环条件用 i < MaxItems。
2. 允许的例外场景：数组索引（如 arr[0]）、位运算（如 value & 0xFF）、格式化字符串（如 fmt.Sprintf("%.2f")）。
3. 常量定义建议：const ( RetryCount = 3; CacheTTL = 3600 )。
4. 测试代码也应避免魔法数字，不推荐 assert.Equal(t, 42, result)，推荐 const ExpectedValue = 42。
5. 工具输出建议：第15行的 "42" 是魔法数字，请定义为常量后再使用。`,
  },{
    title: '查找 Diff 中的 Bug',
    content: `作为资深代码审查专家，请仅关注以下 Diff 中可能引入的 Bug，包括但不限于：
- 逻辑错误（条件判断、分支遗漏等）
- 空指针 / 未初始化使用
- 异常处理缺失
- 错误的 API 使用或参数顺序问题
- 数据边界处理不当
- 多线程/异步相关隐患
- 状态或数据未正确更新

### 输入：
只分析下面的代码 diff（Git 风格 +/- 变更），不要考虑上下文以外的文件。

### 输出格式：
- **Bug 描述**：简要说明问题
- **可能影响**：说明可能导致的后果（如崩溃、错误输出等）
- **建议修复方式**：提出修复建议，必要时附简要代码示例

请用简洁、清晰、准确的语言逐条列出你发现的问题。`,
  }]);

  const searchParams = new URLSearchParams(location.search);
  const id = searchParams.get('id');

  const isEdit = !!id;

  useEffect(() => {
    if (id) {
      getDetail(Number(id)).then(res => {
        const detail = res.data || res;
        form.setFieldsValue({
          ...detail,
          expired: detail.expired ? dayjs(detail.expired * 1000) : undefined,
        });
        setPrompt(detail.prompt || '');
        setUpdateData(detail);
      });
    }
  }, [id]);

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
              rules={isEdit ? [] : [{ required: true, message: '名称为必填项' }]}
            >
              <Input placeholder="请输入Token说明" />
            </Form.Item>
            <Form.Item
              label="API"
              name="api"
              rules={isEdit ? [] : [{ required: true, message: 'API为必填项' }]}
            >
              <Input placeholder="请输入API地址" />
            </Form.Item>
            <Form.Item
              label="Token"
              name="token"
              rules={isEdit ? [] : [{ required: true, message: 'Token为必填项' }]}
            >
              <Input placeholder="请输入Token" />
            </Form.Item>
            <Form.Item
              label="有效期"
              name="expired"
              rules={[
                ...(isEdit ? [] : [{ required: true, message: '有效期为必填项' }]),
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
