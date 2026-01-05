import { useMemo, useState, useRef, useEffect } from 'react';

const FADE_MS = 600;
const GROW_MS = 800;
const VISIBLE_MS = 1200;
const FAST_EASE = 'cubic-bezier(0.2, 0.9, 0.2, 1)';

export type SinKey = 'envy' | 'gluttony' | 'greed' | 'lust' | 'pride' | 'sloth' | 'wrath';

type SinItem = {
    key: SinKey;
    label: string;
    imageSrc: string;
    imageAlt: string;
    audioSrc: string;
    volume: number;
};

type Phase = 'hidden' | 'pre-in' | 'animating-in' | 'visible' | 'animating-out';

export function triggerSin(key: SinKey) {
    const ev = new CustomEvent('sin:show', { detail: { key } });
    window.dispatchEvent(ev);
}

export function SinOverlay() {
    const items: SinItem[] = useMemo(
        () => [
            {
                key: 'envy',
                label: 'Envy',
                imageSrc: '/_static/images/sins/envy.jpg',
                imageAlt: 'Cat_Envy',
                audioSrc: '/_static/audio/sins/envy.mp3',
                volume: 0.5,
            },
            {
                key: 'gluttony',
                label: 'Gluttony',
                imageSrc: '/_static/images/sins/gluttony.webp',
                imageAlt: 'Cat_Gluttony',
                audioSrc: '/_static/audio/sins/gluttony.mp3',
                volume: 0.7,
            },
            {
                key: 'greed',
                label: 'Greed',
                imageSrc: '/_static/images/sins/greed.jpg',
                imageAlt: 'Cat_Greed',
                audioSrc: '/_static/audio/sins/greed.mp3',
                volume: 0.5,
            },
            {
                key: 'lust',
                label: 'Lust',
                imageSrc: '/_static/images/sins/lust.jpg',
                imageAlt: 'Cat_Lust',
                audioSrc: '/_static/audio/sins/lust.mp3',
                volume: 0.2,
            },
            {
                key: 'pride',
                label: 'Pride',
                imageSrc: '/_static/images/sins/pride.jpeg',
                imageAlt: 'Cat_Pride',
                audioSrc: '/_static/audio/sins/pride.mp3',
                volume: 0.2,
            },
            {
                key: 'sloth',
                label: 'Sloth',
                imageSrc: '/_static/images/sins/sloth.webp',
                imageAlt: 'Cat_Sloth',
                audioSrc: '/_static/audio/sins/sloth.mp3',
                volume: 1,
            },
            {
                key: 'wrath',
                label: 'Wrath',
                imageSrc: '/_static/images/sins/wrath.jpg',
                imageAlt: 'Cat_Wrath',
                audioSrc: '/_static/audio/sins/wrath.mp3',
                volume: 0.2,
            },
        ],
        []
    );

    const [active, setActive] = useState<SinItem | null>(null);
    const [phase, setPhase] = useState<Phase>('hidden');
    const audioRef = useRef<HTMLAudioElement | null>(null);
    const t = useRef<number[]>([]);

    function clamp01(x: number) {
        if (x < 0) {
            return 0;
        }
        if (x > 1) {
            return 1;
        }
        return x;
    }

    function clearTimers() {
        t.current.forEach((id) => window.clearTimeout(id));
        t.current = [];
    }

    function stopAudio() {
        if (audioRef.current) {
            audioRef.current.pause();
            audioRef.current.currentTime = 0;
        }
    }

    function show(key: SinKey) {
        const item = items.find((i) => i.key === key);
        if (!item) {
            return;
        }
        clearTimers();
        stopAudio();
        setActive(item);
        setPhase('pre-in');
        requestAnimationFrame(() => {
            setPhase('animating-in');
            const t1 = window.setTimeout(() => {
                setPhase('visible');
                const t2 = window.setTimeout(() => {
                    setPhase('animating-out');
                    const t3 = window.setTimeout(() => {
                        setPhase('hidden');
                        setActive(null);
                        stopAudio();
                    }, FADE_MS);
                    t.current.push(t3);
                }, VISIBLE_MS);
                t.current.push(t2);
            }, GROW_MS);
            t.current.push(t1);
        });
    }

    useEffect(() => {
        function onShow(e: Event) {
            const detail = (e as CustomEvent<{ key: SinKey }>).detail;
            if (!detail) {
                return;
            }
            show(detail.key);
        }
        window.addEventListener('sin:show', onShow);
        return () => {
            window.removeEventListener('sin:show', onShow);
            clearTimers();
            stopAudio();
        };
    }, [items]);

    useEffect(() => {
        if (!active) {
            return;
        }
        const a = new Audio(active.audioSrc);
        a.volume = clamp01(active.volume);
        audioRef.current = a;
        void a.play();
        return () => {
            a.pause();
            a.currentTime = 0;
        };
    }, [active]);

    const visible = phase !== 'hidden';

    const transform =
        phase === 'pre-in'
            ? 'scale(0.2)'
            : phase === 'animating-in'
              ? 'scale(1)'
              : phase === 'visible'
                ? 'scale(1)'
                : phase === 'animating-out'
                  ? 'scale(1.05)'
                  : 'scale(0.2)';

    const opacity =
        phase === 'pre-in'
            ? 0
            : phase === 'animating-in'
              ? 1
              : phase === 'visible'
                ? 1
                : phase === 'animating-out'
                  ? 0
                  : 0;

    const duration = phase === 'pre-in' ? '0ms' : phase === 'animating-in' ? `${GROW_MS}ms` : `${FADE_MS}ms`;

    const timing = phase === 'animating-in' ? FAST_EASE : 'ease-out';

    return (
        <div
            className={`fixed inset-0 z-[999] ${visible ? 'pointer-events-none' : 'pointer-events-none'}`}
            aria-hidden={!visible}>
            {active && (
                <img
                    key={active.key}
                    src={active.imageSrc}
                    alt={active.imageAlt}
                    className='w-auto h-screen object-cover mx-auto'
                    style={{
                        transform,
                        opacity,
                        transitionProperty: 'transform, opacity',
                        transitionDuration: duration,
                        transitionTimingFunction: timing,
                        willChange: 'transform, opacity',
                    }}
                />
            )}
        </div>
    );
}
