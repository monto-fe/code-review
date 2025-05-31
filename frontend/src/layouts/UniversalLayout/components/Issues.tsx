import { memo } from 'react';

export default memo(() => (
  <div className='universallayout-top-notocemenu ant-dropdown-link cursor' onClick={(e) => e.preventDefault()}>
    <a href='https://github.com/monto-fe/code-review/issues' target='_blank'>
      Issues
    </a>
  </div>
));
