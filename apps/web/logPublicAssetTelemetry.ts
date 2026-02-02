import { ipAddress } from '@vercel/functions';

type TelemetryPayload = {
    timestamp: string;
    route: string;
    userAgent: string;
    ip: string | null;
};

export async function logPublicAssetTelemetry(request: Request): Promise<void> {
    const url = new URL(request.url);
    const telemetryUrl = `${url.origin}/api/v1/telemetry`;
    const userAgent = request.headers.get('user-agent') ?? 'unknown';
    const ip = ipAddress(request) ?? null;

    const payload: TelemetryPayload = {
        timestamp: new Date().toISOString(),
        route: url.pathname,
        userAgent,
        ip,
    };

    const headers = new Headers();
    headers.set('Content-Type', 'application/json');
    headers.set('Origin', url.origin);
    if (process.env.TELEMETRY_API_KEY) {
        headers.set('X-Telemetry-Key', process.env.TELEMETRY_API_KEY);
    }

    try {
        await fetch(telemetryUrl, {
            method: 'POST',
            headers,
            body: JSON.stringify(payload),
        });
    } catch {
        return;
    }
}
