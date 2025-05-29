import React, { useContext } from 'react';
import { Form } from 'antd';

import FormModal from '@/pages/component/Form/FormModal';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { ICreateFormProps, IFormItem } from '@/@types/form';
import { FormType } from '@/@types/enum';
import { TableListItem } from '../data';

const CreateForm: React.FC<ICreateFormProps> = (props) => {
  const { visible, setVisible, initialValues, onSubmit, onSubmitLoading, onCancel } = props;

  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [form] = Form.useForm();

  const addFormItems: IFormItem[] = [
    {
      label: t('page.aicodecheck.gitlab.api'),
      name: 'api',
      type: FormType.Input,
      tooltip: '假如gitlab部署域名为https://demo.com，则填写https://demo.com/api'
    },
    {
      label: t('page.aicodecheck.gitlab.token'),
      name: 'token',
      type: FormType.Input,
      tooltip: '获取地址：https://xxx.com/-/profile/personal_access_tokens, 需要所有权限'
    },
    {
      label: t('page.aicodecheck.gitlab.webhook_name'),
      name: 'webhook_name',
      type: FormType.Input,
      tooltip: '企业微信机器人名称或webhook名称'
    },
    {
      label: t('page.aicodecheck.gitlab.webhook_url'),
      name: 'webhook_url',
      type: FormType.Input,
      tooltip: '企业微信机器人地址或webhook地址'
    },
    {
      label: t('page.aicodecheck.gitlab.source_branch'),
      name: 'source_branch',
      type: FormType.Input,
      tooltip: '触发源分支，不填写则匹配全部'
    },
    {
      label: t('page.aicodecheck.gitlab.target_branch'),
      name: 'target_branch',
      type: FormType.Input,
      tooltip: '触发目标分支，不填写则匹配全部'
    },
    {
      label: t('page.aicodecheck.gitlab.expired'),
      name: 'expired',
      type: FormType.Date
    },
  ];

  const onFinish = async (values: TableListItem) => {
    onSubmit({ ...values, expired: values.expired?.unix() }, form);
  };

  return (
    <>
      <FormModal
        visible={visible}
        setVisible={setVisible}
        confirmLoading={onSubmitLoading}
        initialValues={initialValues}
        title={t('page.aicodecheck.gitlab.add')}
        ItemOptions={addFormItems}
        formLayout={{
          labelCol: { span: 8 },
          wrapperCol: { span: 16 },
        }}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
};

export default CreateForm;
