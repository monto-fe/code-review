import type { ColumnsType } from 'antd/lib/table';

export interface ExcelExportConfig {
  filename?: string;
  sheetName?: string;
  includeHeaders?: boolean;
  autoWidth?: boolean;
  selectedColumns?: string[];
}

export interface ExcelExportProps<T = any> {
  query: () => Promise<any>;
  columns: ColumnsType<T>;
  config?: ExcelExportConfig;
  buttonText?: string;
  buttonType?: 'default' | 'primary' | 'dashed' | 'link' | 'text';
  buttonSize?: 'large' | 'middle' | 'small';
  showSettings?: boolean;
  disabled?: boolean;
  loading?: boolean;
  onExport?: (config: ExcelExportConfig) => void;
  onBeforeExport?: (data: T[], config: ExcelExportConfig) => T[] | Promise<T[]>;
}

export interface ExcelDataProcessor<T = any> {
  processData: (data: T[], columns: ColumnsType<T>) => any[];
  getColumnValue: (item: T, column: any) => any;
} 