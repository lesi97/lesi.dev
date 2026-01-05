import { mergeClassNames } from '@/utils';
import { useId, useState, useEffect } from 'react';
import { SinKey, triggerSin } from './sinOverlay';

type TextArcProps = {
    text: string;
    radius: number;
    angleDeg: number;
    arcDegrees?: number;
    durationMs?: number;
    strokeWidth?: number;
    className?: string;
    wrapperClassName?: string;
    id?: string;
    viewBoxPad?: number;
    overlayKey: SinKey;
};

export function ArcedText() {
    const arcs = [
        {
            label: 'CIRCLE GLUTTONY',
            overlayKey: 'gluttony' as SinKey,
            radius: 310,
            angleDeg: -1,
            wrapperClassName:
                'absolute -left-[275px] -top-[10px] w-0 h-0 flex items-start justify-start place-content-start',
        },
        {
            label: 'CIRCLE LUST',
            overlayKey: 'lust' as SinKey,
            radius: 310,
            angleDeg: 1,
            wrapperClassName:
                'absolute -left-[280px] top-[35px] w-0 h-0 flex items-start justify-start place-content-start',
        },
        {
            label: 'CIRCLE PRIDE',
            overlayKey: 'pride' as SinKey,
            radius: 310,
            angleDeg: 0,
            wrapperClassName:
                'absolute -left-[280px] top-[70px] w-0 h-0 overflow-visible flex place-content-start items-start justify-start',
        },
        {
            label: 'CIRCLE SLOTH',
            overlayKey: 'sloth' as SinKey,
            radius: 310,
            angleDeg: 2,
            wrapperClassName: 'absolute -left-[289px] top-[115px] w-0 h-0 overflow-visible',
        },
        {
            label: 'CIRCLE WRATH',
            overlayKey: 'wrath' as SinKey,
            radius: 305,
            angleDeg: 2,
            wrapperClassName: 'absolute -left-[288px] top-[155px] w-0 h-0 overflow-visible',
        },
        {
            label: 'CIRCLE ENVY',
            overlayKey: 'envy' as SinKey,
            radius: 300,
            angleDeg: 3,
            wrapperClassName: 'absolute -left-[290px] top-[197px] w-0 h-0 overflow-visible',
        },
        {
            label: 'CIRCLE GREED',
            overlayKey: 'greed' as SinKey,
            radius: 290,
            angleDeg: 4,
            wrapperClassName:
                'flex items-start justify-start -left-[290px] top-[230px] absolute w-0 h-0 overflow-visible',
        },
    ];
    return (
        <div className='absolute left-1/2 -translate-x-[10px] top-0 z-[19] w-60 h-[400px] text-white'>
            <div className='absolute inset-0 bg-base-100 [clip-path:polygon(10%_0%,100%_0%,30%_100%,0%_100%)] pointer-events-auto z-[30]'>
                {arcs.map((arc) => (
                    <TextArc
                        key={arc.label}
                        text={arc.label}
                        overlayKey={arc.overlayKey}
                        radius={arc.radius}
                        angleDeg={arc.angleDeg}
                        wrapperClassName={arc.wrapperClassName}
                    />
                ))}
            </div>
        </div>
    );
}

function TextArc({
    text,
    radius,
    angleDeg,
    arcDegrees = 90,
    durationMs = 600,
    strokeWidth = 8,
    className = '',
    wrapperClassName = 'inline-block relative',
    id,
    viewBoxPad,
    overlayKey,
}: TextArcProps) {
    const autoId = useId();
    const uid = id ?? autoId;
    const [armed, setArmed] = useState<boolean>(false);

    useEffect(() => {
        setArmed(true);
    }, []);

    function polar(cx: number, cy: number, r: number, deg: number) {
        const rad = ((deg - 90) * Math.PI) / 180;
        return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) };
    }

    const size = radius * 2;
    const pad = viewBoxPad ?? Math.ceil(strokeWidth * 2 + 12);
    const cx = radius;
    const cy = radius;
    const startDeg = 0;
    const endDeg = startDeg + arcDegrees;
    const s = polar(cx, cy, radius, startDeg);
    const e = polar(cx, cy, radius, endDeg);
    const largeArc = arcDegrees > 180 ? 1 : 0;
    const d = `M ${s.x} ${s.y} A ${radius} ${radius} 0 ${largeArc} 1 ${e.x} ${e.y}`;

    return (
        <span className={mergeClassNames(wrapperClassName, 'absolute')} style={{ width: size, height: size }}>
            <svg
                role='button'
                className='block w-full h-full shrink-0 hover:cursor-pointer'
                viewBox={`${-pad} ${-pad} ${size + 2 * pad} ${size + 2 * pad}`}
                preserveAspectRatio='xMidYMid meet'
                aria-hidden='true'
                style={{ overflow: 'visible', pointerEvents: 'none' }}>
                <defs>
                    <path id={`arc-${uid}`} d={d} />
                </defs>

                <g
                    style={{
                        transformBox: 'fill-box',
                        transformOrigin: 'center',
                        transform: `rotate(${armed ? angleDeg : 0}deg)`,
                        transition: `transform ${durationMs}ms ease-out`,
                    }}>
                    <path
                        d={d}
                        fill='none'
                        stroke='transparent'
                        strokeWidth={strokeWidth * 2}
                        style={{ pointerEvents: 'stroke' }}
                        onClick={() => triggerSin(overlayKey)}
                    />
                    <text
                        textAnchor='start'
                        dominantBaseline='middle'
                        style={{ fontSize: 14, fontWeight: 600, paintOrder: 'stroke fill', pointerEvents: 'none' }}
                        className={className}>
                        <textPath
                            href={`#arc-${uid}`}
                            stroke='black'
                            strokeWidth={strokeWidth}
                            strokeLinejoin='round'
                            strokeLinecap='round'
                            fill='none'>
                            {text}
                        </textPath>
                    </text>
                    <text
                        textAnchor='start'
                        dominantBaseline='middle'
                        style={{ fontSize: 14, fontWeight: 600, fill: '#00F07A', pointerEvents: 'none' }}
                        className={className}>
                        <textPath href={`#arc-${uid}`}>{text}</textPath>
                    </text>
                </g>
            </svg>
        </span>
    );
}
