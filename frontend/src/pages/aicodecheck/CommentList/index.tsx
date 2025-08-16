import { memo, useContext, useRef, useState } from 'react';
import { Typography, message, Tag } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useSearchParams } from 'react-router-dom';
import 'github-markdown-css';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
// import { FormType } from '@/@types/enum';
import { useI18n } from '@/store/i18n';
import ExcelExport from '@/components/ExcelExport';
import type { ExcelExportConfig } from '@/components/ExcelExport';
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
  const [tableData, setTableData] = useState<TableListItem[]>([]);

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
      title: '项目',
      dataIndex: 'project_namespace',
      key: 'project_namespace',
      width: 100,
    },
    {
      title: t('page.aicodecheck.comment.result'),
      dataIndex: 'passed',
      key: 'passed',
      width: 80,
      render: (text: number) => text == 1 ? <Tag color='success' >成功</Tag> : <Tag color='error'>失败</Tag>,
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
      title: '评论信息',
      dataIndex: 'result',
      key: 'result',
      render: (text: string) => (
        <Typography className='w-650 markdown-body'>
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{text}</ReactMarkdown>
        </Typography>
      )
    },
    // {
    //   title: t('page.aicodecheck.comment.human_rating'),
    //   dataIndex: 'human_rating',
    //   key: 'human_rating',
    //   width: 200,
    //   render: (_, record:any) => <Rate id={record.id} initialValue={record.human_rating} />
    // },
    {
      title: t('page.aicodecheck.comment.improve_suggestion'),
      dataIndex: 'remark',
      key: 'remark',
      width: 200,
      render: (_, record:any) => {
        return <>
        <div style={{marginBottom: 10}}><Rate id={record.id} initialValue={record.human_rating} /></div>
        <Editable value={record.remark} onChange={(val) => {updateRemark(record, val)}} />
        </>
        
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
    const result = await queryList({
      ...params,
      id: id ? parseInt(id) : undefined
    });
    
    // 更新表格数据用于导出
    if (result.data?.data) {
      setTableData(result.data.data);
    }
    
    return result;
  };

  // 自定义数据处理器，处理特殊列的数据
  const handleBeforeExport = (data: TableListItem[], config: ExcelExportConfig): any[] => {
    return data.map(item => ({
      ...item,
      // 处理评论信息，移除Markdown格式
      result: item.result ? item.result.replace(/[#*`]/g, '').replace(/\n/g, ' ') : '',
      // 处理状态显示
      passed: item.passed ? '成功' : '失败'
    }));
  };

  // const formItems = [
  //   {
  //     label: t('page.resource.name'),
  //     name: 'project_namespace',
  //     type: FormType.Input,
  //     span: 8,
  //   },
  //   {
  //     label: t('page.resource.key'),
  //     name: 'resource',
  //     type: FormType.Input,
  //     span: 8,
  //   },
  // ];

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        columns={columns}
        // filterFormItems={formItems}
        queryList={handleQueryList}
        useTools
        scroll={{ x: 1200 }}
        rightToolsSlot={
          <ExcelExport
            data={tableData}
            columns={columns}
            config={{
              filename: 'comment_list',
              sheetName: '评论列表',
              selectedColumns: ['id', 'project_namespace', 'passed', 'merge_url', 'result', 'human_rating', 'remark', 'create_time']
            }}
            buttonText="导出Excel"
            buttonType="primary"
            buttonSize="small"
            showSettings={false}
            onBeforeExport={handleBeforeExport}
            onExport={(config: ExcelExportConfig) => {
              console.log('导出配置:', config);
            }}
          />
        }
      />
    </div>
  );
}

export default memo(observer(App));
