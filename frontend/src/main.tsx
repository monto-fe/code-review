import ReactDOM from 'react-dom/client';
import { BrowserRouter /* , HashRouter */ } from 'react-router-dom';

import 'dayjs/locale/zh-cn';

// Register icon sprite
import 'virtual:svg-icons-register';
// 全局css
import '@/assets/css/index.less';
// App
import App from '@/App';
import { BasicContext } from '@/store/context';
import { observerRoot } from '@/store';

// 挂载
ReactDOM.createRoot(document.getElementById('app') as HTMLElement).render(
  <BrowserRouter>
    <BasicContext.Provider value={{ storeContext: observerRoot }}>
      <App />
    </BasicContext.Provider>
  </BrowserRouter>,
);
