import { useEffect, useState, forwardRef, useImperativeHandle, Ref, useContext } from 'react';
import { Card, Dropdown, Flex, Input, Space, Table, Tooltip } from 'antd';
import { SizeType } from 'antd/lib/config-provider/SizeContext';
import { ColumnHeightOutlined, ReloadOutlined } from '@ant-design/icons';
import { TableProps } from 'antd/lib/table/InternalTable';
import { SorterResult } from 'antd/lib/table/interface';

import Filters from './Filter';
import { ITable, PaginationConfig } from './data';
import { ResponseData } from '@/utils/request';
import { AnyObject } from 'antd/lib/_util/type';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

export const defaultCurrent = 1;
export const defaultPageSize = 10;

function CommonTable<T extends AnyObject>(props: ITable<T>, ref: Ref<unknown> | undefined) {
  const {
    queryList,
    columns,
    title,
    rowKey,
    useTools = true,
    fuzzySearchKey,
    fuzzySearchPlaceholder,
    filterFormItems,
    scroll,
  } = props;
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const [list, setList] = useState<T[]>([]);
  const [size, setSize] = useState<SizeType>('middle');
  const [filter, setFilter] = useState<any>();
  const [externalFilter, setExternalFilter] = useState<any>();
  const [pagination, setPagination] = useState<PaginationConfig>({
    total: 0,
    current: defaultCurrent,
    pageSize: defaultPageSize,
    showSizeChanger: true,
    showQuickJumper: true,
  });
  const [loading, setLoading] = useState<boolean>(false);
  const [fuzzySearch, setFuzzySearch] = useState<string>('');

  const getList = (current: number, pageSize: number, filter?: any, externalFilter?: any): void => {
    setLoading(true);
    const params = {
      current,
      pageSize: pageSize || defaultPageSize,
      ...filter,
      ...externalFilter,
    };

    if (fuzzySearch && fuzzySearchKey) {
      params[fuzzySearchKey] = fuzzySearch;
    }

    queryList(params)
      .then((response: ResponseData<T[]>) => {
        if (response) {
          console.log("response.data", response.data)
          setList(Array.isArray(response.data.data) ? response.data.data : []);
          setPagination({
            ...pagination,
            current,
            pageSize,
            total: response.data.count || 0,
          });
          setFilter(filter);
        }

        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  };

  const reload = () => {
    getList(pagination.current, pagination.pageSize, filter, externalFilter);
  };

  const rightTools = (
    <Space size='middle'>
      <Tooltip title={<span>{t('app.form.refresh')}</span>}>
        <ReloadOutlined onClick={reload} />
      </Tooltip>
      <Dropdown
        menu={{
          selectedKeys: size ? [size] : undefined,
          items: [
            {
              label: t('app.table.large'),
              key: 'large',
            },
            {
              label: t('app.table.middle'),
              key: 'middle',
            },
            {
              label: t('app.table.small'),
              key: 'small',
            },
          ],
          onClick: (value) => setSize(value.key as SizeType),
        }}
        trigger={['click']}
      >
        <Tooltip title={t('app.table.lineheight')}>
          <ColumnHeightOutlined />
        </Tooltip>
      </Dropdown>
    </Space>
  );

  const handleTableChange: TableProps['onChange'] = (pagination, filters, sorter) => {
    const params: any = {
      order: (sorter as SorterResult<any>).order,
      field: (sorter as SorterResult<any>).field,
    };

    if (filters) {
      Object.getOwnPropertyNames(filters).forEach((key) => {
        if (Array.isArray(filters[key]) && (filters[key] as Array<string>).length) {
          params[key] = (filters[key] as Array<string>)[0] as string;
        }
      });
    }
    getList(pagination.current || defaultCurrent, pagination.pageSize || defaultPageSize, params, externalFilter);
  };

  const handleSearch = (values: any) => {
    setExternalFilter(values);
    getList(pagination.current, pagination.pageSize, filter, values);
  };

  const handleFuzzySearch = () => {
    getList(pagination.current, pagination.pageSize, filter, externalFilter);
  };

  useEffect(() => {
    getList(pagination.current, pagination.pageSize);
  }, []);

  useImperativeHandle(ref, () => ({
    reload,
  }));

  console.log('pagination', pagination);

  return (
    <>
      {filterFormItems ? (
        <Card variant="outlined" className='mb-16'>
          <Filters items={filterFormItems} size={size} handleSearch={handleSearch} />
        </Card>
      ) : null}
      <Card
        variant="outlined"
        title={title}
        extra={
          <Flex align='center'>
            {fuzzySearchKey ? (
              <Input.Search
                placeholder={fuzzySearchPlaceholder || t('app.table.fuzzysearch')}
                style={{ width: '270px', margin: '0px 16px' }}
                onChange={(e) => setFuzzySearch(e.target?.value)}
                onSearch={handleFuzzySearch}
              />
            ) : null}
            {useTools ? rightTools : null}
          </Flex>
        }
      >
        <Table
          rowKey={rowKey || 'id'}
          columns={columns}
          dataSource={list}
          loading={loading}
          scroll={scroll}
          size={size}
          pagination={{
            ...pagination
          }}
          onChange={handleTableChange}
        />
      </Card>
    </>
  );
}

export default forwardRef(CommonTable);
