import { memo, useContext, useRef, useState } from 'react';
import { Button, message, Popconfirm, PopconfirmProps, Space, Switch } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';

import { queryList, removeData, updateData as updateDataService } from './service';
import { TableQueryParam, TableListItem } from './data';
// import CreateForm from './CreateToken';
import StatusTag from './components/StatusTag';

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  const navigate = useNavigate();
  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  // 删除
  const [deleteOpen, setDeleteOpen] = useState<number | undefined>();
  const [deleteLoading, setDeleteLoading] = useState<number | undefined>();

  const handleDelete = (id: number) => setDeleteOpen(id);

  const deleteConfirm = (id: number) => {
    setDeleteLoading(id);
    removeData(id).then(() => {
      message.success(t('app.global.tip.delete.success'));
      reload();
      setDeleteOpen(void 0);
    }).finally(() => {
      setDeleteLoading(void 0);
    });
  };

  const deleteCancel: PopconfirmProps['onCancel'] = () => {
    setDeleteOpen(void 0);
  };

  const handleDetail = (record: TableListItem) => {
    navigate(`/aicodecheck/GitlabTokenDetail?id=${record.id}`);
  };

  // 新增&编辑
  const [createFormVisible, setCreateFormVisible] = useState<boolean>(false);
  const [updateData, setUpdateData] = useState<Partial<TableQueryParam>>({});

  const handleCreate = () => {
    navigate('/aicodecheck/GitlabConfig');
    // setUpdateData({});
    // setCreateFormVisible(true);
  };
  

  const handleUpdate = (record: TableListItem) => {
    // 跳转到GitlabToken新建页面，带上id
    navigate(`/aicodecheck/GitlabConfig?id=${record.id}`);
  };

  const [loadingId, setLoadingId] = useState<number | null>(null);
  const [webhookLoadingId, setWebhookLoadingId] = useState<number | null>(null);
  const [statusLoadingId, setStatusLoadingId] = useState<number | null>(null);

  const handleRuleCheckStatusChange = async (id: number, checked: boolean) => {
    setLoadingId(id);
    try {
      const { ret_code } = await updateDataService({id, rule_check_status: checked ? 1 : 2});
      if (ret_code !== 0) {
        message.error('状态更新失败');
        return;
      }
      message.success('状态更新成功');
      reload();
    } catch (e) {
      message.error('状态更新失败');
    } finally {
      setLoadingId(null);
    }
  };

  const handleWebhookStatusChange = async (record: TableListItem, checked: boolean) => {
    // 增加逻辑，判断下webhook_url是否为空，如果为空，则提示用户需要先配置webhook_url
    const { id, webhook_url } = record;
    if (!webhook_url) {
      message.error('请先配置webhook_url');
      return;
    }
    setWebhookLoadingId(id);
    try {
      const { ret_code } = await updateDataService({id, webhook_status: checked ? 1 : 2});
      if (ret_code !== 0) {
        message.error('状态更新失败');
        return;
      }
      message.success('状态更新成功');
      reload();
    } catch (e) {
      message.error('状态更新失败');
    } finally {
      setWebhookLoadingId(null);
    }
  };

  const handleStatusChange = async (record: TableListItem, checked: boolean) => {
    setStatusLoadingId(record.id);
    try {
      const { ret_code } = await updateDataService({id: record.id, status: checked ? 1 : -1});
      if (ret_code !== 0) {
        message.error('状态更新失败');
        return;
      }
      message.success('状态更新成功');
      reload();
    } catch (e) {
      message.error('状态更新失败');
    } finally {
      setStatusLoadingId(null);
    }
  };

  const columns: ColumnsType<TableListItem> = [
    {
      title: 'id',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: t('page.aicodecheck.gitlab.api'),
      dataIndex: 'api',
      key: 'api',
    },
    {
      title: t('page.aicodecheck.gitlab.webhook_name'),
      dataIndex: 'webhook_name',
      key: 'webhook_name',
    },
    {
      title: t('page.aicodecheck.gitlab.webhook_url'),
      dataIndex: 'webhook_url',
      key: 'webhook_url'
    },
    {
      title: t('page.aicodecheck.gitlab.source_branch'),
      dataIndex: 'source_branch',
      key: 'source_branch',
    },
    {
      title: t('page.aicodecheck.gitlab.target_branch'),
      dataIndex: 'target_branch',
      key: 'target_branch',
    },
    {
      title: t('page.aicodecheck.gitlab.expired'),
      dataIndex: 'expired',
      key: 'expired',
      render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time)
    },
    {
      title: t('page.aicodecheck.gitlab.project_ids_synced'),
      dataIndex: 'project_ids_synced',
      key: 'project_ids_synced',
      render: (text: number) => {
        return <StatusTag status={1} />
      }
    },
    {
      title: t('page.aicodecheck.gitlab.status'),
      dataIndex: 'status',
      key: 'status',
      render: (text: number, record: TableListItem) => {
        return <Switch
        checked={text === 1}
        checkedChildren="开启"
        unCheckedChildren="关闭"
        onChange={(checked) => handleStatusChange(record, checked)}
        loading={statusLoadingId === record.id}
      />
      }
    },
    {
      title: t('page.aicodecheck.gitlab.webhook_status'),
      dataIndex: 'webhook_status',
      key: 'webhook_status',
      render: (text: number, record: TableListItem) => {
        return <Switch
        checked={text === 1}
        checkedChildren="开启"
        unCheckedChildren="关闭"
        onChange={(checked) => handleWebhookStatusChange(record, checked)}
        loading={webhookLoadingId === record.id}
      />
      }
    },
    {
      title: t('page.aicodecheck.gitlab.rule_check_status'),
      dataIndex: 'rule_check_status',
      key: 'rule_check_status',
      render: (text: number, record: TableListItem) => (
        <Switch
          checked={text === 1}
          checkedChildren="开启"
          unCheckedChildren="关闭"
          onChange={(checked) => handleRuleCheckStatusChange(record.id, checked)}
          loading={loadingId === record.id}
        />
      )
    },
    {
      title: t('app.table.action'),
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 220,
      render: (text, record: TableListItem) => (
        <Space size='small'>
          <Button className='btn-group-cell' size='small' type='link' onClick={() => handleDetail(record)}>
            {t('app.global.view')}
          </Button>
          <Button className='btn-group-cell' size='small' type='link' onClick={() => handleUpdate(record)}>
            {t('app.global.edit')}
          </Button>
          <Popconfirm
            open={deleteOpen === record.id}
            title={t('app.global.delete')}
            description={t('app.global.delete.tip')}
            onConfirm={() => deleteConfirm(record.id)}
            onCancel={deleteCancel}
            okText='Yes'
            cancelText='No'
          >
            <Button 
              danger 
              className='btn-group-cell' 
              onClick={() => handleDelete(record.id)} 
              size='small' 
              type='link'
              loading={deleteLoading === record.id}
            >
              {t('app.global.delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className='layout-main-conent'>
      <CommonTable
        ref={tableRef}
        columns={columns}
        queryList={queryList}
        title={
          <Button type='primary' onClick={handleCreate}>
            {t('page.aicodecheck.gitlab.add')}
          </Button>
        }
        scroll={{ x: 1500 }}
        useTools
      />

      {/* <CreateForm
        initialValues={updateData}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      /> */}
    </div>
  );
}

export default memo(observer(App));
