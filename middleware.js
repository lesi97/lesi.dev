import { NextResponse } from 'next/server';

export default function middleware(req) {
    const hostname = req.headers.get('host');
    const url = req.nextUrl.clone();

    if (hostname.startsWith('api-v1.')) {
        if (url.pathname.startsWith('/_next')) {
            //return NextResponse.next();
        }
        if (!url.pathname.startsWith('/api')) {
            url.pathname = `/api/v1${url.pathname}`;
            //return NextResponse.rewrite(url);
        }
    } else {
        if (url.pathname.startsWith('/api')) {
            url.pathname = '/404';
            // return NextResponse.rewrite(url);
        }
    }

    return NextResponse.next();
}
