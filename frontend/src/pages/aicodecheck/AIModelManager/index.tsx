import { memo, useContext, useRef, useState, useEffect } from 'react';
import { Button, FormInstance, message, Popconfirm, PopconfirmProps, Space, Switch } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { getModelOptions, removeData } from './service';

import { createData, queryList, updateData as updateDataService } from './service';
import { TableQueryParam, TableListItem } from './data';
import CreateForm from './components/CreateForm';

const getCurrentModel = (modelOptions: any[], type: string, model: string) => {
  return modelOptions.find((item: any) => item.type === type && item.model === model);
}

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [modelOptions, setModelOptions] = useState<any[]>([]);
  const [typeOptions, setTypeOptions] = useState<any[]>([]);

  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  // 删除
  const [deleteOpen, setDeleteOpen] = useState<number | undefined>();
  const [deleteLoadingId, setDeleteLoadingId] = useState<number | null>(null);
  const handleDelete = (id: number) => setDeleteOpen(id);
  const deleteConfirm = (id: number) => {
    setDeleteLoadingId(id);
    removeData(id)
      .then(() => {
        message.success(t('app.global.tip.delete.success'));
        reload();
        setDeleteOpen(void 0);
      })
      .finally(() => {
        setDeleteLoadingId(null);
      });
  };

  const getAIModelOptions = async () => {
    const res = await getModelOptions();
    const { ret_code, data, message: errorMessage } = res;
    if (ret_code === 0) {
      // Extract unique types and create type options
      const uniqueTypes = [...new Set(data.map((item: any) => item.type))];
      const typeOpts = uniqueTypes.map(type => ({
        label: type,
        value: type
      }));
      setTypeOptions(typeOpts);
      setModelOptions(data);
    } else {
      message.error(errorMessage);
    }
  };

  const deleteCancel: PopconfirmProps['onCancel'] = () => {
    setDeleteOpen(void 0);
  };

  // 新增&编辑
  const [createSubmitLoading, setCreateSubmitLoading] = useState<boolean>(false);
  const [createFormVisible, setCreateFormVisible] = useState<boolean>(false);
  const [updateData, setUpdateData] = useState<Partial<TableQueryParam>>({});
  const [loadingId, setLoadingId] = useState<number | null>(null);

  const handleCreate = () => {
    setUpdateData({});
    setCreateFormVisible(true);
  };
  const createSubmit = (values: any, form: FormInstance) => {
    setCreateSubmitLoading(true);
    const request = updateData.id ? updateDataService : createData;
    const currentModel = getCurrentModel(modelOptions, values.type, values.model);
    const data = {
      ...values,
      api_url: currentModel?.api_url,
      id: updateData?.id,
    }
    request(data)
      .then(() => {
        form.resetFields();
        setCreateFormVisible(false);
        message.success(values.id ? t('app.global.tip.update.success') : t('app.global.tip.create.success'));
        reload();

        setCreateSubmitLoading(false);
      })
      .catch(() => {
        setCreateSubmitLoading(false);
      });
  };

  const handleRuleCheckStatusChange = async (record: TableListItem, checked: boolean) => {
    setLoadingId(record.id);
    updateDataService({ id: record.id, is_active: checked ? 1 : 2 })
      .then(() => {
        message.success(t('app.global.tip.update.success'));
        reload();
      })
      .finally(() => {
        setLoadingId(null);
      });
  };

  const handleUpdate = (record: TableListItem) => {
    setUpdateData({
      ...record,
    });
    setCreateFormVisible(true);
  };

  useEffect(() => {
    getAIModelOptions();
  }, []);

  const columns: ColumnsType<TableListItem> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: t('page.aicodecheck.aimodel.type'),
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: t('page.aicodecheck.aimodel.api_url'),
      dataIndex: 'api_url',
      key: 'api_url',
    },
    {
      title: t('page.aicodecheck.aimodel.model'),
      dataIndex: 'model',
      key: 'model'
    },
    {
      title: t('page.aicodecheck.aimodel.using'),
      dataIndex: 'is_active',
      key: 'is_active',
      render: (text: number, record: TableListItem) => (
        <Switch
          checked={text === 1}
          checkedChildren="开启"
          unCheckedChildren="关闭"
          onChange={(checked) => handleRuleCheckStatusChange(record, checked)}
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
            okButtonProps={{ loading: deleteLoadingId === record.id }}
          >
            <Button
              danger
              className='btn-group-cell'
              onClick={() => handleDelete(record.id)}
              size='small'
              type='link'
              loading={deleteLoadingId === record.id}
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
            {t('page.aicodecheck.aimodel.add')}
          </Button>
        }
        useTools
      />

      <CreateForm
        initialValues={updateData}
        modelOptions={modelOptions}
        typeOptions={typeOptions}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      />
    </div>
  );
}

export default memo(observer(App));
