import ReactDOM from 'react-dom/client';
import { StrictMode } from 'react';
import './index.css';
import { Routes, Route, BrowserRouter } from 'react-router-dom';
import * as Page from './pages';
import { PopoverImageProvider } from '@/context/PopoverImageContext';
import { SeasonProvider } from '@/context/SeasonContext';
import { Ltoe } from '@/components/layout';
import { ErrorBoundary } from './ErrorBoundary';
import { PlainLayout } from '@/components/layout/homeLayout';
import { ContentLayout } from '@/components/layout/contentLayout';

function App() {
    const hexColours = [
        { offset: 0, colour: 'rgba(0,229,255,0.20)' },
        { offset: 0.5, colour: 'rgba(244,63,94,0.20)' },
        { offset: 1, colour: 'rgba(192,38,211, 0.20)' },
    ];
    const gradient = `linear-gradient(to right, ${hexColours.map((obj) => obj.colour).join(', ')})`;

    return (
        <>
            <Routes>
                <Route element={<PlainLayout />}>
                    <Route path='/' element={<Page.Home />} />
                </Route>

                <Route element={<ContentLayout gradient={gradient} stops={hexColours} />}>
                    <Route path='/tools' element={<Page.Tools />} />
                    <Route path='/aim-trainer' element={<Page.AimTrainer />} />
                </Route>
            </Routes>
        </>
    );
}

ReactDOM.createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <ErrorBoundary>
            <SeasonProvider>
                <PopoverImageProvider>
                    <BrowserRouter>
                        <App />
                        <Ltoe />
                    </BrowserRouter>
                </PopoverImageProvider>
            </SeasonProvider>
        </ErrorBoundary>
    </StrictMode>
);
