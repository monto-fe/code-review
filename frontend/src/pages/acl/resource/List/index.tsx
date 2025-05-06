import { memo, useContext, useRef, useState } from 'react';
import { Button, FormInstance, message, Popconfirm, PopconfirmProps, Space } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';

import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';
import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { FormType } from '@/@types/enum';

import { createData, queryList, removeData, updateData as updateDataService } from '../service';
import { TableListItem } from '../data.d';
import { resourceCategoryMap } from '../const';

import CreateForm from './components/CreateForm';

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
    removeData(id).then(() => {
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

  const handleCreate = () => {
    setUpdateData({});
    setCreateFormVisible(true);
  };

  const createSubmit = async (values: TableListItem, form: FormInstance) => {
    setCreateSubmitLoading(true);
    const request = updateData.id ? updateDataService : createData;
    request({ ...values, id: updateData.id as number, category: resourceCategoryMap.Other })
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

  const handleUpdate = async (record: TableListItem) => {
    setUpdateData({
      ...record,
    });
    setCreateFormVisible(true);
  };

  const columns: ColumnsType<TableListItem> = [
    {
      title: 'id',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: t('page.resource.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('page.resource.key'),
      dataIndex: 'resource',
      key: 'resource',
    },
    {
      title: t('page.resource.describe'),
      dataIndex: 'describe',
      key: 'describe',
    },
    {
      title: t('app.table.operator'),
      dataIndex: 'operator',
      key: 'operator',
    },
    {
      title: t('app.table.updatetime'),
      dataIndex: 'update_time',
      key: 'update_time',
      width: 200,
      render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time),
    },
    {
      title: t('app.table.action'),
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 100,
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

  const formItems = [
    {
      label: t('page.resource.name'),
      name: 'name',
      type: FormType.Input,
      span: 8,
    },
    {
      label: t('page.resource.key'),
      name: 'resource',
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
            {t('page.resource.add')}
          </Button>
        }
        useTools
        filterFormItems={formItems}
        scroll={{ x: 1200 }}
      />

      <CreateForm
        initialValues={updateData}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      />
    </div>
  );
}

export default memo(observer(App));
