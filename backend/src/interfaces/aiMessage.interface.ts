export interface AImessage {
    id: number;
    project_id: number;
    merge_url: string;
    merge_id: string;
    ai_model: string;
    rule: 1 | 2; // 1: common 2: custom
    rule_id: number;
    result: string;
    passed?: boolean;  // 默认为 `0`，可以是 `false` 或 `true`
    checked_by?: string | null;  // 可选字段
    create_time: number;
}