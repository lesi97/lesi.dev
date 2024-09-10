import { useRef, useEffect, useState, useCallback } from "react";

export const VideoTimeLine = ({ videoSrc, duration, blobUrl }) => {
    const canvasRef = useRef(null);
    const thumbnailVideoRef = useRef(null);

    function drawImage() {
        const canvas = canvasRef.current;
        const context = canvas.getContext("2d");

        let left = 0;
    }

    return (
        <div className="timeLineContainer">
            <canvas
                className="thumbnails"
                ref={canvasRef}
                style={{ width: "100%" }}
            ></canvas>
            <video
                className="hidden"
                ref={thumbnailVideoRef}
                style={{ display: "none" }}
            ></video>
        </div>
    );
};
