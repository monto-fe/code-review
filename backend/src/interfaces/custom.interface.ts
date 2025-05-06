export interface CustomRule {
    id: number;
    project_name: string;
    project_id: string;
    rule: string;
    status: number;
    operator: string;
    create_time: number;
    update_time: number;
}

export interface CustomRuleCreate {
    project_name: string;
    project_id: string;
    rule: string;
    status: number;
    operator: string;
    create_time?: number;
    update_time?: number;
}