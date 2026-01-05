import { mergeClassNames } from '@/utils';
import { useEffect, useMemo, useRef, useState, CSSProperties, ReactNode } from 'react';

type ColourStop = { offset: number; colour: string };

type HexBackgroundProps = {
    children: ReactNode;
    className?: string;
    gradient?: string;
    density?: number;
    minDurationMs?: number;
    maxDurationMs?: number;
    stops?: ColourStop[];
};

type Pulse = { idx: number; start: number; dur: number };

const CELL_W = 28;
const CELL_H = 49;

function clamp01(v: number): number {
    if (v < 0) {
        return 0;
    }
    if (v > 1) {
        return 1;
    }
    return v;
}

function lerp(a: number, b: number, t: number): number {
    return a + (b - a) * t;
}

function easeInOut(t: number): number {
    if (t < 0.5) {
        return 2 * t * t;
    }
    return 1 - Math.pow(-2 * t + 2, 2) / 2;
}

function drawHexPath(ctx: CanvasRenderingContext2D, left: number, top: number): void {
    ctx.beginPath();
    ctx.moveTo(left + 13.99, top + 9.25);
    ctx.lineTo(left + 26.99, top + 16.75);
    ctx.lineTo(left + 26.99, top + 31.75);
    ctx.lineTo(left + 13.99, top + 39.25);
    ctx.lineTo(left + 1, top + 31.75);
    ctx.lineTo(left + 1, top + 16.75);
    ctx.closePath();
}

export function HexBackground({
    children,
    className,
    gradient = 'linear-gradient(to right, rgba(245,158,11,.30), rgba(244,63,94,.30), rgba(192,38,211,.30))',
    density = 0.12,
    minDurationMs = 1200,
    maxDurationMs = 2400,
    stops,
}: HexBackgroundProps) {
    const canvasRef = useRef<HTMLCanvasElement | null>(null);
    const hostRef = useRef<HTMLDivElement | null>(null);
    const rafRef = useRef<number | null>(null);
    const [box, setBox] = useState<{ w: number; h: number }>({ w: 0, h: 0 });

    const cols = useMemo(() => (box.w > 0 ? Math.ceil(box.w / CELL_W) + 1 : 0), [box.w]);
    const rows = useMemo(() => (box.h > 0 ? Math.ceil(box.h / CELL_H) + 1 : 0), [box.h]);
    const cellCount = useMemo(() => cols * rows, [cols, rows]);

    const defaultStops = useMemo<ColourStop[]>(() => {
        return [
            { offset: 0, colour: 'rgba(245,158,11,0.30)' },
            { offset: 0.5, colour: 'rgba(244,63,94,0.30)' },
            { offset: 1, colour: 'rgba(192,38,211,0.30)' },
        ];
    }, []);

    const styleVars: CSSProperties & { ['--hex-gradient']?: string } = useMemo(() => {
        return { ['--hex-gradient']: gradient };
    }, [gradient]);

    useEffect(function measure() {
        const el = hostRef.current;
        if (!el) {
            return;
        }
        const ro = new ResizeObserver((entries) => {
            const r = entries[0].contentRect;
            setBox({ w: Math.ceil(r.width), h: Math.ceil(r.height) });
        });
        ro.observe(el);
        return function cleanup() {
            ro.disconnect();
        };
    }, []);

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas || box.w === 0 || box.h === 0) {
            return;
        }
        const dpr = window.devicePixelRatio || 1;
        canvas.width = Math.max(1, Math.floor(box.w * dpr));
        canvas.height = Math.max(1, Math.floor(box.h * dpr));
        canvas.style.width = `${box.w}px`;
        canvas.style.height = `${box.h}px`;
        const ctx = canvas.getContext('2d');
        if (!ctx) {
            return;
        }
        ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

        const g = ctx.createLinearGradient(0, 0, box.w, 0);
        const useStops = stops && stops.length > 0 ? stops : defaultStops;
        for (const s of useStops) {
            g.addColorStop(clamp01(s.offset), s.colour);
        }

        const pulses: Pulse[] = [];
        let last = performance.now();
        let spawnAcc = 0;
        const avgDur = (minDurationMs + maxDurationMs) * 0.5;
        const rate = (Math.max(0, density) * cellCount) / Math.max(300, avgDur);

        function pickIdx(): number {
            return Math.floor(Math.random() * cellCount);
        }

        function draw(now: number): void {
            if (!ctx) {
                return;
            }
            const dt = Math.max(0, now - last);
            last = now;
            ctx.clearRect(0, 0, box.w, box.h);

            spawnAcc += rate * dt;
            while (spawnAcc >= 1) {
                const idx = pickIdx();
                const dur = lerp(minDurationMs, maxDurationMs, Math.random());
                pulses.push({ idx, start: now, dur });
                spawnAcc -= 1;
            }

            for (let i = pulses.length - 1; i >= 0; i -= 1) {
                const p = pulses[i];
                const t = (now - p.start) / p.dur;
                if (t >= 1) {
                    pulses.splice(i, 1);
                    continue;
                }
                const tt = easeInOut(t < 0.5 ? t * 2 : (1 - t) * 2);
                const c = p.idx % cols;
                const r = Math.floor(p.idx / cols);
                const left = c * CELL_W;
                const top = r * CELL_H;
                ctx.save();
                ctx.globalAlpha = tt;
                drawHexPath(ctx, left, top);
                ctx.clip();
                ctx.fillStyle = g;
                ctx.fillRect(0, 0, box.w, box.h);
                ctx.restore();
            }

            rafRef.current = requestAnimationFrame(draw);
        }

        rafRef.current = requestAnimationFrame(draw);
        return function cleanup() {
            if (rafRef.current) {
                cancelAnimationFrame(rafRef.current);
            }
        };
    }, [box.w, box.h, cols, cellCount, defaultStops, stops, density, minDurationMs, maxDurationMs]);

    return (
        <div ref={hostRef} className={mergeClassNames('hex-gradient w-full h-full', className)} style={styleVars}>
            <canvas ref={canvasRef} className='absolute inset-0 z-0 pointer-events-none' />
            <div
                className='relative z-10 w-full h-full'
                style={{ position: 'relative', zIndex: 2, width: '100%', height: '100%' }}>
                {children}
            </div>
        </div>
    );
}
