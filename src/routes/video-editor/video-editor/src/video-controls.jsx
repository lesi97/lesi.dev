import "./video-controls.scss";
import { useEffect, useState } from "react";
import { ButtonPlayPause } from "./components/buttonPlayPause/buttonPlayPause";
import { VideoTimeLine } from "./components/videoTimeline/videoTimeLine";

export const VideoControls = ({ video, blobUrl }) => {
    const [timeStamp, setTimeStamp] = useState("00:00.000");
    const [duration, setDuration] = useState("00:00.000");
    const [rawDuration, setRawDuration] = useState(0);

    useEffect(() => {
        if (!video) return;

        function handleMetadataLoaded() {
            const videoDuration = video.duration;
            setRawDuration(videoDuration);
            const durationMinutes = Math.floor(videoDuration / 60);
            const durationSeconds = Math.floor(videoDuration % 60);
            const durationMilliseconds = Math.floor((videoDuration % 1) * 1000);
            const duration = `${pad(durationMinutes, "m")}:${pad(durationSeconds, "s")}.${pad(durationMilliseconds, "ms")}`;
            setDuration(duration);
        }
        video.addEventListener("loadedmetadata", handleMetadataLoaded);
        return () => {
            video.removeEventListener("loadedmetadata", handleMetadataLoaded);
        };
    }, [video]);

    useEffect(() => {
        if (!video) return;
        let animationFrameId = null;
        const updateTimestamp = () => {
            const videoMinutes = Math.floor(video.currentTime / 60);
            const videoSeconds = Math.floor(video.currentTime % 60);
            const videoMilliseconds = Math.floor(
                (video.currentTime % 1) * 1000
            );
            const timeStamp = `${pad(videoMinutes, "m")}:${pad(videoSeconds, "s")}.${pad(videoMilliseconds, "ms")}`;
            setTimeStamp(timeStamp);

            animationFrameId = requestAnimationFrame(updateTimestamp);
        };

        const startUpdate = () => {
            if (animationFrameId === null) {
                animationFrameId = requestAnimationFrame(updateTimestamp);
            }
        };

        const stopUpdate = () => {
            if (animationFrameId !== null) {
                cancelAnimationFrame(animationFrameId);
                animationFrameId = null;
            }
        };

        video.addEventListener("play", startUpdate);
        video.addEventListener("pause", stopUpdate);
        video.addEventListener("ended", stopUpdate);

        return () => {
            video.removeEventListener("play", startUpdate);
            video.removeEventListener("pause", stopUpdate);
            video.removeEventListener("ended", stopUpdate);
            if (animationFrameId !== null) {
                cancelAnimationFrame(animationFrameId);
            }
        };
    }, [video]);

    function pad(value, type) {
        switch (type) {
            case "m":
                return value.toString().padStart(2, "0");
            case "s":
                return value.toString().padStart(2, "0");
            case "ms":
                return value.toString().padStart(3, "0");
            default:
                return value.toString().padStart(2, "0");
        }
    }

    // function segmentVideoDuration() {
    //     const interval = 12;
    //     const segments = [];
    //     console.log(rawDuration)
    //     for (let time = 0; time < rawDuration; time += interval) {
    //         const nextTime = time + interval;
    //         if (nextTime < rawDuration) {
    //             segments.push(nextTime);
    //         } else {
    //             segments.push(rawDuration);
    //             break;
    //         }
    //     }
    //     return segments;
    // }

    return (
        <>
            {video && (
                <>
                    <div className="controlBar">
                        <div className="controls">
                            <label>{timeStamp}</label>
                            <ButtonPlayPause video={video} />
                            <label>{duration}</label>
                        </div>

                        {/* <div className="videoSeekbar">
                        <input type="range" id="seekBar" value="0" min="0" max="100" step="0.1" />
                        <progress value="0" max="100" step="1" id="progressBar" className="progressBar"></progress>
                        <progress value="0" max="100" step="1" id="bufferBar" className="bufferBar"></progress>
                        
                    </div> */}
                        <VideoTimeLine
                            video={video}
                            duration={rawDuration}
                            blobUrl={blobUrl}
                        />
                    </div>
                </>
            )}
        </>
    );
};
