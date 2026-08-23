export type DateMemeClickAction = 'yes' | 'no';

type DateMemeClickOptions = {
    secretEnding?: boolean;
};

export async function sendDateMemeClick(
    route: string,
    action: DateMemeClickAction,
    options: DateMemeClickOptions = {}
): Promise<void> {
    const payload: { route: string; action: DateMemeClickAction; secretEnding?: boolean } = { route, action };
    const telemetryKey = import.meta.env.VITE_TELEMETRY_API_KEY as string | undefined;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (telemetryKey) {
        headers['X-Telemetry-Key'] = telemetryKey;
    }
    if (options.secretEnding) {
        payload.secretEnding = true;
    }

    try {
        await fetch('/api/v1/telemetry/date-meme-click', {
            method: 'POST',
            headers,
            body: JSON.stringify(payload),
            keepalive: true,
        });
    } catch {
        return;
    }
}
