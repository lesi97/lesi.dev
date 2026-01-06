import { useState, useRef, useEffect, Dispatch, SetStateAction, MutableRefObject, RefObject } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
import { convertToMp3 } from '@/lib/ffmpeg/convertToMp3';
import { writeToMemory } from '@/lib/ffmpeg/writeToMemory';

type UseVideoToMp3ReturnType = {
    loading: boolean;
    setLoading: Dispatch<SetStateAction<boolean>>;
    video: any;
    setVideo: Dispatch<SetStateAction<any>>;
    mp3: string;
    setMp3: Dispatch<SetStateAction<string>>;
    progress: number;
    setProgress: Dispatch<SetStateAction<number>>;
    videoRef: RefObject<any>;
    originalVideoBlobUrl: string | undefined;
    setOriginalVideoBlobUrl: Dispatch<SetStateAction<string | undefined>>;
    progressBarTotalRef: RefObject<any>;
    ffmpeg: { ready: boolean; ffmpegRef: RefObject<FFmpeg> };
};

export function useVideoToMp3(ffmpeg: { ready: boolean; ffmpegRef: RefObject<FFmpeg> }): UseVideoToMp3ReturnType {
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
    const [mp3, setMp3] = useState('');
    const [progress, setProgress] = useState(0);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState<string | undefined>();
    const videoRef = useRef<HTMLVideoElement | null>(null);
    const progressBarTotalRef = useRef(null);

    useEffect(() => {
        if (!video) return;
        (async () => {
            try {
                const url = await writeToMemory(video, ffmpeg);
                setOriginalVideoBlobUrl(url);
                if (!videoRef.current) {
                    return;
                }
                convertToMp3(
                    ffmpeg.ffmpegRef,
                    video,
                    videoRef as RefObject<HTMLVideoElement>,
                    setProgress,
                    setLoading,
                    setMp3
                );
            } catch (error) {
                console.error(error);
            }
        })();
    }, [video]);

    useEffect(() => {
        if (videoRef.current) (videoRef.current as HTMLVideoElement).load();
    }, [originalVideoBlobUrl]);

    return {
        loading,
        setLoading,
        video,
        setVideo,
        mp3,
        setMp3,
        progress,
        setProgress,
        videoRef,
        originalVideoBlobUrl,
        setOriginalVideoBlobUrl,
        progressBarTotalRef,
        ffmpeg,
    };
}
