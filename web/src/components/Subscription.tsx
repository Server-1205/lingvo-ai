import { useQuery, useMutation } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { getSubscription, createInvoice } from '../api/client';

export function SubscriptionPlans() {
  const { t } = useTranslation();

  const { data: sub, isLoading } = useQuery({
    queryKey: ['subscription'],
    queryFn: getSubscription,
  });

  const invoiceMutation = useMutation({
    mutationFn: (plan: string) => createInvoice({ plan }),
    onSuccess: (data) => {
      console.debug('[subscription] invoice created', data);
      window.open(data.invoice_link, '_blank');
    },
  });

  const plans = [
    {
      key: 'free',
      title: t('subscription.free_plan'),
      price: 0,
      features: ['10 сообщений/день', 'Базовые исправления'],
    },
    {
      key: 'weekly',
      title: t('subscription.weekly', { stars: 300 }),
      price: 300,
      features: ['Безлимит сообщений', '7 дней', 'Приоритет AI'],
    },
    {
      key: 'monthly',
      title: t('subscription.monthly', { stars: 800 }),
      price: 800,
      features: ['Безлимит сообщений', '30 дней', 'Приоритет AI', 'Экспорт словаря'],
    },
  ];

  return (
    <div className="scroll-area">
      <div className="section-title">{t('subscription.title')}</div>

      {isLoading && (
        <div className="placeholder">
          <div className="placeholder-icon">⏳</div>
          <div>{t('common.loading')}</div>
        </div>
      )}

      {sub?.active && (
        <div className="card" style={{
          border: '2px solid var(--tg-button)',
          textAlign: 'center',
        }}>
          <div style={{ fontSize: 18, marginBottom: 4 }}>{t('subscription.active')}</div>
          <div style={{ fontSize: 13, color: 'var(--tg-hint)' }}>
            {sub.plan} — {sub.expires_at}
          </div>
        </div>
      )}

      {plans.map(plan => (
        <div
          key={plan.key}
          className="card"
          style={{
            border: plan.key === 'monthly' ? '2px solid var(--tg-button)' : '1px solid var(--tg-border)',
          }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <span style={{ fontSize: 16, fontWeight: 600 }}>{plan.title}</span>
            {plan.key === 'monthly' && (
              <span style={{
                fontSize: 11,
                background: 'var(--tg-button)',
                color: 'var(--tg-button-text)',
                padding: '2px 8px',
                borderRadius: 10,
              }}>
                BEST VALUE
              </span>
            )}
          </div>
          <ul style={{ fontSize: 13, color: 'var(--tg-hint)', listStyle: 'none', marginBottom: 12 }}>
            {plan.features.map((f, i) => (
              <li key={i} style={{ padding: '2px 0' }}>✅ {f}</li>
            ))}
          </ul>
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
      ))}
    </div>
  );
}
