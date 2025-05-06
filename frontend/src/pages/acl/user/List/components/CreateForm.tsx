import { useContext } from 'react';
import { FormInstance } from 'antd/lib/form';
import { Form } from 'antd';

import { TableListItem as RoleTableListItem } from '@/pages/acl/role/data.d';
import FormModal from '@/pages/component/Form/FormModal';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { TableListItem, TableQueryParam } from '../data.d';
import { IFormItem } from '@/@types/form';
import { FormType } from '@/@types/enum';

interface ICreateFormProps {
  visible: boolean;
  setVisible: Function;
  initialValues?: Partial<TableQueryParam>;
  roleList: RoleTableListItem[];
  onSubmitLoading: boolean;
  onSubmit: (values: TableListItem, form: FormInstance) => void;
  onCancel?: () => void;
}

function CreateForm(props: ICreateFormProps) {
  const { visible, setVisible, roleList, initialValues, onSubmit, onSubmitLoading, onCancel } = props;

  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [form] = Form.useForm();

  const addFormItems: IFormItem[] = [
    {
      label: t('page.user.enname'),
      name: 'user',
      required: true,
      type: FormType.Input,
    },
    {
      label: t('page.user.cnname'),
      name: 'name',
      required: true,
      type: FormType.Input,
    },
    {
      label: t('page.user.password'),
      name: 'password',
      option: {
        placeholder: 'default: 12345678',
      },
      type: FormType.Input,
    },
    {
      label: t('page.user.job'),
      name: 'job',
      type: FormType.Input,
    },
    {
      label: t('page.user.email'),
      name: 'email',
      type: FormType.Input,
    },
    {
      label: t('page.user.phone'),
      name: 'phone_number',
      type: FormType.Input,
    },
    {
      label: t('page.user.role'),
      name: 'role_ids',
      type: FormType.SelectMultiple,
      options: (roleList || []).map((role: RoleTableListItem) => ({
        label: `${role.name} (${role.role})`,
        value: role.id,
      })),
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
        title={initialValues?.id ? t('app.global.edit') : t('page.user.add')}
        ItemOptions={addFormItems}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
}

export default CreateForm;
