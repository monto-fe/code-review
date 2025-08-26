import { memo, useContext, useEffect, useRef, useState } from 'react';
import { Typography, message, Tag, Flex, Space, Tooltip } from 'antd';
import type { GetProps, DatePicker, TableColumnsType } from 'antd';
import { observer } from 'mobx-react-lite';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useSearchParams } from 'react-router-dom';
import dayjs from 'dayjs';
import { NodeExpandOutlined, NodeCollapseOutlined } from '@ant-design/icons';

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
import { queryTokenList } from '@/pages/aicodecheck/GitlabToken/service';

import Rate from './rate';
import Editable from './editable';
import { TableListItem, TableQueryParam } from './data.d';

type RangePickerProps = GetProps<typeof DatePicker.RangePicker>;

const disabledDate: RangePickerProps['disabledDate'] = (current) => {
  // Can not select days before today and today
  return current && current > dayjs().add(1, 'd').endOf('d');
};

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [searchParams] = useSearchParams();
  const id = searchParams.get('id');

  const [tableData, setTableData] = useState<TableListItem[]>([]);
  const [filters, setFilters] = useState<TableQueryParam>();
  const [projectNamespaceOptions, setProjectNamespaceOptions] = useState<string[]>([]);
  const [tokenOptions, setTokenOptions] = useState<string[]>([]);
  const [expandedRowKeys, setExpandedRowKeys] = useState<number[]>([]);

  const updateRemark = (record: any, val: string) => {
    updateRating(record.id, record.human_rating, val).then(() => {
      message.success(t('app.global.tip.update.success'));
    });
  }

  const columns: TableColumnsType<TableListItem> = [
    {
      title: 'ID',
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
      width: 80,
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
    {
      title: t('page.aicodecheck.comment.createtime'),
      dataIndex: 'create_time',
      key: 'create_time',
      width: 160,
      render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time),
    },
    {
      title: t('page.aicodecheck.comment.improve_suggestion'),
      dataIndex: 'remark',
      key: 'remark',
      width: 120,
      fixed: 'right',
      render: (_, record: any) => {
        return <>
          <div style={{ marginBottom: 10 }}><Rate id={record.id} initialValue={record.human_rating} /></div>
          <Editable value={record.remark} onChange={(val) => { updateRemark(record, val) }} />
        </>
      }
    },
  ];

  const handleQueryList = async (params?: TableQueryParam) => {
    const data: any = {
      passed: params?.passed,
      page_size: params?.page_size,
      current: params?.current,
      human_ratings: params?.human_ratings,
      project_namespaces: params?.project_namespaces,
      project_ids: params?.project_ids && typeof params.project_ids === 'string' ? params.project_ids.split(',').map(Number) : void 0,
      id: id ? parseInt(id) : undefined
    }
    if (Array.isArray(params?.date)) {
      data.start_date = dayjs(params?.date[0]).startOf('d').unix();
      data.end_date = dayjs(params?.date[1]).endOf('d').unix();
    }

    setFilters(data);
    const result = await queryList(data);

    // 更新表格数据用于导出
    if (result.data?.data) {
      setTableData(result.data.data);
      setExpandedRowKeys(result.data.data.map((item: TableListItem) => item.id));
    }

    return result;
  };

  const handleDownloadList = async () => {
    return await queryList({ ...filters, page_size: 99999, current: 1 }).then(res => res.data?.data || []);
  };

  // 自定义数据处理器，处理特殊列的数据
  const handleBeforeExport = (data: TableListItem[], config: ExcelExportConfig): any[] => {
    return data.map(item => ({
      ...item,
      // 处理评论信息，移除Markdown格式
      result: item.result ? item.result.replace(/[#*`]/g, '').replace(/\n/g, ' ') : '',
      // 处理状态显示
      passed: item.passed ? '成功' : '失败',
      create_time: renderDateFromTimestamp(item.create_time as number, timeFormatType.time),
    }));
  };

  const formItems = [
    {
      label: '创建时间',
      name: 'date',
      type: FormType.DateRange,
      option: {
        disabledDate
      },
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
      type: FormType.SelectMultiple,
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
      type: FormType.SelectMultiple,
      options: projectNamespaceOptions,
      span: 8,
    },
    {
      label: 'Token Projects',
      name: 'project_ids',
      type: FormType.Select,
      options: tokenOptions,
      span: 8,
    },
  ];

  useEffect(() => {
    queryProjectNamespaceList().then((res) => {
      setProjectNamespaceOptions((res.data?.data || []).map((item: any) => ({
        label: item.project_namespace,
        value: item.project_namespace,
      })));
    });

    queryTokenList().then((res) => {
      setTokenOptions((res.data?.data || []).map((item: any) => ({
        label: item.name,
        value: item.project_ids,
      })));
    });
  }, []);

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        expandable={{
          expandedRowRender: (record: TableListItem) => (
            <Flex align='center' justify='center' className='mx-24'>
              <Typography className='w-full markdown-body'>
                <ReactMarkdown remarkPlugins={[remarkGfm]}>{record.result}</ReactMarkdown>
              </Typography>
            </Flex>
          ),
          rowExpandable: (record: TableListItem) => true,
          expandedRowKeys,
          onExpandedRowsChange: (keys: number[]) => setExpandedRowKeys(keys)
        }}
        columns={columns.filter(column => column.key !== 'result')}
        filterFormItems={formItems}
        queryList={handleQueryList}
        useTools
        scroll={{ x: 1200 }}
        rightToolsSlot={
          <Space size='middle'>
            <Tooltip title="展开全部">
              <NodeExpandOutlined onClick={() => setExpandedRowKeys(tableData.map((item: TableListItem) => item.id))} />
            </Tooltip>
            <Tooltip title="折叠全部">
              <NodeCollapseOutlined onClick={() => setExpandedRowKeys([])} />
            </Tooltip>
            <ExcelExport
              query={handleDownloadList}
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
          </Space>
        }
      />
    </div>
  );
}

export default memo(observer(App));
