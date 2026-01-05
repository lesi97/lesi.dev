import { createContext, useContext, useState, ReactNode, useEffect, useRef } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
import { toBlobURL } from '@ffmpeg/util';

interface FfmpegContextType {
    ready: boolean;
    ffmpegRef: React.MutableRefObject<FFmpeg>;
}

export const FfmpegContext = createContext<FfmpegContextType | undefined>(undefined);

export function FfmpegProvider({ children }: { children: ReactNode }) {
    const ffmpegRef = useRef(new FFmpeg());
    const [ready, setReady] = useState(false);

    useEffect(() => {
        load();
    }, []);

    const load = async () => {
        const baseURL = 'https://unpkg.com/@ffmpeg/core@0.12.6/dist/umd';
        const ffmpeg = ffmpegRef.current;
        ffmpeg.on('log', ({ message }) => {
            console.log(message);
        });
        // toBlobURL is used to bypass CORS issue, urls with the same domain can be used directly.
        // await ffmpeg.load({
        //     coreURL: await toBlobURL('/ffmpeg/ffmpeg-core.js', ' text/javascript'),
        //     wasmURL: await toBlobURL('/ffmpeg/ffmpeg-core.wasm', 'application/wasm'),
        // });
        await ffmpeg.load({
            coreURL: await toBlobURL(`${baseURL}/ffmpeg-core.js`, 'text/javascript'),
            wasmURL: await toBlobURL(`${baseURL}/ffmpeg-core.wasm`, 'application/wasm'),
        });
        setReady(true);
    };

    return <FfmpegContext.Provider value={{ ready, ffmpegRef }}>{children}</FfmpegContext.Provider>;
}

export function useFfmpeg() {
    const context = useContext(FfmpegContext);
    if (context === undefined) {
        throw new Error('useFfmpeg must be used within a FFmpeg Provider');
    }
    return context;
}
