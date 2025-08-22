import { memo, useContext, useEffect, useRef, useState } from 'react';
import { Typography, message, Tag } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useSearchParams } from 'react-router-dom';
import 'github-markdown-css';

import { FormType } from '@/@types/enum';
import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
// import { FormType } from '@/@types/enum';
import { useI18n } from '@/store/i18n';
import ExcelExport from '@/components/ExcelExport';
import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';
import type { ExcelExportConfig } from '@/components/ExcelExport';

import { queryList, updateRating, queryProjectNamespaceList } from './service';
import Rate from './rate';
import Editable from './editable';
import { TableListItem } from './data.d';

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  
  const [searchParams] = useSearchParams();
  const id = searchParams.get('id');
  const [tableData, setTableData] = useState<TableListItem[]>([]);
  const [projectNamespaceOptions, setProjectNamespaceOptions] = useState<string[]>([]);

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
    console.log('params', params);
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

  const formItems = [
    {
      label: '创建时间',
      name: 'date',
      type: FormType.DateRange,
      span: 8,
    },
    {
      label: '状态',
      name: 'passed',
      type: FormType.Select,
      options: [
        {
          label: '成功',
          value: 1,
        },
        {
          label: '失败',
          value: -1,
        },
      ],
      span: 8,
    },
    {
      label: '评分',
      name: 'human_ratings',
      type: FormType.Select,
      options: [
        {
          label: '1星',
          value: 1,
        },
        {
          label: '2星',
          value: 2,
        },
        {
          label: '3星',
          value: 3,
        },
        {
          label: '4星',
          value: 4,
        },
        {
          label: '5星',
          value: 5,
        },
      ],
      span: 8,
    },
    {
      label: '项目命名空间',
      name: 'project_namespaces',
      type: FormType.Select,
      options: projectNamespaceOptions,
      span: 8,
    },
  ];

  useEffect(() => {
    queryProjectNamespaceList().then((res) => {
      setProjectNamespaceOptions(res.data?.data || []);
    });
  }, []);

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        columns={columns}
        filterFormItems={formItems}
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
