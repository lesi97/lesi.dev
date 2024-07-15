import "./video-to-mp3.scss";
import { DropBox } from "../../../components/dropbox/dropbox";
// import { LoaderNewtonsCradle, LoaderRing } from "../../../components/loading/";
import { createFFmpeg, fetchFile } from "@ffmpeg/ffmpeg";
import { useEffect, useState, useRef } from "react";
import { writeToMemory } from "../video-editor/src/write-to-memory";

const ffmpeg = createFFmpeg({ log: true });

export const VideoToMp3 = ({ setError }) => {
    const [ready, setReady] = useState(false);
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
    const [mp3, setMp3] = useState(null);
    const [progress, setProgress] = useState(0);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState();
    // const progressBarRef = useRef(null);
    const progressBarTotalRef = useRef(null);

    setError(false);

    useEffect(() => {
        document.title = `Lesi | Video To MP3`;
    }, []);

    async function load() {
        if (!ffmpeg.isLoaded()) {
            await ffmpeg.load();
        }
        setReady(true);
    }

    useEffect(() => {
        load();
    }, []);

    useEffect(() => {
        if (!video) return;
        convertToMp3();
    }, [video])

    const timeToSeconds = (time) => {
        const parts = time.split(':');
        return parseFloat(parts[0]) * 3600 + parseFloat(parts[1]) * 60 + parseFloat(parts[2]);
    };

    const loadMetadata = (vidPlayer) => {
        return new Promise((resolve, reject) => {
            vidPlayer.onloadedmetadata = () => {
                console.log('Metadata loaded!');
                resolve(vidPlayer.duration);
            };
            vidPlayer.onerror = (e) => {
                console.error('Failed to load metadata:', e);
                reject(e);
            };
        });
    };

    async function convertToMp3() {
        const label = document.getElementById("uploadText");
        label.innerHTML = `Converting file, please wait a moment<br><br><span class="wrongFileType">${video.name}</span>`;
        setLoading(true);

        const vidPlayer = document.getElementById('video');
        const duration = await loadMetadata(vidPlayer);

        if (!ffmpeg.isLoaded()) {
            await ffmpeg.load()
        }

        ffmpeg.setLogger(({ message }) => {
            const timeMatch = message.match(/time=(\d{2}:\d{2}:\d{2}\.\d{2})/);
            if (timeMatch) {
                const currentTime = timeToSeconds(timeMatch[1]);
                const percentage = (currentTime / duration) * 100;
                setProgress(percentage);
                return label.innerHTML = `Converting video, please wait a moment<br><br><span class="wrongFileType">${video.name}</span><br>${percentage.toFixed(2)}% Complete`;
            }
        });


        ffmpeg.FS("writeFile", video.name, await fetchFile(video));
        const mp3Name = video.name.replace(/\.(mp4|avi|mov|wmv|mkv)$/i, ".mp3");
        await ffmpeg.run("-i", video.name, "-vn", "-ar", "44100", "-ac", "2", "-b:a", "192k", mp3Name);

        setLoading(false)
        label.innerHTML = `<span class="wrongFileType">${video.name}</span><br>
            has now finished converting, check your Downloads folder for this file<br><br>
            You can also drag and drop a new file to convert`;

        const data = ffmpeg.FS("readFile", mp3Name);
        const url = URL.createObjectURL(new Blob([data.buffer], { type: "audio/mpeg" }));
        setMp3(url);

        const downloadLink = document.createElement("a");
        downloadLink.download = mp3Name;
        downloadLink.href = url;
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
    }

    async function loadVideo() {
        const fileInput = document.getElementById("fileInput");
        const file = fileInput?.files?.[0];
        const label = document.getElementById("uploadText");
        if (!file.type.startsWith("video/")) {
            label.innerHTML = `
                <span class="wrongFileType">File type not valid!</span>
                <span>
                    Please use a valid <span class="wrongFileType">video</span> file<br /><br />
                    Or click here to browse<br /> 
                    your PC for a <span class="wrongFileType">video</span> file to upload
                </span>
            `;
            label.classList.add("invalid-animation");
            setTimeout(() => {
                label.classList.remove("invalid-animation");
            }, 450);
            return;
        }

        label.innerHTML = file.name;
        const url = await writeToMemory(file, ffmpeg);
        setOriginalVideoBlobUrl(url);
        setVideo(file);
    }

    return (
        <main>
            <div className="videoToMp3">
                {/* {loading && <LoaderRing />} */}
                <div className="description">
                    <h1>Video To MP3 Converter</h1>
                    <h2>
                        Drag and drop a video file to convert it to MP3 and download it<br />
                        Press the play button once loaded to play the audio in your browser
                    </h2>
                </div>
                {ready && <DropBox loadVideo={loadVideo} mp3url={mp3} type="videoToMp3" loading={loading} progressBar={progress} progressBarTotal={progressBarTotalRef} />}
                {/*eslint-disable-next-line jsx-a11y/media-has-caption*/}
                {video && <video id="video" src={originalVideoBlobUrl} />}
            </div>
        </main>
    );
}