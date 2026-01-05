import { useState, useRef, useEffect, Dispatch, SetStateAction, MutableRefObject } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
import { cropVideo, writeToMemory } from '../lib/web-app';
import { useFfmpeg } from '../context/FfmpegContext';

type UseVideoCropperReturnType = {
    loading: boolean;
    setLoading: Dispatch<SetStateAction<boolean>>;
    video: any;
    setVideo: Dispatch<SetStateAction<any>>;
    progress: number;
    setProgress: Dispatch<SetStateAction<number>>;
    videoRef: MutableRefObject<any>;
    originalVideoBlobUrl: string | undefined;
    setOriginalVideoBlobUrl: Dispatch<SetStateAction<string | undefined>>;
    progressBarTotalRef: MutableRefObject<any>;
    ffmpeg: { ready: boolean; ffmpegRef: MutableRefObject<FFmpeg> };
};

export default function useVideoCropper(ffmpeg: {
    ready: boolean;
    ffmpegRef: React.MutableRefObject<FFmpeg>;
}): UseVideoCropperReturnType {
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
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
                cropVideo(ffmpeg.ffmpegRef, video, videoRef, setProgress, setLoading);
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
        progress,
        setProgress,
        videoRef,
        originalVideoBlobUrl,
        setOriginalVideoBlobUrl,
        progressBarTotalRef,
        ffmpeg,
    };
}
