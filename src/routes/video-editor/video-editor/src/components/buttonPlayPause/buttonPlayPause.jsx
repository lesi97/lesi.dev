import "./buttonPlayPause.scss";
import { useState } from "react";

export const ButtonPlayPause = ({ video }) => {
    const [paused, setPaused] = useState(video?.paused ?? true);

    function pauseUnpause() {
        paused ? video.play() : video.pause();
        setPaused(!paused);
    }

    return (
        <div
            className="buttonPlayPauseContainer"
            onClick={pauseUnpause}
            role="button"
            tabIndex="0"
            onKeyDown={(e) => {
                if (e.code === "Enter") {
                    pauseUnpause();
                }
            }}
        >
            <div
                className={`buttonPlayPauseLine ${paused ? "playing" : "paused"}`}
            ></div>
        </div>
    );
};
