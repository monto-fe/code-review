export interface AIDetectionStatsItem {
  level: number;
  value: number;
}

export function getAIDetectionStats(): Promise<{ data: AIDetectionStatsItem[] }>; 