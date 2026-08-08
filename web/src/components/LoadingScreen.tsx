import { LoadingDots } from './LoadingDots';

interface LoadingScreenProps {
  text?: string;
}

export function LoadingScreen({ text }: LoadingScreenProps) {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '56px 24px',
      color: 'var(--tg-hint)',
      textAlign: 'center',
      gap: 16,
    }}>
      <LoadingDots size={12} />
      {text && (
        <div style={{ fontSize: 15, fontWeight: 500 }}>{text}</div>
      )}
    </div>
  );
}
