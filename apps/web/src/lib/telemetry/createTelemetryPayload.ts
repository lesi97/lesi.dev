export type TelemetryPayload = {
    timestamp: string;
    route: string;
    userAgent: string;
    ip?: string | null;
    error?: string;
};

export function createTelemetryPayload(route: string, ip?: string | null, error?: string): TelemetryPayload {
    const timestamp = new Date().toISOString();
    const userAgent = typeof navigator === 'undefined' ? 'unknown' : navigator.userAgent;

    if (!error) {
        return { timestamp, route, userAgent, ip };
    }

    return { timestamp, route, userAgent, ip, error };
}
