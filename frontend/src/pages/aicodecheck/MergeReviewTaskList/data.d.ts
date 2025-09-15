export interface MergeReviewTask {
  id: number;
  project_id: number;
  project_name: string;
  merge_id: number;
  merge_title: string;
  status: number;
  status_text: string;
  create_time: number;
  update_time: number;
}

export interface BotResult {
  botId: string;
  botName: string;
  category: string;
  output: string;
  score?: number;
}

export interface TableQueryParam {
  current: number;
  page_size: number;
}

export interface TableListItem extends MergeReviewTask {} 