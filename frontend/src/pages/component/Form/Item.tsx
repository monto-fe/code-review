import { IFormItem } from '@/@types/form';
import { Input, InputNumber, Select, DatePicker, Radio, Switch } from 'antd';

const { Option } = Select;
const { RangePicker } = DatePicker;

/**
 *  表单弹窗 
 * @param ItemOptions 表单配置
 * @description ItemOptions = {
    label: "标题",
    name: "字段名",
    required: "必填-boolean",
    content: "表单",
    hidden: "隐藏-boolean",
    title: "作为文本标题展示-boolean"
 * }
 * 选项的类型：Input-输入框，NumberInput-数字输入框，Select-下拉框，SelectMultiple-下拉框多选，
 * Date-日期，DateRange-日期范围，DateRangeButton-日期范围按钮
 * 
 * @example 
 * @param initialValues 表单整体的默认值，不再使用单个表单默认值
 * */

class FormItemComponent {
  static Input(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Input
        style={{ width: '100%' }}
        placeholder={field.placeholder}
        disabled={field.disabled}
        {...field.option}
      />
    )
  }

  static TextArea(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Input.TextArea style={{ width: '100%' }} placeholder={field.label} {...field.option} />
    )
  }

  static NumberInput(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <InputNumber style={{ width: '100%' }} placeholder={field.label} {...field.option} />
    )
  }

  static Select(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Select
        allowClear={true}
        {...field.option}
        style={{ width: '100%' }}
        placeholder={field.label}
        showSearch={true}
        disabled={field.disabled || false}
        filterOption={(input, option: any) => {
          try {
            return option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0;
          } catch {
            return option.props.children.indexOf(input) >= 0;
          }
        }}
      >
        {field.options
          ? field.options.map((item: any, index: number) => (
            <Option key={index} value={item.value}>
              {item.name || item.label}
            </Option>
          ))
          : null}
      </Select>
    )
  }

  static SelectMultiple(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Select
        {...field.option}
        style={{ width: '100%' }}
        mode='multiple'
        allowClear={true}
        placeholder={field.label}
        filterOption={(input, option: any) => {
          try {
            return option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0;
          } catch {
            return option.props.children.indexOf(input) >= 0;
          }
        }}
      >
        {field.options
          ? field.options.map((item: any, index: number) => (
            <Option key={index} value={item.value}>
              {item.name || item.label}
            </Option>
          ))
          : null}
      </Select>
    )
  }

  static Radio(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Radio.Group options={field.options} optionType='button' buttonStyle='solid' {...field.option} />
    )
  }

  static Switch(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <Switch {...field.option} />
    )
  }

  static Date(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <DatePicker {...field.option} />
    )
  }

  static DateRange(props: { field: IFormItem }) {
    const { field } = props;

    return (
      <RangePicker {...field.option} />
    )
  }

  static Content(props: { field: IFormItem }) {
    const { field } = props;

    return field.content;
  }
}

export default FormItemComponent;
