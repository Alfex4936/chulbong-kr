import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// Import your translation files
import translationEN from './locales/en/translation.json';
import translationKO from './locales/ko/translation.json';

// the translations
const resources = {
    en: {
        translation: translationEN,
    },
    ko: {
        translation: translationKO,
    },
};

i18n
    // detect the user language
    // .use(LanguageDetector)
    // pass the i18n instance to the react-i18next components.
    .use(initReactI18next)
    // init i18next
    //  https://www.i18next.com/overview/configuration-options
    .init({
        resources,
        fallbackLng: 'en',
        debug: true,

        interpolation: {
            escapeValue: false, // not needed for react as it escapes by default
        },
    });

export default i18n;
