import { PageModel } from './common.interface';
// namespace
export interface Namespace {
    id: number;
    namespace: string;
    parent: string;
    name: string;
    describe: string;
    operator: string;
    create_time: number;
    update_time: number;
}

export interface NamespaceReqList extends PageModel {
    namespace?: string;
}

export interface NamespaceReq {
    namespace: string;
    parent: string;
    name: string;
    describe: string;
    operator: string;
}

export interface TimeModel {
    create_time: number;
    update_time: number;
}

export interface NamespaceParams extends TimeModel {
    namespace: string;
    parent: string;
    name: string;
    describe: string;
    operator: string;
}
