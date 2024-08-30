import "./crop-to-shorts.scss";
import { DropBox } from "../../../components/dropbox/dropbox";
// import { LoaderNewtonsCradle, LoaderRing } from "../../../components/loading/";
import { createFFmpeg, fetchFile } from "@ffmpeg/ffmpeg";
import { useEffect, useState, useRef } from "react";
import { writeToMemory } from "../video-editor/src/write-to-memory";
import { useLocation } from "react-router-dom";

const ffmpeg = createFFmpeg({ log: true });

export const CropToShort = ({ setError }) => {

    const [ready, setReady] = useState(false);
    const [loading, setLoading] = useState(false);
    const [video, setVideo] = useState();
    // const [croppedVid, setCroppedVid] = useState();
    const [progress, setProgress] = useState(0);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState();
    const progressBarRef = useRef(null);
    const progressBarTotalRef = useRef(null);
    const location = useLocation();

    useEffect(() => {
        document.title = `Lesi | Video Cropper`;
        setError(false);
        load();
    }, []);

    async function load() {
        if (!ffmpeg.isLoaded()) {
            await ffmpeg.load();
        }
        setReady(true);
    }

    useEffect(() => {
        return () => {
            if (originalVideoBlobUrl) {
                URL.revokeObjectURL(originalVideoBlobUrl);
            }
        };
    }, [location, originalVideoBlobUrl]);

    useEffect(() => {
        if (!video) return;
        cropVideo();
    }, [video])

    useEffect(() => {
        const totalLength = 350;
        const progressLength = (progress / 100) * totalLength;
        const startX = 128.83144;
        const endX = startX + progressLength;
        // const progressPath = `M${endX},356.80298v12H128.83144c-3.31,0-6-2.69-6-6s2.69-6,6-6h231.5Z`;
        const progressPath = `M478.83144,356.80298H128.83144c-3.31,0-6,2.69-6,6s2.69,6,6,6H${endX}c3.31,0,6-2.69,6-6s-2.69-6-6-6Z`;


        if (progressBarRef.current) {
            progressBarRef.current.setAttribute("d", progressPath);
            if (progress === 0) {
                progressBarRef.current.setAttribute("class", "hidden");
            } else {
                progressBarRef.current.removeAttribute("class");
            }
        }
        // console.log(progressBarTotalRef);
        if (progressBarTotalRef.current && progress !== 0) {
            progressBarTotalRef.current.removeAttribute("class");
        }

    }, [progress]);



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

    async function cropVideo() {
        const label = document.getElementById("uploadText");
        label.innerHTML = `Converting file, please wait a moment<br><br><span class="wrongFileType">${video.name}</span><br>${progress.toFixed(2)}% Complete`;
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
                return label.innerHTML = `Cropping video, please wait a moment<br><br><span class="wrongFileType">${video.name}</span><br>${percentage.toFixed(2)}% Complete`;
            }
        });



        ffmpeg.FS("writeFile", video.name, await fetchFile(video));
        const newName = video.name.replace(/\.(mp4|avi|mov|wmv|mkv)$/i, '-lesi-edited.$1');
        await ffmpeg.run(
            "-i", video.name,
            "-vf", "crop=(ih*9/16):ih:(iw-ih*9/16)/2:0",
            "-c:v", "libx264",
            "-profile:v", "baseline",
            "-crf", "18",
            "-preset", "slow",
            "-b:a", "192k",
            newName);

        setLoading(false)
        label.innerHTML = `<span class="wrongFileType">${video.name}</span><br>
            has now finished converting, check your Downloads folder for this file<br><br>
            You can also drag and drop a new file to convert`;
        ffmpeg.FS('readdir', '/')
        const data = ffmpeg.FS("readFile", newName);
        const url = URL.createObjectURL(new Blob([data.buffer], { type: "video/mp4" }));

        const downloadLink = document.createElement("a");
        downloadLink.download = newName;
        downloadLink.href = url;
        document.body.appendChild(downloadLink);
        downloadLink.click();
        document.body.removeChild(downloadLink);
    }

    async function loadVideo() {
        const fileInput = document.getElementById("fileInput");
        console.log(fileInput);
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
        setTimeout(setVideo(file), 1);
    }

    return (
        <main>
            <div className="videoToMp3">
                {/* {loading && <LoaderRing />} */}
                <div className="description">
                    <h1>Video Cropper</h1>
                    <h2>
                        Drag and drop a video file to crop it to a 9:16 aspect ratio<br />&nbsp;
                    </h2>
                </div>
                {ready && <DropBox loadVideo={loadVideo} type="videoCrop" loading={loading} progressBar={progressBarRef} progressBarTotal={progressBarTotalRef} />}
                {video && <video id="video" muted src={originalVideoBlobUrl} />}
            </div>
        </main>
    );
}