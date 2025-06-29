export function maskSecrets(text: string): string {
  const apiKeyPattern = /(sk-[a-zA-Z0-9]{10,})/g;
  return text.replace(apiKeyPattern, (m) => '*'.repeat(m.length));
}
