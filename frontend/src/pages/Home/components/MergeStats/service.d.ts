export interface MergeStatsResponse {
  total: number;
}

export function getMergeStats(): Promise<{ data: MergeStatsResponse }>; 