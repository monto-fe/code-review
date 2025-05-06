import { getLocale, defaultLang } from '@/utils/i18n';
import { I18nKey, I18nVal } from '@/@types/i18n.d';

// 导入 全局自定义 语言
import globalLocales from '@/locales';

// 导入 antd 语言
import zhCN from 'antd/es/locale/zh_CN';
import zhTW from 'antd/es/locale/zh_TW';
import enUS from 'antd/es/locale/en_US';

// antd 语言包
export const antdMessages: { [key in I18nKey]: any } = {
  'zh-CN': zhCN,
  'zh-TW': zhTW,
  'en-US': enUS,
};

// 语言配置
const sysLocale = getLocale();
export const i18nLocaleDefault = antdMessages[sysLocale] ? sysLocale : defaultLang;
export const getI18nLocale: (i18n: I18nKey) => I18nVal = (i18n: I18nKey) =>
  (globalLocales[i18n] || globalLocales['zh-CN']) as I18nVal;
export const getAntdI18nMessage = (i18n: I18nKey) => antdMessages[i18n] || antdMessages['zh-CN'];

export const useI18n = (i18n: I18nKey) => (key: string) =>
  getI18nLocale(i18n)[key] || getAntdI18nMessage(i18n)[key] || key;
