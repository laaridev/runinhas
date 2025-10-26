/**
 * Safe JSON parsing utilities
 */

/**
 * Safely parse JSON with fallback value
 * @param response - JSON string to parse
 * @param fallback - Fallback value if parsing fails
 * @returns Parsed object or fallback
 */
export function safeJsonParse<T>(response: string, fallback: T): T {
  try {
    return JSON.parse(response) as T;
  } catch (error) {
    console.error('Failed to parse JSON:', error);
    return fallback;
  }
}

/**
 * Safely parse JSON with error callback
 * @param response - JSON string to parse
 * @param onError - Callback function if parsing fails
 * @returns Parsed object or null
 */
export function safeJsonParseWithError<T>(
  response: string,
  onError?: (error: Error) => void
): T | null {
  try {
    return JSON.parse(response) as T;
  } catch (error) {
    if (onError && error instanceof Error) {
      onError(error);
    }
    return null;
  }
}

/**
 * Stringify JSON with fallback
 * @param data - Data to stringify
 * @param fallback - Fallback string if stringification fails
 * @returns JSON string or fallback
 */
export function safeJsonStringify(data: any, fallback: string = '{}'): string {
  try {
    return JSON.stringify(data);
  } catch (error) {
    console.error('Failed to stringify JSON:', error);
    return fallback;
  }
}
