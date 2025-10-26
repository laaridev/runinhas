export interface RetryOptions {
  maxRetries?: number;
  delayMs?: number;
  backoff?: boolean;
  onRetry?: (attempt: number, error: any) => void;
}

/**
 * Executa uma função com retry automático em caso de falha
 */
export async function withRetry<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> {
  const {
    maxRetries = 3,
    delayMs = 1000,
    backoff = true,
    onRetry
  } = options;

  let lastError: any;

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error;

      // Se foi a última tentativa, lança o erro
      if (attempt === maxRetries) {
        throw error;
      }

      // Callback de retry (opcional)
      if (onRetry) {
        onRetry(attempt + 1, error);
      }

      // Calcula delay com backoff exponencial se habilitado
      const delay = backoff ? delayMs * Math.pow(2, attempt) : delayMs;
      
      // Aguarda antes de tentar novamente
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }

  throw lastError;
}
