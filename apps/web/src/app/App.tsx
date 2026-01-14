import ReactDOM from 'react-dom/client';
import { StrictMode, useEffect } from 'react';
import './index.css';
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import * as Page from './pages';
import * as Docs from './pages/docs';
import { PopoverImageProvider } from '@/context/PopoverImageContext';
import { SeasonProvider } from '@/context/SeasonContext';
import { ErrorBoundary } from './ErrorBoundary';
import { DefaultLayout, FfmpegLayout, WideLayout, Ltoe, Nav } from '@/components/layout';
import { PopoverLayout } from '@/components/layout/popoverLayout';

function Router() {
    return (
        <Routes>
            <Route element={<DefaultLayout />}>
                <Route path='/' element={<Page.Home />} />
                <Route path='/aspect-ratio-calculator' element={<Page.AspectRatio />} />
                <Route path='/pdf-to-png' element={<Page.PdfToPng />} />
                <Route path='/ico-converter' element={<Page.ImageToIcon />} />
                <Route path='/password-generator' element={<Page.PasswordGenerator />} />
                <Route path='/weight-converter' element={<Page.WeightConverter />} />
                <Route path='/minifier' element={<Page.Minifier />} />
                <Route path='/countdown' element={<Page.Countdown />} />
                <Route path='/settings' element={<Page.Settings />} />
                <Route path='/aim-trainer/release-notes' element={<Docs.AimTrainerReleaseNotes />} />
            </Route>

            <Route element={<FfmpegLayout hasAudio={true} />}>
                <Route path='/video-to-mp3' element={<Page.VideoToMp3 />} />
            </Route>
            <Route element={<FfmpegLayout hasAudio={false} />}>
                <Route path='/video-cropper' element={<Page.VideoCropper />} />
            </Route>

            <Route element={<WideLayout />}>
                <Route path='/aim-trainer' element={<Page.AimTrainer />} />
            </Route>

            <Route element={<PopoverLayout />}>
                <Route path='/docs/browser-gpu-acceleration' element={<Docs.BrowserGpuAcceleration />} />
            </Route>

            <Route element={<DefaultLayout />}>
                <Route path='*' element={<Page.NotFound />} />
            </Route>
        </Routes>
    );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <SeasonProvider>
            <ErrorBoundary>
                <PopoverImageProvider>
                    <BrowserRouter>
                        <Nav />
                        <Router />
                        <Ltoe />
                    </BrowserRouter>
                </PopoverImageProvider>
            </ErrorBoundary>
        </SeasonProvider>
    </StrictMode>
);
