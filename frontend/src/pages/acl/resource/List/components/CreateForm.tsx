import { useContext } from 'react';
import { Form } from 'antd';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import FormModal from '@/pages/component/Form/FormModal';
import { TableListItem } from '../../data.d';

import { ICreateFormProps, IFormItem } from '@/@types/form';
import { FormType } from '@/@types/enum';

function CreateForm(props: ICreateFormProps) {
  const { visible, setVisible, initialValues, onSubmit, onSubmitLoading, onCancel } = props;

  const [form] = Form.useForm();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const addFormItems: IFormItem[] = [
    {
      label: t('page.resource.name'),
      name: 'name',
      required: true,
      type: FormType.Input,
    },
    {
      label: t('page.resource.key'),
      name: 'resource',
      required: true,
      type: FormType.Input,
    },
    {
      label: t('page.resource.describe'),
      name: 'describe',
      type: FormType.Input,
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
        title={t('page.resource.add')}
        formLayout={{
          labelCol: { span: 8 },
          wrapperCol: { span: 14 }
        }}
        ItemOptions={addFormItems}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
}

export default CreateForm;
