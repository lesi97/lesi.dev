import './globals.css';
import { Herr_Von_Muellerhoff, Inter } from 'next/font/google';
import SnowfallWrapper from '../_components/snowfall-wrapper';
import Nav from '@/app/_ui/nav';
import getSeason from '@/app/_lib/getSeason';
const inter = Inter({ subsets: ['latin'] });

export const metadata = {
    title: 'Home | Lesi',
    description: 'Home',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
    const currentSeason: string | null = getSeason();
    return (
        <html
            lang="en"
            className="dark"
            //data-theme="forest"
            data-theme="lesi-dark"
        >
            <body className={`${inter.className} bg-base-100`}>
                <Nav season={currentSeason} />
                {currentSeason ? <SnowfallWrapper season={currentSeason} /> : null}
                <main className="relative top-8 mb-8 min-w-50% w-50% bg-primary py-8 px-12 rounded-lg outline outline-neutral outline-3 left-1/2 -translate-x-1/2 h-fit">
                    <div className="flex flex-col gap-4">{children}</div>
                </main>
            </body>
        </html>
    );
}
