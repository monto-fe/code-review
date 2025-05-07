import { FormType } from "./enum";

interface IFormItem {
  option?: any;
  label: React.ReactNode | string;
  disabled?: boolean;
  key?: string;
  name: string;
  required?: boolean;
  placeholder?: string;
  options?: unknown[];
  content?: unknown;
  hidden?: boolean;
  span?: number;
  validators?: any;
  type: FormType;
  tooltip?: React.ReactNode | string;
}

interface IForm {
  items: IFormItem[];
  size?: SizeType;
  handleSearch?: Function;
}

interface IFormModal {
  formInstance?: any;
  visible: boolean;
  setVisible: Function;
  width?: number;
  title: React.ReactElement | string | undefined;
  description?: React.ReactElement | string | undefined;
  footer?: React.ReactElement | string | undefined;
  confirmLoading?: boolean;
  ItemOptions: IFormItem[];
  initialValues?: any;
  onFieldsChange?: Function;
  formLayout?: {
    labelCol: { span: number };
    wrapperCol: { span: number };
  };
  onFinish?: Function;
  onCancel?: Function;
}

interface FieldData {
  name: string | number | (string | number)[];
  value?: any;
  touched?: boolean;
  validating?: boolean;
  errors?: string[];
}

interface ICreateFormProps {
  visible: boolean;
  setVisible: Function;
  initialValues?: Partial<TableListItem>;
  onSubmitLoading: boolean;
  onSubmit: (values: TableListItem, form: FormInstance) => void;
  onCancel?: () => void;
}