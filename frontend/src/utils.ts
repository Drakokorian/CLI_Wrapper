export function maskSecrets(text: string): string {
  // Match OpenAI API keys which follow the "sk-" prefix and consist of
  // at least 48 alphanumeric characters. Word boundaries ensure that the
  // entire key is replaced rather than just a substring.
  const apiKeyPattern = /\bsk-[A-Za-z0-9]{48,}\b/g;
  return text.replace(apiKeyPattern, (m) => '*'.repeat(m.length));
}
