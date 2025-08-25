import type { TableColumnsType } from 'antd';
import { FormType } from "@/@types/enum";

export interface PaginationConfig {
  total: number;
  current: number;
  page_size: number;
  showSizeChanger: boolean;
  showQuickJumper: boolean;
}

export interface ITableFilterItem {
  label: string;
  name: string;
  type: FormType;
  option?: unknown;
  options?: any;
  span?: number;
  required?: boolean;
}

export interface ITable<T> {
  queryList: Function;
  columns: ColumnsType<T>;
  title?: React.ReactElement | string;
  rowKey?: string;
  expandable?: TableColumnsType.ExpandableConfig;
  useTools?: boolean;
  fuzzySearchKey?: string;
  fuzzySearchPlaceholder?: string;
  filterFormItems?: ITableFilterItem[];
  scroll?: { x?: number; y?: number };
  reload?: Function;
  rightToolsSlot?: React.ReactNode;
}

export interface ITableFilter {
  items: ITableFilterItem[];
  size: SizeType;
  handleSearch?: Function;
}
