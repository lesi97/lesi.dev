export type DateMemeClickAction = 'yes' | 'no';

export async function sendDateMemeClick(route: string, action: DateMemeClickAction): Promise<void> {
    const telemetryKey = import.meta.env.VITE_TELEMETRY_API_KEY as string | undefined;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (telemetryKey) {
        headers['X-Telemetry-Key'] = telemetryKey;
    }

    try {
        await fetch('/api/v1/telemetry/date-meme-click', {
            method: 'POST',
            headers,
            body: JSON.stringify({ route, action }),
            keepalive: true,
        });
    } catch {
        return;
    }
}
