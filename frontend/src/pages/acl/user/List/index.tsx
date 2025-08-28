import { memo, useContext, useRef, useState } from 'react';
import { Button, FormInstance, message, Popconfirm, PopconfirmProps, Space } from 'antd';
import type { TableColumnsType } from 'antd';
import { observer } from 'mobx-react-lite';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { FormType } from '@/@types/enum';

import CreateForm from './components/CreateForm';
import Preview from './components/Preview';
import { createData, queryList, removeData, updateData as updateDataService } from './service';
import { TableListItem } from './data.d';

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  // 删除
  const [deleteOpen, setDeleteOpen] = useState<number | undefined>();
  const handleDelete = (id: number) => setDeleteOpen(id);
  const deleteConfirm = (id: number, user: string) => {
    removeData(id, user).then(() => {
      message.success(t('app.global.tip.delete.success'));
      reload();
      setDeleteOpen(void 0);
    });
  };

  const deleteCancel: PopconfirmProps['onCancel'] = () => {
    setDeleteOpen(void 0);
  };

  // 新增&编辑
  const [createSubmitLoading, setCreateSubmitLoading] = useState<boolean>(false);
  const [createFormVisible, setCreateFormVisible] = useState<boolean>(false);
  const [updateData, setUpdateData] = useState<Partial<TableListItem>>({});

  const [previewVisible, setPreviewVisible] = useState<boolean>(false);
  const [previewData, setPreviewData] = useState<Partial<TableListItem>>({});

  const handleCreate = () => {
    setUpdateData({});
    setCreateFormVisible(true);
  };

  const createSubmit = async (values: TableListItem, form: FormInstance) => {
    setCreateSubmitLoading(true);
    const request = updateData.id ? updateDataService : createData;
    request({ ...values, id: updateData.id as number })
      .then(() => {
        form.resetFields();
        setCreateFormVisible(false);
        message.success(t(values.id ? 'app.global.tip.update.success' : 'app.global.tip.create.success'));
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
      password: '',
    });
    setCreateFormVisible(true);
  };

  const handlePreview = (record: TableListItem) => {
    setPreviewData({
      ...record,
      password: '',
    });
    setPreviewVisible(true);
  };

  const columns: TableColumnsType<TableListItem> = [
    {
      title: 'Id',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: t('page.user.enname'),
      dataIndex: 'user',
      key: 'user',
    },
    {
      title: t('page.user.cnname'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('page.user.job'),
      dataIndex: 'job',
      key: 'job',
    },
    {
      title: t('page.user.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: t('page.user.phone'),
      dataIndex: 'phone_number',
      key: 'phone_number',
    },
    {
      title: t('app.table.action'),
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 150,
      render: (text, record: TableListItem) => (
        <Space size='small'>
          <Button className='btn-group-cell' size='small' type='link' onClick={() => handlePreview(record)}>
            {t('app.global.view')}
          </Button>
          <Button className='btn-group-cell' size='small' type='link' onClick={() => handleUpdate(record)}>
            {t('app.global.edit')}
          </Button>
          <Popconfirm
            open={deleteOpen === record.id}
            title={t('app.global.delete')}
            description={t('app.global.delete.tip')}
            onConfirm={() => deleteConfirm(record.id, record.user)}
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

  const formItems = [
    {
      label: t('page.user.enname'),
      name: 'user',
      type: FormType.Input,
      span: 8,
    },
    {
      label: t('page.user.cnname'),
      name: 'userName',
      type: FormType.Input,
      span: 8,
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
            {t('page.user.add')}
          </Button>
        }
        filterFormItems={formItems}
        useTools
      />
      <CreateForm
        initialValues={updateData}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      />
      <Preview visible={previewVisible} setVisible={setPreviewVisible} data={previewData} />
    </div>
  );
}

export default memo(observer(App));
