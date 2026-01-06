import { createContext, useContext, useState, ReactNode, useEffect, useRef, RefObject } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
import { toBlobURL } from '@ffmpeg/util';

interface FfmpegContextType {
    ready: boolean;
    ffmpegRef: RefObject<FFmpeg>;
}

export const FfmpegContext = createContext<FfmpegContextType | undefined>(undefined);

export function FfmpegProvider({ children }: { children: ReactNode }) {
    const ffmpegRef = useRef(new FFmpeg());
    const [ready, setReady] = useState(false);

    useEffect(() => {
        loadFfmpeg();
    }, []);

    async function loadFfmpeg() {
        try {
            const baseURL = 'https://unpkg.com/@ffmpeg/core@0.12.6/dist/esm';
            const ffmpeg = ffmpegRef.current;
            ffmpeg.on('log', ({ message }) => {
                console.log(message);
            });

            const [coreURL, wasmURL] = await Promise.all([
                toBlobURL(`${baseURL}/ffmpeg-core.js`, 'text/javascript'),
                toBlobURL(`${baseURL}/ffmpeg-core.wasm`, 'application/wasm'),
            ]);

            await ffmpeg.load({ coreURL, wasmURL });
            setReady(true);
        } catch (err) {
            console.error('ffmpeg loading error', err);
        }
    }

    return <FfmpegContext.Provider value={{ ready, ffmpegRef }}>{children}</FfmpegContext.Provider>;
}

export function useFfmpeg() {
    const context = useContext(FfmpegContext);
    if (context === undefined) {
        throw new Error('useFfmpeg must be used within a FFmpeg Provider');
    }
    return context;
}
