import { useState, useRef, useEffect, Dispatch, SetStateAction, MutableRefObject, RefObject } from 'react';
import { FFmpeg } from '@ffmpeg/ffmpeg';
import { writeToMemory } from '@/lib/ffmpeg/writeToMemory';
import { cropVideo } from '@/lib/ffmpeg/cropVideo';

type UseVideoCropperReturnType = {
    loading: boolean;
    setLoading: Dispatch<SetStateAction<boolean>>;
    video: any;
    setVideo: Dispatch<SetStateAction<any>>;
    progress: number;
    setProgress: Dispatch<SetStateAction<number>>;
    videoRef: RefObject<any>;
    originalVideoBlobUrl: string | undefined;
    setOriginalVideoBlobUrl: Dispatch<SetStateAction<string | undefined>>;
    progressBarTotalRef: RefObject<any>;
    ffmpeg: { ready: boolean; ffmpegRef: RefObject<FFmpeg> };
};

export function useVideoCropper(ffmpeg: { ready: boolean; ffmpegRef: RefObject<FFmpeg> }): UseVideoCropperReturnType {
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
    const [progress, setProgress] = useState(0);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState<string | undefined>();
    const videoRef = useRef<HTMLVideoElement>(null);
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
                cropVideo(ffmpeg.ffmpegRef, video, videoRef as RefObject<HTMLVideoElement>, setProgress, setLoading);
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
