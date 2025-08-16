import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';

export interface ColumnValueExtractor {
  extractValue: (value: any, record: any) => string | number | null;
}

// 默认值提取器
export const defaultExtractor: ColumnValueExtractor = {
  extractValue: (value: any) => {
    if (value === null || value === undefined) return '';
    if (typeof value === 'string' || typeof value === 'number') return value;
    return String(value);
  }
};

// 时间戳提取器
export const timestampExtractor: ColumnValueExtractor = {
  extractValue: (value: any) => {
    if (!value) return '';
    return renderDateFromTimestamp(value, timeFormatType.time);
  }
};

// 状态标签提取器
export const statusExtractor: ColumnValueExtractor = {
  extractValue: (value: any) => {
    if (value === 1) return '成功';
    if (value === 0) return '失败';
    return value ? '成功' : '失败';
  }
};

// 评分提取器
export const ratingExtractor: ColumnValueExtractor = {
  extractValue: (value: any) => {
    if (value === null || value === undefined) return '';
    return `${value}分`;
  }
};

// 链接提取器
export const linkExtractor: ColumnValueExtractor = {
  extractValue: (value: any) => {
    if (!value) return '';
    return value;
  }
};

// 创建自定义提取器
export const createCustomExtractor = (extractFn: (value: any, record: any) => string | number | null): ColumnValueExtractor => ({
  extractValue: extractFn
});

// 获取列的值
export const getColumnValue = (item: any, column: any, extractors: Record<string, ColumnValueExtractor> = {}) => {
  const { dataIndex, key } = column;
  
  if (!dataIndex) return '';
  
  const value = item[dataIndex];
  
  // 使用自定义提取器
  if (key && extractors[key]) {
    return extractors[key].extractValue(value, item);
  }
  
  // 根据数据类型使用默认提取器
  if (column.render) {
    // 对于有render函数的列，尝试获取原始值
    return defaultExtractor.extractValue(value, item);
  }
  
  return defaultExtractor.extractValue(value, item);
};

// 处理表格数据为Excel格式
export const processTableData = (
  data: any[], 
  columns: any[], 
  selectedColumns: string[],
  extractors: Record<string, ColumnValueExtractor> = {}
) => {
  const filteredColumns = columns.filter(col => 
    selectedColumns.includes(col.key as string)
  );
  
  return data.map(item => {
    const row: any = {};
    filteredColumns.forEach(col => {
      const columnKey = col.key as string;
      const columnTitle = col.title as string;
      row[columnTitle] = getColumnValue(item, col, extractors);
    });
    return row;
  });
}; 