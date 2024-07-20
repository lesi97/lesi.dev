
import { useRef, useEffect, forwardRef, useImperativeHandle } from "react";
import WaveSurfer from "wavesurfer.js";
import "./waveform.scss"

export const Waveform = forwardRef(({ url, setPlaying }, ref) => {
    const waveformRef = useRef(null);
    const wavesurfer = useRef(null);

    useImperativeHandle(ref, () => ({
        playPause() {
            wavesurfer.current.playPause();
        },
    }));

    useEffect(() => {
        let isSeeking = false;
        wavesurfer.current = WaveSurfer.create({
            container: waveformRef.current,
            waveColor: "#cc3369",
            progressColor: "#6d1652",
            // cursorColor: 'navy',
            interact: false,
            scrollParent: true,
            barWidth: 1,
        });

        wavesurfer.current.on("finish", () => {
            setPlaying(false);
        });

        wavesurfer.current.load(url);
        document.addEventListener('keydown', handleSpacebar);

        const seek = (e) => {
            if (isSeeking) {
                const bbox = waveformRef.current.getBoundingClientRect();
                const clientX = e.clientX || e.touches[0].clientX;
                const x = e.clientX - bbox.left;
                const progress = x / bbox.width;
                wavesurfer.current.seekTo(progress);
            }
        };

        const mouseDown = () => {
            isSeeking = true;
        };

        const mouseUp = () => {
            isSeeking = false;
        };

        const touchStart = (e) => {
            e.preventDefault();
            mouseDown();
        };

        const touchMove = (e) => {
            e.preventDefault();
            seek(e);
        };

        const touchEnd = () => {
            mouseUp();
        };

        waveformRef.current.addEventListener('mousedown', mouseDown);
        waveformRef.current.addEventListener('mouseup', mouseUp);
        waveformRef.current.addEventListener('mouseleave', mouseUp);
        waveformRef.current.addEventListener('mousemove', seek);

        waveformRef.current.addEventListener('touchstart', touchStart);
        waveformRef.current.addEventListener('touchmove', touchMove);
        waveformRef.current.addEventListener('touchend', touchEnd);



        return () => {
            wavesurfer.current.destroy();
            waveformRef.current?.removeEventListener('mousedown', mouseDown);
            waveformRef.current?.removeEventListener('mouseup', mouseUp);
            waveformRef.current?.removeEventListener('mouseleave', mouseUp);
            waveformRef.current?.removeEventListener('mousemove', seek);

            waveformRef.current?.removeEventListener('touchstart', touchStart);
            waveformRef.current?.removeEventListener('touchmove', touchMove);
            waveformRef.current?.removeEventListener('touchend', touchEnd);

            document?.removeEventListener('keydown', handleSpacebar);
        };

    }, [url]);

    const handleSpacebar = (e) => {
        if (e.code === "Space") {
            e.preventDefault();
            wavesurfer.current.playPause();
            setPlaying(wavesurfer.current.isPlaying());
        }
    };



    return <div id="waveform" className="waveformMp3" ref={waveformRef}></div>;
})

Waveform.displayName = "Waveform";