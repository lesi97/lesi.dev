import ReactDOM from 'react-dom/client';
import { StrictMode, lazy, Suspense, useEffect } from 'react';
import './index.css';
import { Routes, Route, BrowserRouter, Navigate, useLocation } from 'react-router-dom';
import { Home } from './pages/home';
import { AspectRatio } from './pages/aspectRatio';
import { ImageToIcon } from './pages/imageToIcon';
import { PasswordGenerator } from './pages/passwordGenerator';
import { UuidGenerator } from './pages/uuidGenerator';
import { WeightConverter } from './pages/weightConverter';
import { Minifier } from './pages/minifier';
import { Countdown } from './pages/countdown';
import { DiscordTimestampConverter } from './pages/discordTimestampConverter';
import { AimTrainer } from './pages/aimTrainer';
import { Settings } from './pages/settings';
import { NotFound } from './pages/notFound';
import { AimTrainerReleaseNotes } from './pages/docs/aimTrainerReleaseNotes';
import { BrowserGpuAcceleration } from './pages/docs/browserGpuAcceleration';
import { loadPdfToPng } from './pages/pdfToPng';
import { loadVideoToMp3 } from './pages/videoToMp3';
import { loadVideoCropper } from './pages/videoCropper';
import { PopoverImageProvider } from '@/context/PopoverImageContext';
import { SeasonProvider } from '@/context/SeasonContext';
import { ErrorBoundary } from './ErrorBoundary';
import { DefaultLayout, FfmpegLayout, WideLayout, Ltoe, Nav } from '@/components/layout';
import { PopoverLayout } from '@/components/layout/popoverLayout';
import { DateMeme, WordpressAdmin } from './pages';
import { sendTelemetry } from '@/lib/telemetry/sendTelemetry';

const PdfToPng = lazy(loadPdfToPng);
const VideoToMp3 = lazy(loadVideoToMp3);
const VideoCropper = lazy(loadVideoCropper);

function Router() {
    const location = useLocation();

    const isMemePage =
        location.pathname.startsWith('/wp-admin') ||
        location.pathname.startsWith('/wp-content') ||
        location.pathname.startsWith('/date') ||
        location.pathname.startsWith('/audrey');

    useEffect(() => {
        sendTelemetry(location.pathname);
    }, [location.pathname]);

    return (
        <>
            {!isMemePage && <Nav />}
            <Suspense fallback={null}>
                <Routes>
                    <Route element={<DefaultLayout />}>
                        <Route path='/' element={<Home />} />
                        <Route path='/aspect-ratio-calculator' element={<AspectRatio />} />
                        <Route path='/pdf-to-png' element={<PdfToPng />} />
                        <Route path='/ico-converter' element={<ImageToIcon />} />
                        <Route path='/password-generator' element={<PasswordGenerator />} />
                        <Route path='/uuid-generator' element={<Navigate to='/uuid-generator/v4' replace />} />
                        <Route path='/uuid-generator/v1' element={<UuidGenerator version='1' />} />
                        <Route path='/uuid-generator/v4' element={<UuidGenerator version='4' />} />
                        <Route path='/uuid-generator/v6' element={<UuidGenerator version='6' />} />
                        <Route path='/uuid-generator/v7' element={<UuidGenerator version='7' />} />
                        <Route path='/uuid-generator/nil' element={<UuidGenerator version='nil' />} />
                        <Route path='/weight-converter' element={<WeightConverter />} />
                        <Route path='/minifier' element={<Minifier />} />
                        <Route path='/countdown' element={<Countdown />} />
                        <Route path='/discord-timestamp-converter' element={<DiscordTimestampConverter />} />
                        <Route path='/settings' element={<Settings />} />
                        <Route path='/aim-trainer/release-notes' element={<AimTrainerReleaseNotes />} />
                    </Route>
                    <Route path='/audrey' element={<DateMeme />} />
                    <Route path='/date/:slug' element={<DateMeme />} />

                    <Route element={<FfmpegLayout hasAudio={true} />}>
                        <Route path='/video-to-mp3' element={<VideoToMp3 />} />
                    </Route>
                    <Route element={<FfmpegLayout hasAudio={false} />}>
                        <Route path='/video-cropper' element={<VideoCropper />} />
                    </Route>

                    <Route element={<WideLayout />}>
                        <Route path='/aim-trainer' element={<AimTrainer />} />
                    </Route>

                    <Route element={<PopoverLayout />}>
                        <Route path='/docs/browser-gpu-acceleration' element={<BrowserGpuAcceleration />} />
                    </Route>

                    <Route element={<WordpressAdmin />}>
                        <Route path='/wp-admin/*' />
                        <Route path='/wp-content/*' />
                    </Route>

                    <Route element={<DefaultLayout />}>
                        <Route path='*' element={<NotFound />} />
                    </Route>
                </Routes>
            </Suspense>
            {!isMemePage && <Ltoe />}
        </>
    );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <SeasonProvider>
            <ErrorBoundary>
                <PopoverImageProvider>
                    <BrowserRouter>
                        <Router />
                    </BrowserRouter>
                </PopoverImageProvider>
            </ErrorBoundary>
        </SeasonProvider>
    </StrictMode>
);
