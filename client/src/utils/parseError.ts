import { ZodError } from 'zod/v4';

export function parseError(error: ZodError): Record<string, string>;
export function parseError(error: unknown): string;

/**
 * Parses an error and returns the error as a string
 *
 * @export
 * @param {unknown} error
 * @returns {string}
 *
 * @example
 * try {
 *  throw new Error('Test');
 * } catch (error) {
 *  const parsedError = parseError(error);
 *  console.log(parsedError); // Test
 * }
 */
export function parseError(error: unknown): string | Record<string, string> {
    if (typeof error === 'string') {
        return error;
    }

    if (error instanceof ZodError) {
        const parsed = JSON.parse(error.message);
        const errors = parsed.reduce((acc: Record<string, unknown>, err: Record<string, string[]>) => {
            const key = err.path[0];
            const val = err.message;
            acc[key] = val;
            return acc;
        }, {});
        return errors;
    }

    if (error instanceof Error) {
        return error.message;
    }

    return 'An error has occurred';
}
