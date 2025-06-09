import { memo, useContext, useRef } from 'react';
import { Typography, message } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useSearchParams } from 'react-router-dom';
import 'github-markdown-css';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import Rate from './rate';
import { queryList, updateRating } from './service';
import { TableListItem } from './data.d';
import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';
import Editable from './editable';

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  const [searchParams] = useSearchParams();
  const id = searchParams.get('id');

  const updateRemark = (record: any, val: string) => {
    updateRating(record.id, record.human_rating, val).then(() => {
      message.success(t('app.global.tip.update.success'));
    });
  }

  const columns: ColumnsType<TableListItem> = [
    {
      title: 'id',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: 'MergeUrl',
      dataIndex: 'merge_url',
      key: 'merge_url',
      render: (text: string) => (
        <a href={text} target='_blank' rel="noreferrer">查看</a>
      )
    },
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
      title: t('page.aicodecheck.comment.human_rating'),
      dataIndex: 'human_rating',
      key: 'human_rating',
      width: 200,
      render: (_, record:any) => <Rate id={record.id} initialValue={record.human_rating} />
    },
    {
      title: t('page.aicodecheck.comment.improve_suggestion'),
      dataIndex: 'remark',
      key: 'remark',
      width: 200,
      render: (_, record:any) => {
        return <Editable value={record.remark} onChange={(val) => {updateRemark(record, val)}} />
      }
    },
    {
      title: t('page.aicodecheck.comment.createtime'),
      dataIndex: 'create_time',
      key: 'create_time',
      width: 160,
      render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time),
    }
    // {
    //   title: t('page.aicodecheck.comment.status'),
    //   fixed: 'right',
    //   dataIndex: 'passed',
    //   key: 'passed',
    //   render: (text: boolean) => text ? <Tag color='success' >{t('page.aicodecheck.comment.status.pass')}</Tag> : <Tag color='error'>{t('page.aicodecheck.comment.status.fail')}</Tag>,
    // },
  ];

  const handleQueryList = async (params?: any) => {
    return queryList({
      ...params,
      id: id ? parseInt(id) : undefined
    });
  };

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        columns={columns}
        queryList={handleQueryList}
        useTools
        scroll={{ x: 1200 }}
      />
    </div>
  );
}

export default memo(observer(App));
