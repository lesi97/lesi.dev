import { useRef, useEffect, useState, useCallback } from "react";

export const VideoTimeLine = ({ videoSrc, duration, blobUrl }) => {
    const canvasRef = useRef(null);
    const thumbnailVideoRef = useRef(null);
    const [canvasWidth, setCanvasWidth] = useState(0);

    useEffect(() => {
        const updateCanvasDimensions = () => {
            if (canvasRef.current) {
                const containerWidth = canvasRef.current.parentElement.offsetWidth;
                setCanvasWidth(containerWidth);
            }
        };

        updateCanvasDimensions();  // Initialize dimensions
        const resizeObserver = new ResizeObserver(updateCanvasDimensions);
        if (canvasRef.current && canvasRef.current.parentElement) {
            resizeObserver.observe(canvasRef.current.parentElement);
        }

        return () => resizeObserver.disconnect();
    }, []);

    const calculateNumberOfThumbnails = (duration) => {
        if (duration < 30) {
            return Math.ceil(duration / 5); // Fewer thumbnails for very short videos
        } else if (duration < 120) {
            return Math.ceil(duration / 10); // Medium number of thumbnails for short videos
        } else {
            return 12; // Max 12 thumbnails for longer videos
        }
    };

    const drawImage = useCallback((thumbnailWidth, thumbnailHeight) => {
        const canvas = canvasRef.current;
        const context = canvas.getContext('2d');
        let x = 0;
        return () => {
            // Calculate vertical centering
            const dy = (canvas.height - thumbnailHeight);
            context.drawImage(thumbnailVideoRef.current, 0, 0, thumbnailVideoRef.current.videoWidth, 70, x, 0, thumbnailWidth, 70);
            x += thumbnailWidth;
            if (x < canvasWidth) {
                thumbnailVideoRef.current.currentTime = thumbnailWidth * (x / thumbnailWidth);
            }
        };
    }, [canvasWidth]);

    useEffect(() => {
        if (!canvasWidth) return; // Skip if canvas width is not set yet

        const thumbnailVideo = thumbnailVideoRef.current;
        thumbnailVideo.src = blobUrl;

        const maxThumbnails = calculateNumberOfThumbnails(duration) < 10 ? 10 : calculateNumberOfThumbnails(duration);
        console.log(maxThumbnails);
        const thumbnailHeight = 70; // Maximum height set to 70px
        const thumbnailWidth = (thumbnailHeight * 9) / 16; // Width calculated based on 9:16 aspect ratio

        const timestamps = Array.from({ length: maxThumbnails }, (_, i) => (duration / maxThumbnails) * i);

        const handleSeeked = drawImage(thumbnailWidth, thumbnailHeight);
        thumbnailVideo.addEventListener('seeked', handleSeeked);

        thumbnailVideo.addEventListener('loadedmetadata', () => {
            thumbnailVideo.currentTime = timestamps[0]; // Start capturing thumbnails
        });

        return () => {
            thumbnailVideo.removeEventListener('seeked', handleSeeked);
            thumbnailVideo.removeEventListener('loadedmetadata', handleSeeked);
        };

    }, [videoSrc, duration, blobUrl, canvasWidth, drawImage]);

    return (
        <div className="timeLineContainer">
            <canvas className="thumbnails" ref={canvasRef} style={{ width: '100%' }}></canvas>
            <video className="hidden" ref={thumbnailVideoRef} style={{ display: 'none' }}></video>
        </div>
    );
};
