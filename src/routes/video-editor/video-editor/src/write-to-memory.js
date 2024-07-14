import { fetchFile } from "@ffmpeg/ffmpeg";

export async function writeToMemory(video, ffmpeg) {

    ffmpeg.FS("writeFile", video.name, await fetchFile(video));
    const data = ffmpeg.FS("readFile", video.name);
    return URL.createObjectURL(new Blob([data.buffer], { type: video.type }));
}