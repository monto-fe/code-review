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
      label: t('page.aicodecheck.aimodel.name'),
      name: 'name',
      required: true,
      type: FormType.Input,
      disabled: true,
    },
    {
      label: t('page.aicodecheck.aimodel.api_url'),
      name: 'api_url',
      required: true,
      type: FormType.Input,
      disabled: true,
    },
    {
      label: t('page.aicodecheck.aimodel.api_key'),
      name: 'api_key',
      required: true,
      type: FormType.Input,
      tooltip: '获取地址: https://console.ucloud.cn/modelverse/experience/api-keys'
    },
    {
      label: t('page.aicodecheck.aimodel.model'),
      name: 'model',
      required: true,
      type: FormType.Input,
      disabled: true
    }
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
        title={t('page.aicodecheck.aimodel.add')}
        ItemOptions={addFormItems}
        formLayout={{
          labelCol: { span: 4 },
          wrapperCol: { span: 18 },
        }}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
};

export default CreateForm;
