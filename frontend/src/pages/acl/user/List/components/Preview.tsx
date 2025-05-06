import { useContext } from 'react';
import { Button, Descriptions, Modal } from 'antd';
import { TableListItem as RoleTableListItem } from '@/pages/acl/role/data.d';
import { BasicContext } from '@/store/context';
import { useI18n } from '@/store/i18n';

import { TableListItem } from '../data.d';

interface PreviewProps {
  visible: boolean;
  setVisible: Function;
  data: Partial<TableListItem>;
  roleList: RoleTableListItem[];
}

function Preview(props: PreviewProps) {
  const { visible, setVisible, data } = props;

  const context = useContext(BasicContext) as any;
  const { i18nLocale } = context.storeContext;
  const t = useI18n(i18nLocale);

  const handleCancel = () => {
    setVisible(false);
  };

  return (
    <Modal
      open={visible}
      title={t('page.user.view')}
      width={800}
      onCancel={handleCancel}
      footer={[
        <Button key='close' onClick={handleCancel}>
          {t('app.global.close')}
        </Button>,
      ]}
    >
      <Descriptions className='mt-32' title='' layout='vertical' bordered size='small'>
        <Descriptions.Item label={t('page.user.enname')}>{data?.user || '-'}</Descriptions.Item>
        <Descriptions.Item label={t('page.user.cnname')}>{data?.name || '-'}</Descriptions.Item>
        <Descriptions.Item label={t('page.user.job')}>{data?.job || '-'}</Descriptions.Item>
        <Descriptions.Item label={t('page.user.email')}>{data?.email || '-'}</Descriptions.Item>
        <Descriptions.Item label={t('page.user.phone')}>{data?.phone_number || '-'}</Descriptions.Item>
        <Descriptions.Item label={t('page.user.role')}>
          {Array.isArray(data?.role) ? (
            <span>
              {data?.role?.length ? (
                <div>
                  {data?.role.map((item: any) => (
                    <div key={item.role}>{`${item.name} (${item.role})`}</div>
                  ))}
                </div>
              ) : (
                '-'
              )}
            </span>
          ) : (
            ''
          )}
        </Descriptions.Item>
      </Descriptions>
    </Modal>
  );
}

export default Preview;
