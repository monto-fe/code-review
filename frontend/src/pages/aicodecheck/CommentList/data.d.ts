import type { Dayjs } from 'dayjs';

export interface TableQueryParam {
  id?: number;
  current?: number;
  page_size?: number;
  start_date?: number;
  end_date?: number;
  passed?: number;
  human_ratings?: number[];
  project_ids?: number[] | string;
  project_namespaces?: string[];
  date?: Dayjs
}

export interface TableListItem {
  id: number;
  project_id: string;
  merge_id: string;
  ai_model: string;
  rule: string;
  rule_id: string;
  result: string;
  passed: boolean;
  description: string;
  human_rating?: number;
  remark?: string;
  create_time?: number;
  update_time?: number;
}
