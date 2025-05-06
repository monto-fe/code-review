import { FormType } from "@/@types/enum";

export interface PaginationConfig {
  total: number;
  current: number;
  pageSize: number;
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
  useTools?: boolean;
  fuzzySearchKey?: string;
  fuzzySearchPlaceholder?: string;
  filterFormItems?: ITableFilterItem[];
  scroll?: { x?: number; y?: number };
  reload?: Function;
}

export interface ITableFilter {
  items: ITableFilterItem[];
  size: SizeType;
  handleSearch?: Function;
}
