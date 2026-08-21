import { type RefObject } from 'react';
import { loadMetadata } from './loadMetaData';
import { fetchFile } from '@ffmpeg/util';
import { replaceVideoExtensionWithMp3 } from './replaceVideoExtensionWithMp3';
import { downloadFile, timeToSeconds } from '@/utils';

export async function convertToMp3(
    ffmpegRef: any,
    videoFile: any,
    videoRef: RefObject<HTMLVideoElement>,
    setProgress: (value: number) => void,
    setLoading: (value: boolean) => void,
    setMp3: (value: string) => void
) {
    try {
        setLoading(true);
        const ffmpeg = ffmpegRef.current;

        const videoElement = videoRef.current;
        if (videoElement === null) {
            setLoading(false);
            return;
        }
        let duration = 0;
        try {
            duration = (await loadMetadata(videoElement)) as number;
        } catch (error) {
            console.warn('Could not load video metadata; continuing conversion without duration-based progress.', error);
        }

        if (!ffmpeg.Loaded) {
            await ffmpeg.load();
        }

        ffmpeg.on('log', ({ message }: { message: string }) => {
            const timeMatch = message.match(/time=(\d{2}:\d{2}:\d{2}\.\d{2})/);
            if (timeMatch && duration > 0) {
                const currentTime = timeToSeconds(timeMatch[1]);
                const percentage = (currentTime / duration) * 100;
                setProgress(Math.min(percentage, 100));
            }
        });

        // Write the video file to FFmpeg's virtual filesystem
        await ffmpeg.writeFile(videoFile.name, await fetchFile(videoFile));

        // Convert to MP3
        const mp3Name = replaceVideoExtensionWithMp3(videoFile.name);
        await ffmpeg.exec([
            '-i',
            videoFile.name,
            '-vn', // Disable video
            '-ar',
            '44100', // Audio sample rate
            '-ac',
            '2', // Audio channels (stereo)
            '-b:a',
            '192k', // Audio bitrate
            mp3Name,
        ]);

        // Read the output file and create a URL
        const data = await ffmpeg.readFile(mp3Name);
        const url = URL.createObjectURL(new Blob([data.buffer], { type: 'audio/mpeg' }));
        setMp3(url);
        setProgress(100);
        setLoading(false);
        downloadFile(mp3Name, url);

        // Clean up FFmpeg filesystem
        ffmpeg.deleteFile(videoFile.name);
        ffmpeg.deleteFile(mp3Name);
    } catch (error) {
        console.error('Error converting video to MP3:', error);
        setLoading(false);
        setProgress(0);
    }
}
