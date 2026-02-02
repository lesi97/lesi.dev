export async function getPublicIP(): Promise<string | null> {
    if (typeof window === 'undefined') {
        return null;
    }

    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 1500);

    try {
        const response = await fetch('https://api.ipify.org?format=json', { signal: controller.signal });
        if (!response.ok) {
            return null;
        }

        const data = (await response.json()) as { ip?: string };
        if (!data.ip) {
            return null;
        }

        return data.ip;
    } catch {
        return null;
    } finally {
        window.clearTimeout(timeout);
    }
}
