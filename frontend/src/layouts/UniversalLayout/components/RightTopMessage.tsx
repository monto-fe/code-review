import { memo } from 'react';
import { Popover } from 'antd';
// import { BasicContext } from '@/store/context';
// import { useI18n } from '@/store/i18n';

export default memo(() => (
  // const context = useContext(BasicContext) as any;
  // const { i18nLocale } = context.storeContext;
  // const t = useI18n(i18nLocale);

  <Popover content={<div>二维码展示区</div>} title='微信群'>
    <div className='universallayout-top-notocemenu ant-dropdown-link cursor' onClick={(e) => e.preventDefault()}>
      技术交流
    </div>
  </Popover>
));
