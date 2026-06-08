import { useTranslation } from 'react-i18next';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const toggleLang = () => {
    const next = i18n.language === 'uz' ? 'ru' : 'uz';
    i18n.changeLanguage(next);
  };

  return (
    <button
      onClick={toggleLang}
      style={{
        padding: '6px 14px',
        border: '1px solid rgba(255,255,255,0.3)',
        borderRadius: 10,
        background: 'rgba(255,255,255,0.15)',
        color: '#fff',
        fontSize: 13,
        fontWeight: 500,
        cursor: 'pointer',
      }}
    >
      {i18n.language === 'uz' ? 'RU' : 'UZ'}
    </button>
  );
}
