import { useState, useEffect } from 'react';
import { Steps, Form, Input, Button, message, Modal, Select } from 'antd';
// 假设你有如下API
import { createData as saveAIModelConfig, MODEL_TYPE_OPTIONS } from '@/pages/aicodecheck/AIModelManager/service';
// import { getGitlabConfig, saveGitlabConfig } from '@/pages/aicodecheck/GitlabToken/service';

const { Step } = Steps; 

const GuideStepsForm = () => {
  const [current, setCurrent] = useState(0);
  const [aiForm] = Form.useForm();
  const [gitlabForm] = Form.useForm();
  const [modalVisible, setModalVisible] = useState(false);

  // 1. 获取初始值
  useEffect(() => {
    // getAIModelConfig().then(data => aiForm.setFieldsValue(data));
    // getGitlabConfig().then(data => gitlabForm.setFieldsValue(data));
  }, []);

  // 2. 步骤内容
  const steps = [
    {
      title: 'AI Model 配置',
      content: (
        <Form form={aiForm} layout="vertical">
          <Form.Item name="type" label="模型类型" rules={[{ required: true, message: '请选择模型类型' }]}>
            <Select options={MODEL_TYPE_OPTIONS} />
          </Form.Item>
          <Form.Item name="name" label="模型名称" rules={[{ required: true, message: '请输入模型名称' }]}><Input /></Form.Item>
          <Form.Item name="model" label="模型" rules={[{ required: true, message: '请选择模型' }]}><Input /></Form.Item>
          <Form.Item name="api_url" label="API Base URL" rules={[{ required: true, message: '请输入API Base URL' }]}><Input /></Form.Item>
          <Form.Item name="api_key" label="API Key" rules={[{ required: true, message: '请输入API Key' }]}><Input /></Form.Item>
        </Form>
      ),
    },
    {
      title: 'Gitlab Token 配置',
      content: (
        <Form form={gitlabForm} layout="vertical">
          <Form.Item name="token" label="Gitlab Token" rules={[{ required: true, message: '请输入Gitlab Token' }]}><Input /></Form.Item>
          <Form.Item name="gitlabUrl" label="Gitlab 地址" rules={[{ required: true, message: '请输入Gitlab地址' }]}><Input /></Form.Item>
          <Form.Item name="projectId" label="项目ID" rules={[{ required: true, message: '请输入项目ID' }]}><Input /></Form.Item>
        </Form>
      ),
    },
  ];

  // 3. 下一步
  const next = async () => {
    if (current === 0) {
      try {
        const values = await aiForm.validateFields();
        console.log("values", values);
        // await saveAIModelConfig(values);
        message.success('AI Model 配置保存成功');
        setCurrent(current + 1);
      } catch {}
    } else if (current === 1) {
      try {
        const values = await gitlabForm.validateFields();
        // await saveGitlabConfig(values);
        message.success('Gitlab Token 配置保存成功');
        setModalVisible(true);
      } catch {}
    }
  };

  // 4. 上一步
  const prev = () => setCurrent(current - 1);

  // 5. 完成
  const handleFinish = () => {
    setModalVisible(false);
    setCurrent(0);
  };

  return (
    <div>
      <Steps current={current} style={{ marginBottom: 24 }}>
        {steps.map(item => (<Step key={item.title} title={item.title} />))}
      </Steps>
      {steps[current].content}
      <div style={{ marginTop: 24 }}>
        {current < steps.length - 1 && (
          <Button type="primary" onClick={next}>下一步</Button>
        )}
        {current === steps.length - 1 && (
          <Button type="primary" onClick={next}>完成</Button>
        )}
        {current > 0 && (
          <Button style={{ marginLeft: 8 }} onClick={prev}>上一步</Button>
        )}
      </div>
      <Modal open={modalVisible} onCancel={handleFinish} onOk={handleFinish} title="配置结果">
        <div>配置完成！</div>
      </Modal>
    </div>
  );
};

export default GuideStepsForm; 