import "./video-editor.scss";
import { DropBox } from "../../../components/dropbox/dropbox";
// import { LoaderNewtonsCradle, LoaderRing } from "../../../components/loading";
import { createFFmpeg } from "@ffmpeg/ffmpeg";
import { useEffect, useState } from "react";
// import { Navigate, useSearchParams, useNavigate } from "react-router-dom";
import { writeToMemory } from "./src/write-to-memory";
import { VideoControls } from "./src/video-controls";

const ffmpeg = createFFmpeg({ log: true });

export const VideoEditor = ({ setError }) => {
    const [ready, setReady] = useState(false);
    // const [loading, setLoading] = useState(false);
    const [originalVideoBlobUrl, setOriginalVideoBlobUrl] = useState();
    const [video, setVideo] = useState();

    setError(false);


    useEffect(() => {
        document.title = `Lesi | Video Editor`;
        load();
    }, []);

    async function load() {
        if (!ffmpeg.isLoaded()) {
            await ffmpeg.load();
        }
        setReady(true);
    }

    useEffect(() => {
        if (!originalVideoBlobUrl) return;
        setVideo(document.getElementById("video"));
    }, [originalVideoBlobUrl]);

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

        const url = await writeToMemory(file, ffmpeg);
        setOriginalVideoBlobUrl(url);
    }


    return (
        <div className="videoEditor">
            {!originalVideoBlobUrl &&
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
                        {ready && <DropBox loadVideo={loadVideo} type="videoToMp3" /*loading={loading}*/ />}
                    </div>
                </main>}

            {originalVideoBlobUrl &&
                <>
                    <video id="video" muted src={originalVideoBlobUrl} controls />
                    <VideoControls video={video} blobUrl={originalVideoBlobUrl} />
                </>
            }
        </div>
    );
}
