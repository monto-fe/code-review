export enum HumanRating {
    EXCELLENT = 1,      // 发现Bug
    PARTIAL = 2,        // 部分有效：AI 建议部分正确
    NEEDS_IMPROVEMENT = 3, // 需要改进：AI 建议需要进一步优化
    INVALID = 4        // 无效
}

export interface AImessage {
    id: number;
    project_id: number;
    merge_url: string;
    merge_id: string;
    ai_model: string;
    rule: 1 | 2; // 1: common 2: custom
    rule_id: number;
    result: string;
    human_rating: HumanRating; // 使用枚举类型替代原来的数字
    remark: string; // 备注
    passed?: boolean;  // 默认为 `0`，可以是 `false` 或 `true`
    checked_by?: string | null;  // 可选字段
    create_time: number;
}