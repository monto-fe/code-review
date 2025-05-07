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
      required: true,
      type: FormType.Input,
      disabled: true
    },
    {
      label: t('page.aicodecheck.gitlab.token'),
      name: 'token',
      required: true,
      type: FormType.Input,
    },
    {
      label: t('page.aicodecheck.gitlab.expired'),
      name: 'expired',
      required: true,
      type: FormType.Date
    },
  ];

  const onFinish = async (values: TableListItem) => {
    onSubmit({ ...values }, form);
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
          labelCol: { span: 7 },
          wrapperCol: { span: 15 },
        }}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
};

export default CreateForm;
