import { useRef, useEffect, forwardRef, useImperativeHandle } from 'react';
import WaveSurfer from 'wavesurfer.js';
import { useMusic } from '@/context/MusicContext';

export const Waveform = forwardRef(({ url }: { url: string }, ref) => {
    const waveformRef = useRef<WaveSurfer | null>(null);
    const containerRef = useRef<HTMLDivElement | null>(null);

    const { isPlaying, setIsPlaying } = useMusic();

    useEffect(() => {
        if (isPlaying) {
            waveformRef.current?.play();
        } else {
            waveformRef.current?.pause();
        }
    }, [isPlaying]);

    useImperativeHandle(ref, () => ({
        playPause: () => waveformRef.current?.playPause(),
    }));

    useEffect(() => {
        if (!containerRef.current) return;

        let isSeeking = false;

        if (waveformRef.current) {
            try {
                waveformRef.current.destroy();
            } catch (error) {
                console.warn('Previous WaveSurfer instance cleanup failed:', error);
            }
        }

        waveformRef.current = WaveSurfer.create({
            container: containerRef.current,
            waveColor: '#cc3369',
            progressColor: '#6d1652',
            interact: true,
            barWidth: 1,
        });

        waveformRef.current.on('finish', () => setIsPlaying(false));

        try {
            waveformRef.current.load(url);
        } catch (error) {
            console.warn('WaveSurfer load error:', error);
        }

        const handleSpacebar = (e: KeyboardEvent) => {
            if (e.code === 'Space') {
                e.preventDefault();
                setIsPlaying((prev: boolean) => !prev);
            }
        };

        document.addEventListener('keydown', handleSpacebar);

        const seek = (e: MouseEvent | TouchEvent) => {
            if (isSeeking && containerRef.current) {
                const bbox = containerRef.current.getBoundingClientRect();
                const clientX = (e as MouseEvent).clientX || (e as TouchEvent).touches[0].clientX;
                const x = clientX - bbox.left;
                const progress = x / bbox.width;
                waveformRef.current?.seekTo(progress);
            }
        };

        const mouseDown = () => (isSeeking = true);
        const mouseUp = () => (isSeeking = false);

        if (containerRef.current) {
            containerRef.current.addEventListener('mousedown', mouseDown);
            containerRef.current.addEventListener('mouseup', mouseUp);
            containerRef.current.addEventListener('mouseleave', mouseUp);
            containerRef.current.addEventListener('mousemove', seek);
            containerRef.current.addEventListener('touchstart', mouseDown, { passive: true });
            containerRef.current.addEventListener('touchmove', seek, { passive: true });
            containerRef.current.addEventListener('touchend', mouseUp);
        }

        return () => {
            if (waveformRef.current) {
                try {
                    waveformRef.current.destroy();
                } catch (error) {
                    console.warn('WaveSurfer cleanup error:', error);
                }
            }
            document.removeEventListener('keydown', handleSpacebar);

            if (containerRef.current) {
                containerRef.current.removeEventListener('mousedown', mouseDown);
                containerRef.current.removeEventListener('mouseup', mouseUp);
                containerRef.current.removeEventListener('mouseleave', mouseUp);
                containerRef.current.removeEventListener('mousemove', seek);
                containerRef.current.removeEventListener('touchstart', mouseDown);
                containerRef.current.removeEventListener('touchmove', seek);
                containerRef.current.removeEventListener('touchend', mouseUp);
            }
        };
    }, [url]);

    return (
        <div
            id='waveform'
            className='absolute bottom-10 left-0 z-10 w-full cursor-grab px-10 text-red-500'
            ref={containerRef}></div>
    );
});

Waveform.displayName = 'Waveform';

export default Waveform;
