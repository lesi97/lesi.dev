import { next, waitUntil } from '@vercel/functions';
import { logPublicAssetTelemetry } from './logPublicAssetTelemetry';

export default function middleware(request: Request) {
    const url = new URL(request.url);
    const pathname = url.pathname;

    if (request.method !== 'GET' && request.method !== 'HEAD') {
        return next();
    }

    if (pathname.startsWith('/api/')) {
        return next();
    }

    if (!pathname.includes('.')) {
        return next();
    }

    if (pathname.startsWith('/assets/')) {
        return next();
    }

    if (pathname.startsWith('/@')) {
        return next();
    }

    waitUntil(logPublicAssetTelemetry(request));

    return next();
}

export const config = {
    matcher: '/:path*',
};
