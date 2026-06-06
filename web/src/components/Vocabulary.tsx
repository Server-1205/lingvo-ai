import { useTranslation } from 'react-i18next';

export function Vocabulary() {
  const { t } = useTranslation();

  return (
    <div className="scroll-area">
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '12px 16px',
      }}>
        <span className="section-title" style={{ padding: 0 }}>{t('vocab.title')}</span>
        <button className="btn btn-primary" style={{ padding: '6px 14px', fontSize: 13 }}>
          {t('vocab.add')}
        </button>
      </div>
      <div className="placeholder">
        <div className="placeholder-icon">📚</div>
        <div>{t('vocab.empty')}</div>
      </div>
    </div>
  );
}
