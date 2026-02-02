import { createTelemetryPayload } from './createTelemetryPayload';
import { getPublicIP } from './getPublicIP';

export async function sendTelemetry(route: string, error?: string): Promise<void> {
    const ip = await getPublicIP();
    const payload = createTelemetryPayload(route, ip, error);
    const telemetryKey = import.meta.env.VITE_TELEMETRY_API_KEY as string | undefined;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (telemetryKey) {
        headers['X-Telemetry-Key'] = telemetryKey;
    }

    try {
        await fetch('/api/v1/telemetry', {
            method: 'POST',
            headers,
            body: JSON.stringify(payload),
            keepalive: true,
        });
    } catch {
        return;
    }
}
