import { config } from 'dotenv';

if (process.env.NODE_ENV !== 'production') {
  config({ path: '../../.env.dev' });
} else {
  config({ path: '../../.env' });
}

export function env(names: string[]) {
  let values: Record<string, string> = {};
  for (const name of names) {
    const value = process.env[name];
    if (!value) {
      throw new Error(`Environment variable ${name} is not set.`);
    }
    values[name] = value;
  }
  return values;
}

