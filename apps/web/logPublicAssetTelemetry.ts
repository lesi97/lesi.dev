import { ipAddress } from '@vercel/functions';

type TelemetryPayload = {
    timestamp: string;
    route: string;
    userAgent: string;
    ip: string | null;
};

type Options = {
    sampleRate?: number;
    timeoutMs?: number;
};

export async function logPublicAssetTelemetry(request: Request, options: Options = {}): Promise<void> {
    const url = new URL(request.url);

    if (url.pathname.startsWith('/api/v1/telemetry')) {
        return;
    }

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
    headers.set('X-Telemetry-Source', 'server');
    if (process.env.TELEMETRY_API_KEY) {
        headers.set('X-Telemetry-Key', process.env.TELEMETRY_API_KEY);
    }

    const controller = new AbortController();
    const timeoutMs = options.timeoutMs ?? 800;
    const timeoutId = setTimeout(() => {
        controller.abort();
    }, timeoutMs);

    try {
        await fetch(telemetryUrl, {
            method: 'POST',
            headers,
            body: JSON.stringify(payload),
            signal: controller.signal,
            keepalive: true,
        });
    } catch {
        return;
    } finally {
        clearTimeout(timeoutId);
    }
}
