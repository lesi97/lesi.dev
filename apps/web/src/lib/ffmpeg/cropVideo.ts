'use client';
import { timeToSeconds, downloadFile } from '@/utils';
import { loadMetadata } from './loadMetaData';
import { fetchFile } from '@ffmpeg/util';

export async function cropVideo(
    ffmpegRef: any,
    videoFile: any,
    videoRef: React.RefObject<HTMLVideoElement>,
    setProgress: (value: number) => void,
    setLoading: (value: boolean) => void
) {
    try {
        setLoading(true);
        const ffmpeg = ffmpegRef.current;

        const videoElement = videoRef.current;
        if (videoElement === null) return;
        const duration = (await loadMetadata(videoElement)) as number;

        if (!ffmpeg.Loaded) {
            await ffmpeg.load();
        }

        ffmpeg.on('log', ({ message }: { message: string }) => {
            const timeMatch = message.match(/time=(\d{2}:\d{2}:\d{2}\.\d{2})/);
            if (timeMatch) {
                const currentTime = timeToSeconds(timeMatch[1]);
                const percentage = (currentTime / duration) * 100;
                setProgress(Math.min(percentage, 100));
            }
        });

        // Write the video file to FFmpeg's virtual filesystem
        ffmpeg.writeFile(videoFile.name, await fetchFile(videoFile));

        // Convert to MP3
        const newName = videoFile.name.replace(/\.(mp4|avi|mov|wmv|mkv)$/i, '_lesi.mp4');
        await ffmpeg.exec([
            '-i',
            videoFile.name,
            '-vf',
            'crop=(ih*9/16):ih:(iw-ih*9/16)/2:0',
            '-c:v',
            'libx264',
            '-profile:v',
            'baseline',
            '-crf',
            '18',
            '-preset',
            'slow',
            '-b:a',
            '192k',
            newName,
        ]);

        // Read the output file and create a URL
        const data = await ffmpeg.readFile(newName);
        const url = URL.createObjectURL(new Blob([data.buffer], { type: 'audio/mpeg' }));
        setLoading(false);
        downloadFile(newName, url);

        // Clean up FFmpeg filesystem
        ffmpeg.deleteFile(videoFile.name);
        ffmpeg.deleteFile(newName);
    } catch (error) {
        console.error('Error converting video to new aspect ratio:', error);
        setLoading(false);
        setProgress(0);
    }
}
