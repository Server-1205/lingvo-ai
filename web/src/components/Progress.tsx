import { useTranslation } from 'react-i18next';

export function ProgressView() {
  const { t } = useTranslation();

  return (
    <div className="scroll-area">
      <div className="section-title">{t('progress.title')}</div>

      <div className="card" style={{ display: 'flex', justifyContent: 'space-around', padding: 20 }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>A1</div>
          <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.level')}</div>
        </div>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>0</div>
          <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.messages')}</div>
        </div>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>0</div>
          <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.words')}</div>
        </div>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: 28, fontWeight: 700, color: 'var(--tg-button)' }}>0</div>
          <div style={{ fontSize: 12, color: 'var(--tg-hint)', marginTop: 4 }}>{t('progress.streak')}</div>
        </div>
      </div>
    </div>
  );
}
