import { useQuery, useMutation } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { getSubscription, createInvoice } from '../api/client';
import { LoadingDots } from './LoadingDots';
import { debug } from '../lib/debug';

const checkIcon = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';
const crossIcon = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>';

export function SubscriptionPlans() {
  const { t } = useTranslation();

  const { data: sub, isLoading } = useQuery({
    queryKey: ['subscription'],
    queryFn: getSubscription,
  });

  const invoiceMutation = useMutation({
    mutationFn: (plan: string) => createInvoice({ plan }),
    onSuccess: (data) => {
      debug('[subscription] invoice created', data);
      window.open(data.invoice_link, '_blank');
    },
  });

  const plans = [
    {
      key: 'free',
      title: t('subscription.free_plan'),
      price: 0,
      features: [
        { text: t('subscription.feature_messages', { count: 10 }), included: true },
        { text: t('subscription.feature_basic_corrections'), included: true },
        { text: t('subscription.feature_unlimited'), included: false },
        { text: t('subscription.feature_priority'), included: false },
        { text: t('subscription.feature_ielts'), included: false },
        { text: t('subscription.feature_errors'), included: false },
        { text: t('subscription.feature_export'), included: false },
      ],
    },
    {
      key: 'weekly',
      title: t('subscription.weekly', { stars: 300 }),
      price: 300,
      features: [
        { text: t('subscription.feature_unlimited'), included: true },
        { text: t('subscription.feature_priority'), included: true },
        { text: t('subscription.feature_basic_corrections'), included: true },
        { text: t('subscription.feature_ielts'), included: true },
        { text: t('subscription.feature_errors'), included: false },
        { text: t('subscription.feature_export'), included: false },
      ],
    },
    {
      key: 'monthly',
      title: t('subscription.monthly', { stars: 800 }),
      price: 800,
      features: [
        { text: t('subscription.feature_unlimited'), included: true },
        { text: t('subscription.feature_priority'), included: true },
        { text: t('subscription.feature_basic_corrections'), included: true },
        { text: t('subscription.feature_ielts'), included: true },
        { text: t('subscription.feature_errors'), included: true },
        { text: t('subscription.feature_export'), included: true },
      ],
    },
  ];

  return (
    <div className="scroll-area">
      <div className="section-title">{t('subscription.title')}</div>

      {isLoading && (
        <div className="placeholder">
          <LoadingDots size={14} />
          <div>{t('common.loading')}</div>
        </div>
      )}

      {sub?.active && (
        <div style={{
          margin: '0 16px 16px',
          padding: 'var(--card-padding)',
          borderRadius: 'var(--round-lg)',
          border: '2px solid var(--c-primary)',
          background: 'var(--c-surface-container-lowest)',
          textAlign: 'center',
        }}>
          <div style={{ fontSize: 18, fontWeight: 600, marginBottom: 4, color: 'var(--c-primary)' }}>
            {t('subscription.active')}
          </div>
          <div style={{ fontSize: 13, color: 'var(--c-on-surface-variant)' }}>
            {sub.plan} — {sub.expires_at}
          </div>
        </div>
      )}

      {plans.map(plan => {
        const isMonthly = plan.key === 'monthly';
        return (
          <div
            key={plan.key}
            style={{
              margin: '0 16px 16px',
              padding: 'var(--card-padding)',
              borderRadius: 'var(--round-lg)',
              border: isMonthly ? '2px solid var(--c-secondary)' : '1px solid var(--c-outline-variant)',
              background: 'var(--c-surface-container-lowest)',
              position: 'relative',
              overflow: 'hidden',
            }}
          >
            {isMonthly && (
              <div style={{
                position: 'absolute',
                top: 12,
                right: -28,
                transform: 'rotate(45deg)',
                background: 'var(--c-secondary)',
                color: '#fff',
                fontSize: 10,
                fontWeight: 700,
                padding: '2px 32px',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
              }}>
                {t('subscription.best_value')}
              </div>
            )}
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
              <span style={{ fontSize: 16, fontWeight: 700, color: 'var(--c-on-surface)' }}>
                {plan.title}
              </span>
              {plan.price > 0 && (
                <span style={{
                  fontSize: 14,
                  fontWeight: 600,
                  color: 'var(--c-secondary)',
                }}>
                  ⭐ {plan.price}
                </span>
              )}
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 6, marginBottom: 16 }}>
              {plan.features.map((f, i) => (
                <div key={i} style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 8,
                  fontSize: 13,
                  color: f.included ? 'var(--c-on-surface)' : 'var(--c-outline)',
                }}>
                  <span
                    style={{
                      display: 'flex',
                      color: f.included ? 'var(--c-primary)' : 'var(--c-outline)',
                      flexShrink: 0,
                    }}
                    dangerouslySetInnerHTML={{ __html: f.included ? checkIcon : crossIcon }}
                  />
                  {f.text}
                </div>
              ))}
            </div>
            {plan.key !== 'free' && (
              <button
                className="btn btn-primary"
                style={{ width: '100%' }}
                onClick={() => invoiceMutation.mutate(plan.key)}
                disabled={invoiceMutation.isPending}
              >
                {invoiceMutation.isPending ? t('common.loading') : t('subscription.subscribe')}
              </button>
            )}
          </div>
        );
      })}
    </div>
  );
}
