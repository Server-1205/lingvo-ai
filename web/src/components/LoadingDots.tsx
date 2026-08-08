const dotStyle: React.CSSProperties = {
  width: 10,
  height: 10,
  borderRadius: '50%',
  background: 'var(--c-primary)',
  display: 'inline-block',
  animation: 'lds-bounce 1.4s ease-in-out infinite both',
};

const containerStyle: React.CSSProperties = {
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '4px 0',
};

const delays = ['0s', '0.16s', '0.32s'];

export function LoadingDots({ size = 10 }: { size?: number }) {
  return (
    <div style={containerStyle}>
      {[0, 1, 2].map(i => (
        <span
          key={i}
          style={{
            ...dotStyle,
            width: size,
            height: size,
            animationDelay: delays[i],
          }}
        />
      ))}
    </div>
  );
}
