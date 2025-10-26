import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// Import translations
import commonPtBR from '../locales/pt-BR/common.json';
import eventsPtBR from '../locales/pt-BR/events.json';
import settingsPtBR from '../locales/pt-BR/settings.json';
import welcomePtBR from '../locales/pt-BR/welcome.json';

import commonEN from '../locales/en/common.json';
import eventsEN from '../locales/en/events.json';
import settingsEN from '../locales/en/settings.json';
import welcomeEN from '../locales/en/welcome.json';

const resources = {
  'pt-BR': {
    common: commonPtBR,
    events: eventsPtBR,
    settings: settingsPtBR,
    welcome: welcomePtBR,
  },
  en: {
    common: commonEN,
    events: eventsEN,
    settings: settingsEN,
    welcome: welcomeEN,
  },
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: 'pt-BR', // Default language (will be overridden by saved preference)
    fallbackLng: 'en',
    defaultNS: 'common',
    interpolation: {
      escapeValue: false, // React already escapes
    },
  });

export default i18n;
