import { useContext, useEffect } from 'react';
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

  // 处理编辑时的初始值，排除密码字段
  const getFormInitialValues = () => {
    if (initialValues?.id) {
      // 编辑时排除密码字段
      const { password, ...valuesWithoutPassword } = initialValues;
      return valuesWithoutPassword;
    }
    return initialValues;
  };

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
      required: !initialValues?.id, // 编辑时密码不是必填
      option: {
        placeholder: initialValues?.id ? '不修改请留空' : 'default: 12345678',
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
    // 编辑时如果密码为空，则从提交数据中移除密码字段
    if (initialValues?.id && !values.password) {
      const { password, ...valuesWithoutPassword } = values;
      onSubmit(valuesWithoutPassword, form);
    } else {
      onSubmit({ ...values }, form);
    }
  };

  return (
    <>
      <FormModal
        visible={visible}
        setVisible={setVisible}
        confirmLoading={onSubmitLoading}
        initialValues={getFormInitialValues()}
        title={initialValues?.id ? t('app.global.edit') : t('page.user.add')}
        ItemOptions={addFormItems}
        onFinish={onFinish}
        onCancel={onCancel}
      />
    </>
  );
}

export default CreateForm;
