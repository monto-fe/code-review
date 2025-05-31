import React, { useState } from 'react';
import { Button, message } from 'antd';
import FormModal from '@/pages/component/Form/FormModal';
import { IFormItem } from '@/@types/form';
import { FormType } from '@/@types/enum';

const FeedbackForm = () => {
  const [visible, setVisible] = useState(false);
  const [submitLoading, setSubmitLoading] = useState(false);

  const formItems: IFormItem[] = [
    {
      label: '反馈类型',
      name: 'type',
      required: true,
      type: FormType.Select,
      options: [
        { label: '功能建议', value: 'feature' },
        { label: 'Bug反馈', value: 'bug' },
        { label: '其他', value: 'other' },
      ],
    },
    {
      label: '反馈内容',
      name: 'content',
      required: true,
      type: FormType.TextArea,
    },
    {
      label: '联系方式',
      name: 'contact',
      type: FormType.Input,
    },
  ];

  const handleSubmit = (values: any) => {
    setSubmitLoading(true);
    // 模拟提交
    setTimeout(() => {
      message.success('反馈提交成功！');
      setVisible(false);
      setSubmitLoading(false);
    }, 1000);
  };

  return (
    <div>
      <Button type='primary' onClick={() => setVisible(true)}>提交反馈</Button>
      <FormModal
        visible={visible}
        setVisible={setVisible}
        title='产品建议反馈'
        ItemOptions={formItems}
        onFinish={handleSubmit}
        confirmLoading={submitLoading}
      />
    </div>
  );
};

export default FeedbackForm; 