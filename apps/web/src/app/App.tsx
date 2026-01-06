import ReactDOM from 'react-dom/client';
import { StrictMode, useEffect } from 'react';
import './index.css';
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import * as Page from './pages';
import { PopoverImageProvider } from '@/context/PopoverImageContext';
import { SeasonProvider } from '@/context/SeasonContext';
import { ErrorBoundary } from './ErrorBoundary';
import { DefaultLayout, FfmpegLayout, Ltoe, Nav } from '@/components/layout';

function Router() {
    return (
        <Routes>
            <Route element={<DefaultLayout />}>
                <Route path='/' element={<Page.Home />} />
                <Route path='/aspect-ratio-calculator' element={<Page.AspectRatio />} />
                <Route path='/pdf-to-png' element={<Page.PdfToPng />} />
                <Route path='/ico-converter' element={<Page.ImageToIcon />} />
                <Route path='/aim-trainer' element={<Page.AimTrainer />} />
                <Route path='/settings' element={<Page.Settings />} />
            </Route>

            <Route element={<FfmpegLayout hasAudio={true} />}>
                <Route path='/video-to-mp3' element={<Page.VideoToMp3 />} />
            </Route>
            <Route element={<FfmpegLayout hasAudio={false} />}>
                <Route path='/video-cropper' element={<Page.VideoCropper />} />
            </Route>
        </Routes>
    );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <ErrorBoundary>
            <SeasonProvider>
                <PopoverImageProvider>
                    <BrowserRouter>
                        <Nav />
                        <Router />
                        <Ltoe />
                    </BrowserRouter>
                </PopoverImageProvider>
            </SeasonProvider>
        </ErrorBoundary>
    </StrictMode>
);
