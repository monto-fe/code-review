import { memo, useContext, useRef, useState } from 'react';
import { Button, FormInstance, message, Popconfirm, PopconfirmProps, Space } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';

import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { createData, MODEL_TYPE_OPTIONS, queryList, updateData as updateDataService } from './service';
import { TableQueryParam, TableListItem } from './data';
import CreateForm from './components/CreateForm';

const initialValues = {
    name: 'DeepSeek',
    api_url: 'https://deepseek.modelverse.cn/v1/chat/completions',
    api_key: '',
    model: 'deepseek-ai/DeepSeek-V3-0324'
}

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  // 删除
  const [deleteOpen, setDeleteOpen] = useState<number | undefined>();
  const handleDelete = (id: number) => setDeleteOpen(id);
  const deleteConfirm = (id: number) => {
    message.success(t('app.global.doing'));
    // removeData(id).then(() => {
    //   message.success(t('app.global.tip.delete.success'));
    //   reload();
    //   setDeleteOpen(void 0);
    // });
  };

  const deleteCancel: PopconfirmProps['onCancel'] = () => {
    setDeleteOpen(void 0);
  };

  // 新增&编辑
  const [createSubmitLoading, setCreateSubmitLoading] = useState<boolean>(false);
  const [createFormVisible, setCreateFormVisible] = useState<boolean>(false);
  const [updateData, setUpdateData] = useState<Partial<TableQueryParam>>({});

  const handleCreate = () => {
    setUpdateData({});
    setCreateFormVisible(true);
  };
  const createSubmit = (values: TableListItem, form: FormInstance) => {
    setCreateSubmitLoading(true);
    const request = updateData.id ? updateDataService : createData;
    request({ ...values, id: updateData.id as number })
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

  const handleUpdate = (record: TableListItem) => {
    setUpdateData({
      ...record,
    });
    setCreateFormVisible(true);
  };

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
      render: (text: string) => {
        return MODEL_TYPE_OPTIONS.find(item => item.value === text)?.label;
      }
    },
    {
      title: t('page.aicodecheck.aimodel.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('page.aicodecheck.aimodel.api_url'),
      dataIndex: 'api_url',
      key: 'api_url',
    },
    {
      title: t('page.aicodecheck.aimodel.api_key'),
      dataIndex: 'api_key',
      key: 'api_key',
    },
    {
      title: t('page.aicodecheck.aimodel.model'),
      dataIndex: 'model',
      key: 'model'
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
          >
            <Button danger className='btn-group-cell' onClick={() => handleDelete(record.id)} size='small' type='link'>
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
        initialValues={initialValues}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      />
    </div>
  );
}

export default memo(observer(App));
