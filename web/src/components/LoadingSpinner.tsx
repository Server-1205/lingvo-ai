const spinnerStyle: React.CSSProperties = {
  width: 36,
  height: 36,
  border: '3px solid var(--c-outline-variant)',
  borderTopColor: 'var(--c-primary)',
  borderRadius: '50%',
  animation: 'spin 0.8s linear infinite',
};

interface LoadingSpinnerProps {
  size?: number;
  stroke?: number;
}

export function LoadingSpinner({ size = 36, stroke = 3 }: LoadingSpinnerProps) {
  return (
    <span
      style={{
        ...spinnerStyle,
        width: size,
        height: size,
        borderWidth: stroke,
      }}
    />
  );
}
