import React, { useEffect, useRef } from 'react';
import { gsap } from 'gsap';
import { SplitText } from 'gsap/SplitText';
import { ScrambleTextPlugin } from 'gsap/ScrambleTextPlugin';
import { mergeClassNames } from '@/utils';

gsap.registerPlugin(SplitText, ScrambleTextPlugin);

export interface ScrambledTextProps {
    radius?: number;
    duration?: number;
    speed?: number;
    scrambleChars?: string;
    className?: string;
    style?: React.CSSProperties;
    children?: React.ReactNode;
    active?: boolean;
}

export function ScrambledText({
    radius = 100,
    duration = 1.2,
    speed = 0.5,
    scrambleChars = '.:',
    className = '',
    style = {},
    children,
    active = true,
}: ScrambledTextProps) {
    const rootRef = useRef<HTMLSpanElement | null>(null);
    const textElRef = useRef<HTMLElement | null>(null);
    const splitRef = useRef<SplitText | null>(null);
    const moveHandlerRef = useRef<(e: PointerEvent) => void>(() => null);

    useEffect(
        function init() {
            if (!active) {
                return;
            }
            if (!rootRef.current) {
                return;
            }
            textElRef.current = rootRef.current.querySelector('[data-role="text"]') as HTMLElement | null;
            if (!textElRef.current) {
                return;
            }
            if (splitRef.current) {
                splitRef.current.revert();
                splitRef.current = null;
            }
            splitRef.current = SplitText.create(textElRef.current, {
                type: 'chars',
                charsClass: 'inline-block will-change-transform',
            });
            splitRef.current.chars.forEach((el) => {
                const c = el as HTMLElement;
                gsap.set(c, { attr: { 'data-content': c.textContent || '' } });
            });
            const handleMove = (e: PointerEvent) => {
                if (!splitRef.current) {
                    return;
                }
                splitRef.current.chars.forEach((el) => {
                    const c = el as HTMLElement;
                    const r = c.getBoundingClientRect();
                    const dx = e.clientX - (r.left + r.width / 2);
                    const dy = e.clientY - (r.top + r.height / 2);
                    const dist = Math.hypot(dx, dy);
                    if (dist < radius) {
                        gsap.to(c, {
                            overwrite: true,
                            duration: duration * (1 - dist / radius),
                            scrambleText: { text: c.dataset.content || '', chars: scrambleChars, speed },
                            ease: 'none',
                        });
                    }
                });
            };
            moveHandlerRef.current = handleMove;
            textElRef.current.addEventListener('pointermove', handleMove);
            return function cleanup() {
                if (textElRef.current && moveHandlerRef.current) {
                    textElRef.current.removeEventListener('pointermove', moveHandlerRef.current);
                }
                if (splitRef.current) {
                    splitRef.current.revert();
                    splitRef.current = null;
                }
            };
        },
        [active, radius, duration, speed, scrambleChars]
    );

    useEffect(
        function syncChildren() {
            if (!active) {
                return;
            }
            if (!textElRef.current) {
                return;
            }
            const next = String(children ?? '');
            if (splitRef.current && splitRef.current.chars.length === next.length) {
                const spans = splitRef.current.chars as HTMLElement[];
                for (let i = 0; i < spans.length; i += 1) {
                    const c = spans[i];
                    const ch = next[i] ?? '';
                    c.dataset.content = ch;
                    if (c.textContent !== ch) {
                        gsap.to(c, {
                            overwrite: 'auto',
                            duration: 0.18,
                            scrambleText: { text: ch, chars: scrambleChars, speed },
                            ease: 'none',
                        });
                    }
                }
                return;
            }
            if (splitRef.current) {
                splitRef.current.revert();
                splitRef.current = null;
            }
            textElRef.current.textContent = next;
            splitRef.current = SplitText.create(textElRef.current, {
                type: 'chars',
                charsClass: 'inline-block will-change-transform',
            });
            splitRef.current.chars.forEach((el) => {
                const c = el as HTMLElement;
                gsap.set(c, { attr: { 'data-content': c.textContent || '' } });
            });
        },
        [children, active, scrambleChars]
    );

    return (
        <span ref={rootRef} className={mergeClassNames('inline-block text-primary', className)} style={style}>
            <span data-role='text'>{children}</span>
        </span>
    );
}
