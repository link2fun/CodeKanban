import { useI18n } from 'vue-i18n';
import dayjs from 'dayjs';

export type LocaleType = 'zh-CN' | 'en-US';

export function useLocale() {
  const { locale, t } = useI18n();

  const setLocale = (newLocale: LocaleType) => {
    locale.value = newLocale;
    localStorage.setItem('app-locale', newLocale);

    // 同步更新 dayjs 的语言
    dayjs.locale(newLocale === 'zh-CN' ? 'zh-cn' : 'en');
  };

  const toggleLocale = () => {
    const newLocale: LocaleType = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN';
    setLocale(newLocale);
  };

  return {
    locale,
    setLocale,
    toggleLocale,
    t,
  };
}
