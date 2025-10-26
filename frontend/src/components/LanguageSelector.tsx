import { useState, useEffect } from 'react';
import { Globe, Check } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { GetLanguage, SetLanguage } from '../../wailsjs/go/main/App';

export function LanguageSelector() {
  const { i18n } = useTranslation();
  const [currentLang, setCurrentLang] = useState('pt-BR');
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    // Load language from backend on mount
    loadLanguage();
  }, []);

  const loadLanguage = async () => {
    try {
      const lang = await GetLanguage();
      if (lang && (lang === 'pt-BR' || lang === 'en')) {
        setCurrentLang(lang);
        i18n.changeLanguage(lang);
      }
    } catch (error) {
      console.error('Failed to load language:', error);
    }
  };

  const changeLanguage = async (lang: string) => {
    try {
      // Update backend
      await SetLanguage(lang);
      
      // Update frontend
      i18n.changeLanguage(lang);
      setCurrentLang(lang);
      setIsOpen(false);
    } catch (error) {
      console.error('Failed to change language:', error);
    }
  };

  const languages = [
    { code: 'pt-BR', name: 'Português' },
    { code: 'en', name: 'English' },
  ];

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="group relative w-9 h-9 rounded-lg bg-white/80 dark:bg-gray-800/80 border border-gray-200/50 dark:border-gray-700/50 flex items-center justify-center transition-all duration-300 hover:shadow-sm hover:scale-105"
        title="Change language"
      >
        <Globe className="w-3.5 h-3.5 text-gray-600 dark:text-gray-400 transition-colors duration-300" />
      </button>

      {isOpen && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 z-40"
            onClick={() => setIsOpen(false)}
          />
          
          {/* Dropdown */}
          <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 z-50 overflow-hidden">
            {languages.map((lang) => (
              <button
                key={lang.code}
                onClick={() => changeLanguage(lang.code)}
                className={`w-full px-4 py-3 flex items-center gap-3 transition-colors ${
                  currentLang === lang.code
                    ? 'bg-purple-50 dark:bg-purple-900/20 text-purple-700 dark:text-purple-300'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/50'
                }`}
              >
                <Globe className="w-4 h-4" />
                <span className="text-sm font-medium flex-1 text-left">{lang.name}</span>
                {currentLang === lang.code && (
                  <Check className="w-4 h-4 text-purple-600 dark:text-purple-400" />
                )}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
