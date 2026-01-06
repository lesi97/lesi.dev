import { fetchFile } from '@ffmpeg/util';

export async function writeToMemory(video: any, ffmpeg: any) {
    await ffmpeg.ffmpegRef.current.writeFile(video.name, await fetchFile(video));
    const data = await ffmpeg.ffmpegRef.current.readFile(video.name);
    return URL.createObjectURL(new Blob([data.buffer], { type: video.type }));
}
