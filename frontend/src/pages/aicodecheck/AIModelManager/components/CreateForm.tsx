import React, { useContext, useState } from 'react';
import { Form, Modal, Input, Select } from 'antd';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { ICreateFormProps } from '@/@types/form';

const { Option } = Select;

const CreateForm: React.FC<ICreateFormProps> = (props) => {
  const { visible, setVisible, initialValues, onSubmit, onSubmitLoading, onCancel, modelOptions=[], typeOptions=[] } = props;
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [form] = Form.useForm();
  const [modelSelectOptions, setModelSelectOptions] = useState<any[]>([]);

  // type change handler
  const handleTypeChange = (value: string) => {
    const filteredModels = modelOptions
      .filter((item: any) => item.type === value)
      .map((item: any) => ({
        label: item.model,
        value: item.model
      }));
    setModelSelectOptions(filteredModels);
    form.setFieldsValue({ model: undefined });
  };

  // submit handler
  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      onSubmit({ ...values }, form);
    } catch (e) {}
  };

  // cancel handler
  const handleCancel = () => {
    form.resetFields();
    setModelSelectOptions([]);
    onCancel && onCancel();
    setVisible(false);
  };

  // initialValues 联动
  React.useEffect(() => {
    if (visible && initialValues && initialValues.type) {
      handleTypeChange(initialValues.type);
      form.setFieldsValue(initialValues);
    } else if (visible) {
      form.resetFields();
      setModelSelectOptions([]);
    }
    // eslint-disable-next-line
  }, [visible]);

  console.log("selected", modelSelectOptions);

  return (
    <Modal
      open={visible}
      title={t('page.aicodecheck.aimodel.add')}
      onOk={handleOk}
      confirmLoading={onSubmitLoading}
      onCancel={handleCancel}
      destroyOnClose
      maskClosable={false}
      okText={t('app.global.confirm')}
      cancelText={t('app.global.close')}
    >
      <Form
        form={form}
        layout="horizontal"
        labelCol={{ span: 5 }}
        wrapperCol={{ span: 18 }}
        initialValues={initialValues}
      >
        <Form.Item
          label={t('page.aicodecheck.aimodel.resource')}
          name="type"
          rules={[{ required: true, message: t('app.form.required') }]}
        >
          <Select
            placeholder={t('page.aicodecheck.aimodel.resource')}
            onChange={handleTypeChange}
            allowClear
          >
            {typeOptions.map((item: any) => (
              <Option key={item.value} value={item.value}>
                {item.label}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label={t('page.aicodecheck.aimodel.model')}
          name="model"
          rules={[{ required: true, message: t('app.form.required') }]}
        >
          <Select
            placeholder={t('page.aicodecheck.aimodel.model')}
            allowClear
            disabled={modelSelectOptions.length === 0}
          >
            {modelSelectOptions.map((item: any) => (
              <Option key={item.value} value={item.value}>
                {item.label}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label={t('page.aicodecheck.aimodel.api_key')}
          name="api_key"
          tooltip={<div>获取方法: <a href="https://github.com/monto-fe/code-review/wiki/AI-Model%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E" target="_blank" rel="noopener noreferrer">https://github.com/monto-fe/code-review/wiki/AI-Model%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E</a></div>}
          rules={[{ required: true, message: t('app.form.required') }]}
        >
          <Input placeholder={t('page.aicodecheck.aimodel.api_key')} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateForm;
