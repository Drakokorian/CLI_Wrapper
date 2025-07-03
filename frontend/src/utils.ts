export function maskSecrets(text: string): string {
  // Match OpenAI API keys which follow the sk-<48+ alphanumeric chars> format
  // Keys may occasionally vary in length, so match 48 or more characters
  const apiKeyPattern = /(sk-[A-Za-z0-9]{48,})/g;
  return text.replace(apiKeyPattern, (m) => '*'.repeat(m.length));
}
