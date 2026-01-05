import { useState, useRef, useEffect, Dispatch, SetStateAction, MutableRefObject } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
// import useFfmpeg from './useFfmpeg';
import { convertToMp3, writeToMemory } from '../lib/web-app';
import { useFfmpeg } from '../context/FfmpegContext';

type UseVideoToMp3ReturnType = {
    loading: boolean;
    setLoading: Dispatch<SetStateAction<boolean>>;
    video: any;
    setVideo: Dispatch<SetStateAction<any>>;
    mp3: string;
    setMp3: Dispatch<SetStateAction<string>>;
    progress: number;
    setProgress: Dispatch<SetStateAction<number>>;
    videoRef: MutableRefObject<any>;
    originalVideoBlobUrl: string | undefined;
    setOriginalVideoBlobUrl: Dispatch<SetStateAction<string | undefined>>;
    progressBarTotalRef: MutableRefObject<any>;
    ffmpeg: { ready: boolean; ffmpegRef: MutableRefObject<FFmpeg> };
};

export default function useVideoToMp3(ffmpeg: {
    ready: boolean;
    ffmpegRef: React.MutableRefObject<FFmpeg>;
}): UseVideoToMp3ReturnType {
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
    const [mp3, setMp3] = useState('');
    const [progress, setProgress] = useState(0);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState<string | undefined>();
    const videoRef = useRef(null);
    const progressBarTotalRef = useRef(null);

    useEffect(() => {
        if (!video) return;
        (async () => {
            try {
                const url = await writeToMemory(video, ffmpeg);
                setOriginalVideoBlobUrl(url);
                if (!videoRef.current) return;
                convertToMp3(ffmpeg.ffmpegRef, video, videoRef, setProgress, setLoading, setMp3);
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
