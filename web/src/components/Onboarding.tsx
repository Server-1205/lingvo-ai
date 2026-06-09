import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { initData } from '@telegram-apps/sdk';

const slides = [
  {
    icon: '🤖',
    titleKey: 'onboarding.slide1_title',
    descKey: 'onboarding.slide1_desc',
    items: [
      'onboarding.slide1_item1',
      'onboarding.slide1_item2',
      'onboarding.slide1_item3',
    ],
  },
  {
    icon: '📚',
    titleKey: 'onboarding.slide2_title',
    descKey: 'onboarding.slide2_desc',
    items: [
      'onboarding.slide2_item1',
      'onboarding.slide2_item2',
      'onboarding.slide2_item3',
    ],
  },
  {
    icon: '📊',
    titleKey: 'onboarding.slide3_title',
    descKey: 'onboarding.slide3_desc',
    items: [
      'onboarding.slide3_item1',
      'onboarding.slide3_item2',
      'onboarding.slide3_item3',
    ],
  },
];

interface OnboardingProps {
  onDone: () => void;
}

export function Onboarding({ onDone }: OnboardingProps) {
  const { t } = useTranslation();
  const [slide, setSlide] = useState(0);
  const rawUser = initData.user();
  const userName = rawUser?.first_name ?? '';

  const current = slides[slide];
  const isLast = slide === slides.length - 1;

  const handleNext = () => {
    if (isLast) {
      localStorage.setItem('lingvo_onboarding_done', 'true');
      onDone();
    } else {
      setSlide(s => s + 1);
    }
  };

  const handleSkip = () => {
    localStorage.setItem('lingvo_onboarding_done', 'true');
    onDone();
  };

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100svh',
      background: 'linear-gradient(180deg, var(--c-primary) 0%, var(--c-surface) 100%)',
      padding: '24px 20px',
    }}>
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        {!isLast && (
          <button
            onClick={handleSkip}
            style={{
              background: 'rgba(255,255,255,0.2)',
              border: 'none',
              borderRadius: 'var(--round-full)',
              padding: '8px 16px',
              color: '#fff',
              fontSize: 14,
              fontWeight: 600,
              cursor: 'pointer',
            }}
          >
            {t('onboarding.skip')}
          </button>
        )}
      </div>

      <div style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        textAlign: 'center',
        gap: 16,
      }}>
        {slide === 0 && (
          <div style={{
            fontSize: 14,
            color: 'rgba(255,255,255,0.85)',
            fontWeight: 500,
            marginBottom: 4,
          }}>
            {t('onboarding.welcome', { name: userName })}
          </div>
        )}

        <div style={{ fontSize: 64, marginBottom: 8 }}>{current.icon}</div>

        <div style={{
          fontSize: 24,
          fontWeight: 700,
          color: slide === 2 ? 'var(--c-on-surface)' : '#fff',
          lineHeight: 1.3,
        }}>
          {t(current.titleKey)}
        </div>

        <div style={{
          fontSize: 15,
          color: slide === 2 ? 'var(--c-on-surface-variant)' : 'rgba(255,255,255,0.8)',
          lineHeight: 1.5,
          maxWidth: 300,
        }}>
          {t(current.descKey)}
        </div>

        <div style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 10,
          marginTop: 12,
          textAlign: 'left',
        }}>
          {current.items.map((item, i) => (
            <div key={i} style={{
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              fontSize: 15,
              color: slide === 2 ? 'var(--c-on-surface-variant)' : 'rgba(255,255,255,0.85)',
            }}>
              <span style={{ fontSize: 18, color: 'var(--c-primary-fixed-dim)' }}>✦</span>
              <span>{t(item)}</span>
            </div>
          ))}
        </div>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 20, paddingBottom: 20 }}>
        <div style={{ display: 'flex', gap: 8 }}>
          {slides.map((_, i) => (
            <div key={i} style={{
              width: i === slide ? 24 : 8,
              height: 8,
              borderRadius: 'var(--round-full)',
              background: i === slide
                ? (slide === 2 ? 'var(--c-primary)' : '#fff')
                : 'rgba(255,255,255,0.3)',
              transition: 'all 0.3s',
            }} />
          ))}
        </div>

        <button
          onClick={handleNext}
          style={{
            width: '100%',
            padding: '14px 0',
            fontSize: 17,
            fontWeight: 700,
            borderRadius: 'var(--round-full)',
            border: 'none',
            cursor: 'pointer',
            background: slide === 2 ? 'var(--c-primary)' : 'rgba(255,255,255,0.95)',
            color: slide === 2 ? '#fff' : 'var(--c-primary)',
          }}
        >
          {isLast ? t('onboarding.start') : t('onboarding.next')}
        </button>
      </div>
    </div>
  );
}
