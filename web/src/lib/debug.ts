const isDev = import.meta.env.DEV;

export function debug(...args: unknown[]) {
  if (isDev) {
    console.debug(...args);
  }
}
