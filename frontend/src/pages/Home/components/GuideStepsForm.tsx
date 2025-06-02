import { useState, useEffect, useCallback } from 'react';
import dayjs from 'dayjs';
import { Steps, Form, Button, message, Cascader, Input, DatePicker } from 'antd';
import { createData as saveAIModelConfig } from '@/pages/aicodecheck/AIModelManager/service';
import { queryInnerAIModel } from '../service';
import { createData as saveGitlabConfig } from '@/pages/aicodecheck/GitlabToken/service';

const { Step } = Steps;

const generateAIModelOptions = (data: any[]) => {
  const uniqueTypes = [...new Set(data.map(item => item.type))];
  const options = uniqueTypes.map(type => ({
    label: type,
    value: type,
    children: data.filter(item => item.type === type).map(item => ({
      label: item.model,
      value: item.model,
    })),
  }));

  return options;
}

const GuideStepsForm = ({ AIConfig, GitlabConfig, callback }: { AIConfig: any[], GitlabConfig: any[], callback: (key: string) => void }) => {
  const [current, setCurrent] = useState(0);
  const [aiForm] = Form.useForm();
  const [gitlabForm] = Form.useForm();
  const [AIModelList, setAIModelList] = useState<any[]>([]);
  const [AIModelOptions, setAIModelOptions] = useState<any[]>([]);
  const [nextLoading, setNextLoading] = useState(false);
  const [modalLoading, setModalLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);

  const getAIModelList = useCallback(() => {
    queryInnerAIModel().then((res) => {
      const { data } = res;
      setAIModelList(data);
      setAIModelOptions(generateAIModelOptions(data));
    });
  }, [queryInnerAIModel]);

  const getApiUrl = useCallback((type: string[]) => {
    return AIModelList.find(item => item.type === type[0] && item.model === type[1])?.api_url;
  }, [AIModelList]);

  useEffect(() => {
    getAIModelList();
  }, [getAIModelList]);

  // 根据配置数组长度决定显示哪些步骤
  const getVisibleSteps = () => {
    const steps = [];
    
    if (AIConfig.length === 0) {
      steps.push({
        title: 'AI Model 配置',
        content: (
          <Form form={aiForm} layout="vertical">
            <Form.Item name="type" label="模型类型" rules={[{ required: true, message: '请选择模型类型' }]}>
              <Cascader
                options={AIModelOptions}
                placeholder="请选择AI模型"
              />
            </Form.Item>
            <Form.Item name="api_key" label="API Key" rules={[{ required: true, message: '请输入API Key' }]}><Input /></Form.Item>
          </Form>
        ),
      });
    }

    if (GitlabConfig.length === 0) {
      steps.push({
        title: 'Gitlab Token 配置',
        content: (
          <Form form={gitlabForm} layout="vertical" onFinish={async (values) => {
            try {
              setModalLoading(true);
              const { api, token, expired } = values;
              await saveGitlabConfig({
                api,
                token,
                expired: expired.unix(),
                status: 1,
            });
              message.success('Gitlab Token 配置保存成功');
              callback('GitlabConfig');
              setModalVisible(true);
            } catch {
              message.error('保存失败');
            } finally {
              setModalLoading(false);
            }
          }}>
            <Form.Item name="api" label="Gitlab API" rules={[{ required: true, message: '请输入Gitlab地址' }]}>
                <Input placeholder='https://gitlab.com/api/v4'/>
            </Form.Item>
            <Form.Item
              name="token"
              label="Gitlab Token"
              rules={[
                { required: true, message: '请输入Gitlab Token' },
              ]}
            >
              <Input placeholder='glpat-1234567890'/>
            </Form.Item>
            <Form.Item
              name="expired"
              label="有效期"
              required
              rules={[
                {
                  validator: (_, value) => {
                    if (!value) {
                      return Promise.reject(new Error('请选择有效期'));
                    }
                    if (value.isBefore(dayjs().startOf('day'))) {
                      return Promise.reject(new Error('有效期不能早于今天'));
                    }
                    return Promise.resolve();
                  }
                }
              ]}
            >
              <DatePicker
                style={{ width: '100%' }}
                disabledDate={current => current && current < dayjs().startOf('day')}
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={modalLoading}>
                保存
              </Button>
            </Form.Item>
          </Form>
        ),
      });
    }

    return steps;
  };

  const visibleSteps = getVisibleSteps();

  const next = async () => {
    if (current === 0 && AIConfig.length === 0) {
      try {
        const { api_key, type } = await aiForm.validateFields();
        const apiUrl = getApiUrl(type);
        setNextLoading(true);
        await saveAIModelConfig({
          api_key,
          api_url: apiUrl,
          type: type[0],
          model: type[1],
          name: type[1],
          is_active: 1
        });
        message.success('AI Model 配置保存成功');
        callback('AIConfig');
        setNextLoading(false);
        setCurrent(current + 1);
      } catch {
        setNextLoading(false);
      }
    }
  };

  // 如果没有需要配置的步骤，直接返回null
  if (visibleSteps.length === 0) {
    return null;
  }

  return (
    <div>
      <Steps type="navigation" current={current} style={{ marginBottom: 24 }}>
        {visibleSteps.map(item => (<Step key={item.title} title={item.title} />))}
      </Steps>
      {visibleSteps[current].content}
      <div style={{ marginTop: 24 }}>
        {current < visibleSteps.length - 1 && (
          <Button type="primary" onClick={next} loading={nextLoading}>下一步</Button>
        )}
      </div>
    </div>
  );
};

export default GuideStepsForm; 