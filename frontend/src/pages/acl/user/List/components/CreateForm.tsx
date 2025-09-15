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

  // 验证函数
  const validatePassword = (rule: any, value: string) => {
    if (!value) {
      return Promise.resolve();
    }
    
    if (value.length < 8) {
      return Promise.reject(new Error('密码长度不能少于8位'));
    }
    
    if (value.length > 50) {
      return Promise.reject(new Error('密码长度不能超过50位'));
    }
    
    // 检查是否包含字母
    const hasLetter = /[a-zA-Z]/.test(value);
    if (!hasLetter) {
      return Promise.reject(new Error('密码必须包含字母'));
    }
    
    // 检查字符是否合法
    const validChars = /^[A-Za-z0-9@$!%*?&._]+$/;
    if (!validChars.test(value)) {
      return Promise.reject(new Error('密码只能包含字母、数字和特殊字符(@$!%*?&._)'));
    }
    
    return Promise.resolve();
  };

  const validateEmail = (rule: any, value: string) => {
    if (!value) {
      return Promise.resolve();
    }
    
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(value)) {
      return Promise.reject(new Error('邮箱格式不正确'));
    }
    
    if (value.length > 100) {
      return Promise.reject(new Error('邮箱长度不能超过100个字符'));
    }
    
    return Promise.resolve();
  };

  const validatePhone = (rule: any, value: string) => {
    if (!value) {
      return Promise.resolve();
    }
    
    // 检查是否只包含数字
    if (!/^\d+$/.test(value)) {
      return Promise.reject(new Error('手机号只能包含数字'));
    }
    
    // 检查长度是否为11位
    if (value.length !== 11) {
      return Promise.reject(new Error('手机号长度必须为11位'));
    }
    
    // 检查是否以1开头
    if (value[0] !== '1') {
      return Promise.reject(new Error('手机号必须以1开头'));
    }
    
    return Promise.resolve();
  };

  const validateUsername = (rule: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('用户名不能为空'));
    }
    
    // 检查是否以字母开头
    if (!/^[a-zA-Z]/.test(value)) {
      return Promise.reject(new Error('用户名必须以字母开头'));
    }
    
    // 检查长度
    if (value.length < 3 || value.length > 20) {
      return Promise.reject(new Error('用户名长度必须在3-20位之间'));
    }
    
    // 检查字符是否合法
    if (!/^[a-zA-Z0-9_]+$/.test(value)) {
      return Promise.reject(new Error('用户名只能包含字母、数字和下划线'));
    }
    
    return Promise.resolve();
  };

  const validateName = (rule: any, value: string) => {
    if (!value) {
      return Promise.reject(new Error('姓名不能为空'));
    }
    
    // 检查长度
    if (value.length < 2 || value.length > 50) {
      return Promise.reject(new Error('姓名长度必须在2-50个字符之间'));
    }
    
    // 检查是否包含特殊字符（除了中文、英文、数字、空格）
    if (/[^\u4e00-\u9fa5a-zA-Z0-9\s]/.test(value)) {
      return Promise.reject(new Error('姓名不能包含特殊字符'));
    }
    
    return Promise.resolve();
  };

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
      option: {
        placeholder: '请输入用户名（字母开头，3-20位）',
      },
      validators: [
        {
          validator: validateUsername,
        },
      ],
      type: FormType.Input,
    },
    {
      label: t('page.user.cnname'),
      name: 'name',
      required: true,
      option: {
        placeholder: '请输入姓名（2-50个字符）',
      },
      validators: [
        {
          validator: validateName,
        },
      ],
      type: FormType.Input,
    },
    {
      label: t('page.user.password'),
      name: 'password',
      required: !initialValues?.id, // 编辑时密码不是必填
      option: {
        placeholder: initialValues?.id ? '不修改请留空' : 'default: 12345678',
        type: 'password',
      },
      validators: [
        {
          validator: validatePassword,
        },
      ],
      type: FormType.Input,
    },
    {
      label: t('page.user.job'),
      name: 'job',
      option: {
        placeholder: '请输入职位',
      },
      validators: [
        {
          max: 100,
          message: '职位长度不能超过100个字符',
        },
      ],
      type: FormType.Input,
    },
    {
      label: t('page.user.email'),
      name: 'email',
      option: {
        placeholder: '请输入邮箱地址',
      },
      validators: [
        {
          validator: validateEmail,
        },
      ],
      type: FormType.Input,
    },
    {
      label: t('page.user.phone'),
      name: 'phone_number',
      option: {
        placeholder: '请输入11位手机号',
      },
      validators: [
        {
          validator: validatePhone,
        },
      ],
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
