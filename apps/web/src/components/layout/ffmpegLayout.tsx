import { FfmpegProvider } from '@/context/FfmpegContext';
import { MusicProvider } from '@/context/MusicContext';
import { Outlet } from 'react-router-dom';
import { Footer } from './footer';

export function FfmpegLayout({ hasAudio = false }: { hasAudio?: boolean }) {
    if (hasAudio) {
        return (
            <FfmpegProvider>
                <MusicProvider>
                    <MainContent />
                </MusicProvider>
            </FfmpegProvider>
        );
    }

    return (
        <FfmpegProvider>
            <MainContent />
        </FfmpegProvider>
    );
}

function MainContent() {
    return (
        <>
            <main className='relative top-8 mb-8 flex h-fit w-11/12 justify-center rounded-lg bg-base-100 px-8 py-8 shadow 2xl:w-50% 2xl:min-w-50%'>
                <div className='flex w-11/12 flex-col gap-4'>
                    <Outlet />
                </div>
            </main>
            <Footer />
        </>
    );
}
