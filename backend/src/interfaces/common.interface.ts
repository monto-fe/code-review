export interface CommonRule {
  id: number;
  name: string;
  language: 'Python' | 'Java' | 'JavaScript' | 'Golang' | 'Ruby' | 'C++' | 'other';
  rule: string;
  description?: string | null;  // 可选字段
  create_time: number;
  update_time: number;
}

export interface PageModel {
  limit: number;
  offset: number;
}

export interface CommonRuleCreate {
  name: string;
  language: 'Python' | 'Java' | 'JavaScript' | 'Golang' | 'Ruby' | 'C++' | 'other';
  rule: string;
  description?: string | null;  // 可选字段
  create_time?: number;
  update_time?: number;
}