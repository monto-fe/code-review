import { memo, useContext, useRef } from 'react';
import { Tag, Typography } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import 'github-markdown-css';


import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { queryList } from './service';
import { TableListItem } from './data.d';
import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';

const { Paragraph } = Typography;

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const columns: ColumnsType<TableListItem> = [
    {
      title: 'id',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    // {
    //   title: t('page.aicodecheck.rule.project_id'),
    //   dataIndex: 'project_id',
    //   key: 'project_id',
    // },
    {
      title: 'MergeUrl',
      dataIndex: 'merge_url',
      key: 'merge_url',
      render: (text: string) => (
        <a href={text} target='_blank' rel="noreferrer">查看</a>
      )
    },
    // {
    //   title: t('page.aicodecheck.rule'),
    //   dataIndex: 'rule',
    //   key: 'rule'
    // },
    // {
    //   title: t('page.aicodecheck.rule.id'),
    //   dataIndex: 'rule_id',
    //   key: 'rule_id'
    // },
    {
      title: t('page.aicodecheck.comment.result'),
      dataIndex: 'result',
      key: 'result',
      render: (text: string) => (
        <Typography className='w-650 markdown-body'>
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{text}</ReactMarkdown>
        </Typography>
      )
    },
    {
      title: t('page.aicodecheck.comment.aimodel'),
      dataIndex: 'ai_model',
      key: 'ai_model',
    },
    {
      title: t('page.aicodecheck.comment.createtime'),
      dataIndex: 'create_time',
      key: 'create_time',
      width: 160,
      render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time),
    },
    {
      title: t('page.aicodecheck.comment.status'),
      fixed: 'right',
      dataIndex: 'passed',
      key: 'passed',
      render: (text: boolean) => text ? <Tag color='success' >{t('page.aicodecheck.comment.status.pass')}</Tag> : <Tag color='error'>{t('page.aicodecheck.comment.status.fail')}</Tag>,
    },
  ];

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        columns={columns}
        queryList={queryList}
        useTools
      />
    </div>
  );
}

export default memo(observer(App));
