import { useContext, useEffect, useState } from 'react';
import { Card, Table } from 'antd';
import { ColumnsType, TablePaginationConfig } from 'antd/lib/table';
import { observer } from 'mobx-react-lite';

import { useI18n } from '@/store/i18n';

import { TableListItem } from './data.d';
import { PaginationConfig } from '../../data';

import styles from '../../index.module.less';

import { ResponseData } from '@/utils/request';
import { hotTagsQueryList } from '../../service';
import { BasicContext } from '@/store/context';

const initPagination = {
  total: 0,
  current: 1,
  pageSize: 5,
  showSizeChanger: false,
};

function HotTagsCard() {
  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);
  const [loading, setLoading] = useState<boolean>(false);
  const [list, setList] = useState<TableListItem[]>([]);
  const [pagination, setPagination] = useState<PaginationConfig>({
    ...initPagination,
  });

  const getList = async (current: number) => {
    setLoading(true);
    try {
      const response: ResponseData<TableListItem[]> = await hotTagsQueryList();
      const {
        data: { data, total },
      } = response;
      setList(data || []);
      setPagination({
        ...initPagination,
        current,
        total: total || 0,
      });
    } catch (error: any) {
      console.log(error);
    }
    setLoading(false);
  };

  useEffect(() => {
    getList(1);
  }, []);

  const columns: ColumnsType<TableListItem> = [
    {
      title: t('page.home.hottagscard.card.table-column-number'),
      dataIndex: 'index',
      width: 80,
      render: (_, record, index) => <>{(pagination.current - 1) * pagination.pageSize + index + 1}</>,
    },
    {
      title: t('page.home.hottagscard.card.table-column-name'),
      dataIndex: 'name',
    },
    {
      title: t('page.home.hottagscard.card.table-column-hit'),
      dataIndex: 'hit',
    },
  ];

  return (
    <Card className={styles.homeBoxCard} title={t('page.home.hottagscard.card-title')}>
      <Table
        size='small'
        rowKey='id'
        columns={columns}
        dataSource={list}
        loading={loading}
        pagination={pagination}
        onChange={(p: TablePaginationConfig) => {
          getList(p.current || 1);
        }}
      />
    </Card>
  );
}

export default observer(HotTagsCard);
