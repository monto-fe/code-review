import { memo, useContext, useEffect, useRef, useState } from 'react';
import { Button, FormInstance, message, Popconfirm, PopconfirmProps, Popover, Space } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';

import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';
import { ResponseData } from '@/utils/request';
import { createData, queryList, removeData, updateData as updateDataService } from './service';
import { queryList as queryRoleList } from '@/pages/acl/role/service';
import { TableListItem as RoleTableListItem } from '@/pages/acl/role/data.d';
import CommonTable from '@/pages/component/Table';
import { ITable } from '@/pages/component/Table/data';
import Permission from '@/components/Permission';
import { FormType } from '@/@types/enum';
import { renderDateFromTimestamp, timeFormatType } from '@/utils/timeformat';

import CreateForm from './components/CreateForm';
import Preview from './components/Preview';
import { TableListItem } from './data.d';

function App() {
  const tableRef = useRef<ITable<TableListItem>>();
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [roleList, setRoleList] = useState<RoleTableListItem[]>([]);

  const getRoleList = () => {
    queryRoleList({
      current: 1,
      pageSize: 99999,
    }).then((response: ResponseData<RoleTableListItem[]>) => {
      if (response) {
        setRoleList(response.data.data || []);
      }
    });
  };

  const reload = () => tableRef.current && tableRef.current.reload && tableRef.current.reload();

  useEffect(() => {
    getRoleList();
  }, []);

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
      password: '※※※※※※',
      role_ids: (record.role || []).map((role) => role.id),
    });
    setCreateFormVisible(true);
  };

  const handlePreview = (record: TableListItem) => {
    setPreviewData({
      ...record,
      password: '※※※※※※',
      role_ids: (record.role || []).map((role) => role.id),
    });
    setPreviewVisible(true);
  };

  const columns: ColumnsType<TableListItem> = [
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
      title: t('page.user.role'),
      dataIndex: 'roleName',
      key: 'roleName',
      width: 120,
      render: (text, record: TableListItem) => (
        <div>
          {Array.isArray(record.role) ? (
            <span>
              {record.role.length > 1 ? (
                <Popover
                  placement='right'
                  content={
                    <div>
                      {record.role.map((item: any) => (
                        <div key={item.role}>{`${item.name} (${item.role})`}</div>
                      ))}
                    </div>
                  }
                  trigger='hover'
                >
                  {`${record.role[0]?.name}...`}
                </Popover>
              ) : (
                record.role[0]?.name
              )}
            </span>
          ) : (
            ''
          )}
        </div>
      ),
      filters: roleList.map((item) => ({
        text: `${item.name} (${item.role})`,
        value: item.name,
      })),
      filterMultiple: false,
      filterSearch: true,
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
      width: 200,
    },
    // {
    //   title: t('app.table.updatetime'),
    //   dataIndex: 'update_time',
    //   key: 'update_time',
    //   width: 200,
    //   sorter: true,
    //   render: (text: number) => renderDateFromTimestamp(text, timeFormatType.time),
    // },
    {
      title: t('app.table.action'),
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      width: 150,
      render: (text, record: TableListItem) => (
        <Permission
          role='admin'
          noNode={
            <Button className='btn-group-cell' size='small' type='link' onClick={() => handlePreview(record)}>
              {t('app.global.view')}
            </Button>
          }
        >
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
              <Button
                danger
                className='btn-group-cell'
                onClick={() => handleDelete(record.id)}
                size='small'
                type='link'
              >
                {t('app.global.delete')}
              </Button>
            </Popconfirm>
          </Space>
        </Permission>
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
          <Permission role='admin' noNode={null}>
            <Button type='primary' onClick={handleCreate}>
              {t('page.user.add')}
            </Button>
          </Permission>
        }
        filterFormItems={formItems}
        useTools
        scroll={{ x: 1200 }}
      />
      <CreateForm
        initialValues={updateData}
        visible={createFormVisible}
        setVisible={setCreateFormVisible}
        roleList={roleList}
        onSubmit={createSubmit}
        onSubmitLoading={createSubmitLoading}
      />
      <Preview visible={previewVisible} setVisible={setPreviewVisible} data={previewData} roleList={roleList} />
    </div>
  );
}

export default memo(observer(App));
