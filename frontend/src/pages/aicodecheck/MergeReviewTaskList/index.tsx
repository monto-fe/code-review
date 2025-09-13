import { useContext, useRef } from 'react';
import { Button, Space, Tag, Tooltip } from 'antd';
import type { TableColumnsType } from 'antd';
import { observer } from 'mobx-react-lite';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { queryList } from './service';
import { TableListItem } from './data';
import styles from './index.module.less';

function MergeReviewTaskList() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  const navigate = useNavigate();

  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  // 查看详情
  const handleViewDetail = (id: number) => {
    navigate(`/aicodecheck/merge-review-detail?id=${id}`);
  };

  // 新建任务
  const handleCreate = () => {
    navigate('/aicodecheck/merge-review-create');
  };

  // 表格列定义
  const columns: TableColumnsType<TableListItem> = [
    {
      title: '任务ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
      ellipsis: true,
    },
    {
      title: '项目名称',
      dataIndex: 'project_name',
      key: 'project_name',
      width: 150,
      ellipsis: true,
    },
    {
      title: 'Merge ID',
      dataIndex: 'merge_id',
      key: 'merge_id',
      width: 100,
      ellipsis: true,
    },
    {
      title: 'Merge标题',
      dataIndex: 'merge_title',
      key: 'merge_title',
      width: 250,
      ellipsis: true,
      render: (text: string) => (
        <Tooltip title={text}>
          <span>{text}</span>
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status_text',
      key: 'status_text',
      width: 100,
      render: (statusText: string, record: TableListItem) => {
        const statusConfig: Record<number, { color: string; text: string }> = {
          0: { color: 'orange', text: '待处理' },
          1: { color: 'blue', text: '处理中' },
          2: { color: 'green', text: '完成' },
          3: { color: 'red', text: '失败' },
        };
        const config = statusConfig[record.status] || { color: 'default', text: statusText };
        return <Tag color={config.color} className={styles.statusTag}>{config.text}</Tag>;
      },
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      key: 'create_time',
      width: 160,
      render: (time: number) => new Date(time * 1000).toLocaleString('zh-CN'),
    },
    {
      title: '更新时间',
      dataIndex: 'update_time',
      key: 'update_time',
      width: 160,
      render: (time: number) => new Date(time * 1000).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      fixed: 'right',
      render: (_, record) => (
        <Tooltip title="查看详情">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record.id)}
            className={styles.actionButton}
          />
        </Tooltip>
      ),
    },
  ];



  const rightToolsSlot = (
    <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
      新建任务
    </Button>
  );

  return (
    <div className={styles.mergeReviewTaskList}>
      <div className={styles.tableContainer}>
        <CommonTable
          ref={tableRef}
          queryList={queryList}
          columns={columns}
          title="Merge Review 任务列表"
          rowKey="id"
          rightToolsSlot={rightToolsSlot}
          scroll={{ x: 1100 }}
        />
      </div>
    </div>
  );
}

export default observer(MergeReviewTaskList); 