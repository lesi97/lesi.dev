import { createContext, useContext, useState, type ReactNode, type RefObject } from 'react';
import { useFfmpeg } from './FfmpegContext';
import { FFmpeg } from '@ffmpeg/ffmpeg';

interface MusicContextType {
    isPlaying: boolean;
    setIsPlaying: (isPlaying: boolean | ((prev: boolean) => boolean)) => void;
    ffmpeg: { ready: boolean; ffmpegRef: RefObject<FFmpeg> };
}

export const MusicContext = createContext<MusicContextType | undefined>(undefined);

export function MusicProvider({ children }: { children: ReactNode }) {
    const [isPlaying, setIsPlaying] = useState(false);
    const ffmpeg = useFfmpeg();

    return <MusicContext.Provider value={{ isPlaying, setIsPlaying, ffmpeg }}>{children}</MusicContext.Provider>;
}

export function useMusic() {
    const context = useContext(MusicContext);
    if (context === undefined) {
        throw new Error('useMusic must be used within a MusicProvider');
    }
    return context;
}
