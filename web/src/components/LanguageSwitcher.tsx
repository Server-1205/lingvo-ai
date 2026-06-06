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
        padding: '6px 12px',
        border: '1px solid var(--tg-border)',
        borderRadius: 8,
        background: 'var(--tg-secondary-bg)',
        color: 'var(--tg-text)',
        fontSize: 13,
        fontWeight: 500,
        cursor: 'pointer',
      }}
    >
      {i18n.language === 'uz' ? 'RU' : 'UZ'}
    </button>
  );
}
