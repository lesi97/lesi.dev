import { useEffect, useRef, useState } from 'react';
import { usePageMeta } from '@/hooks';
import { useSeason } from '@/context/SeasonContext';

type MemePrompt = {
    question: string;
    image?: {
        src: string;
        alt: string;
    };
};

type ColourMode = 'light' | 'dark';

const memeImages = {
    angryCat: { src: '/images/date/cat-angry-1.png', alt: 'Angry cat' },
    punchingCatOne: { src: '/images/date/cat-punching-1.jpg', alt: 'Cat ready to punch' },
    punchingCatTwo: { src: '/images/date/cat-punching-2.gif', alt: 'Cat punching' },
    roseCatOne: { src: '/images/date/cat-rose-1.jpg', alt: 'Cat with a rose' },
    roseCatTwo: { src: '/images/date/cat-rose-2.jpg', alt: 'Cat holding a rose' },
    roseCatThree: { src: '/images/date/cat-rose-3.gif', alt: 'Cat offering a rose' },
    rizzFace: { src: '/images/date/rizz-face.webp', alt: 'Rizz face' },
} satisfies Record<string, MemePrompt['image']>;

const acceptedPrompt: MemePrompt = {
    question: 'I knew you would say yes',
    image: memeImages.rizzFace,
};

const initialPrompt: MemePrompt = {
    question: 'Audrey, will you go out on a date with me?',
    image: memeImages.roseCatOne,
};

const prompts: MemePrompt[] = [
    { question: 'Are you sure?', image: memeImages.angryCat },
    { question: 'Really sure? 🤔', image: memeImages.punchingCatOne },
    { question: 'Are you completely sure?', image: memeImages.punchingCatTwo },
    { question: 'What if there is dessert? 🍰', image: memeImages.roseCatTwo },
    { question: 'What if I ask very nicely?', image: memeImages.roseCatThree },
    { question: 'This no button seems suspicious.', image: memeImages.angryCat },
    { question: 'Maybe just a tiny date?', image: memeImages.roseCatOne },
    { question: 'Still thinking no?', image: memeImages.punchingCatTwo },
    { question: 'Final answer?', image: memeImages.punchingCatOne },
    { question: 'The correct button is right there 😔', image: memeImages.roseCatThree },
];

const noLabels = ['No', 'Nope', 'Still no?', 'Try again'];

const noPositions = [
    { x: 82, y: 16, rotate: -5 },
    { x: 18, y: 24, rotate: 7 },
    { x: 88, y: 44, rotate: 5 },
    { x: 12, y: 54, rotate: -9 },
    { x: 72, y: 82, rotate: 8 },
    { x: 28, y: 84, rotate: -6 },
    { x: 50, y: 14, rotate: 4 },
    { x: 88, y: 72, rotate: -8 },
    { x: 14, y: 78, rotate: 9 },
    { x: 66, y: 34, rotate: -4 },
    { x: 34, y: 38, rotate: 6 },
    { x: 50, y: 90, rotate: -3 },
    { x: 8, y: 14, rotate: 10 },
    { x: 92, y: 90, rotate: -10 },
    { x: 50, y: 50, rotate: 3 },
    { x: 24, y: 12, rotate: -7 },
];

const heartRainImages = ['/images/heart-1.webp', '/images/heart-2.webp', '/images/heart-3.webp'];

const heartRainDrops = Array.from({ length: 34 }, (_, index) => ({
    image: heartRainImages[index % heartRainImages.length],
    left: (index * 29) % 100,
    size: 20 + ((index * 11) % 30),
    duration: 5 + ((index * 7) % 5),
    delay: ((index * 0.17) % 2.8).toFixed(2),
    sway: ((index % 2 === 0 ? 1 : -1) * (18 + ((index * 5) % 42))).toString(),
    rotate: ((index * 37) % 48) - 24,
}));

const darkThemes = new Set(['business', 'coffee', 'dim', 'dracula', 'halloween', 'lesi-default', 'sunset']);

const palettes = {
    light: {
        background: 'radial-gradient(circle at top, #ff9acb 0%, #ffc4df 36%, #ffe1ef 68%, #fff4fa 100%)',
        question: '#25151d',
        muted: '#6d4a5c',
        yes: '#ff4fa3',
        yesHover: '#f2338f',
        yesText: '#ffffff',
        no: '#ffffff',
        noText: '#38202c',
        noBorder: '#35202b',
        buttonShadow: '0 0px 34px rgba(255, 79, 163, 0.26)',
        imageShadow: '0 22px 56px rgba(87, 27, 56, 0.22)',
    },
    dark: {
        background: 'radial-gradient(circle at top, #321627 0%, #130a10 46%, #080406 100%)',
        question: '#fff0f7',
        muted: '#f7b8d5',
        yes: '#ff5dac',
        yesHover: '#ff7cbd',
        yesText: '#21000f',
        no: '#2a1a23',
        noText: '#ffe8f3',
        noBorder: '#ffb0d2',
        buttonShadow: '0 0px 38px rgba(255, 93, 172, 0.1)',
        imageShadow: '0 24px 64px rgba(0, 0, 0, 0.48)',
    },
} as const;

const sharedButtonClassName =
    'h-[60px] w-[clamp(120px,31vw,150px)] whitespace-nowrap rounded-full px-6 text-lg font-bold transition-[background-color,transform,left,top] duration-150 ease-out focus-visible:outline focus-visible:outline-4 focus-visible:outline-offset-4';

function getColourMode(): ColourMode {
    if (typeof window === 'undefined') {
        return 'light';
    }

    const currentTheme = document.documentElement.getAttribute('data-theme');

    try {
        const storedTheme = JSON.parse(localStorage.getItem('theme') ?? 'null') as { theme?: string } | null;
        const theme = storedTheme?.theme ?? currentTheme;
        if (theme) {
            return darkThemes.has(theme) ? 'dark' : 'light';
        }
    } catch {
        if (currentTheme) {
            return darkThemes.has(currentTheme) ? 'dark' : 'light';
        }
    }

    return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function ValentineHeartRain() {
    return (
        <div className='pointer-events-none fixed inset-0 z-20 overflow-hidden' aria-hidden='true'>
            <style>
                {`
                    @keyframes date-heart-rain {
                        0% {
                            opacity: 0;
                            transform: translate3d(0, -14vh, 0) rotate(var(--heart-rotate));
                        }
                        10% {
                            opacity: 0.9;
                        }
                        100% {
                            opacity: 0.75;
                            transform: translate3d(var(--heart-sway), 114vh, 0) rotate(calc(var(--heart-rotate) + 260deg));
                        }
                    }
                `}
            </style>
            {heartRainDrops.map((heart, index) => (
                <img
                    key={`${heart.image}-${index}`}
                    src={heart.image}
                    alt=''
                    className='absolute top-0 select-none object-contain'
                    style={{
                        animation: `date-heart-rain ${heart.duration}s linear ${heart.delay}s infinite backwards`,
                        height: `${heart.size}px`,
                        left: `${heart.left}%`,
                        opacity: 0,
                        width: `${heart.size}px`,
                        ['--heart-rotate' as string]: `${heart.rotate}deg`,
                        ['--heart-sway' as string]: `${heart.sway}px`,
                    }}
                />
            ))}
        </div>
    );
}

export function DateMeme() {
    usePageMeta({
        title: 'Date? | Lesi',
        description: 'A deeply serious yes or no question',
    });
    const { effectsEnabled, setEffectsEnabled } = useSeason();

    const [questionIndex, setQuestionIndex] = useState(0);
    const [noPositionIndex, setNoPositionIndex] = useState(0);
    const [noLabelIndex, setNoLabelIndex] = useState(0);
    const [noHoverCount, setNoHoverCount] = useState(0);
    const [colourMode, setColourMode] = useState<ColourMode>(() => getColourMode());
    const [accepted, setAccepted] = useState(false);
    const suppressYesClickUntilRef = useRef(0);

    const noPosition = noPositions[noPositionIndex];
    const activePrompt = accepted ? acceptedPrompt : noHoverCount === 0 ? initialPrompt : prompts[questionIndex];
    const question = activePrompt.question;
    const noLabel = noLabels[noLabelIndex];
    const yesScale = Math.min(1 + noHoverCount * 0.14, 2.05);
    const hasMovedNoButton = noHoverCount > 0;
    const palette = palettes[colourMode];

    useEffect(() => {
        setColourMode(getColourMode());
    }, []);

    function moveNoButton() {
        if (accepted) {
            return;
        }

        setNoPositionIndex((current) => (current + 1) % noPositions.length);
        setQuestionIndex((current) => (noHoverCount === 0 ? 0 : (current + 1) % prompts.length));
        setNoLabelIndex((current) => (current + 1) % noLabels.length);
        setNoHoverCount((current) => current + 1);
    }

    function suppressAccidentalYesClick() {
        suppressYesClickUntilRef.current = Date.now() + 450;
    }

    function acceptDate() {
        setEffectsEnabled(true);
        setAccepted(true);
    }

    const yesButton = (
        <button
            type='button'
            className={sharedButtonClassName}
            style={{
                background: palette.yes,
                boxShadow: palette.buttonShadow,
                color: palette.yesText,
                transform: `scale(${yesScale})`,
            }}
            onMouseEnter={(event) => {
                event.currentTarget.style.background = palette.yesHover;
            }}
            onMouseLeave={(event) => {
                event.currentTarget.style.background = palette.yes;
            }}
            onClick={(event) => {
                if (Date.now() < suppressYesClickUntilRef.current) {
                    event.preventDefault();
                    event.stopPropagation();
                    return;
                }

                acceptDate();
            }}>
            Yes
        </button>
    );

    const noButton = (
        <button
            type='button'
            className={`${sharedButtonClassName} ${hasMovedNoButton ? 'fixed z-30' : ''}`}
            style={{
                background: palette.no,
                border: `2px solid ${palette.noBorder}`,
                boxShadow: palette.buttonShadow,
                color: palette.noText,
                left: hasMovedNoButton ? `clamp(86px, ${noPosition.x}vw, calc(100vw - 86px))` : undefined,
                top: hasMovedNoButton ? `clamp(34px, ${noPosition.y}dvh, calc(100dvh - 34px))` : undefined,
                transform: hasMovedNoButton ? `translate(-50%, -50%) rotate(${noPosition.rotate}deg)` : 'scale(1)',
            }}
            onPointerEnter={(event) => {
                if (event.pointerType === 'mouse') {
                    moveNoButton();
                }
            }}
            onPointerDown={(event) => {
                event.preventDefault();
                event.stopPropagation();
                suppressAccidentalYesClick();
                try {
                    event.currentTarget.setPointerCapture(event.pointerId);
                } catch {
                    // The button may move before all browsers allow pointer capture.
                }
                moveNoButton();
            }}
            onFocus={() => {
                if (Date.now() >= suppressYesClickUntilRef.current) {
                    moveNoButton();
                }
            }}
            onClick={(event) => {
                event.preventDefault();
                event.stopPropagation();
                suppressAccidentalYesClick();
            }}>
            {noLabel}
        </button>
    );

    return (
        <section
            className='relative isolate flex min-h-screen w-full flex-col items-center justify-center overflow-hidden px-4 py-10 text-center sm:px-8 '
            style={{
                background: palette.background,
                color: palette.question,
            }}>
            {accepted && effectsEnabled ? <ValentineHeartRain /> : null}
            <div className='relative z-10 grid min-h-[min(92dvh,780px)] w-full max-w-[760px] grid-rows-[clamp(180px,32dvh,260px)_clamp(128px,18dvh,170px)_clamp(230px,34dvh,320px)] items-center gap-5 px-4 py-6 sm:px-8'>
                <div className='flex h-full w-full items-center justify-center'>
                    {activePrompt.image ? (
                        <img
                            src={activePrompt.image.src}
                            alt={activePrompt.image.alt}
                            className='h-full max-w-full rounded-[18px] object-contain'
                            style={{ boxShadow: palette.imageShadow }}
                        />
                    ) : null}
                </div>

                <p className='flex h-full items-center justify-center text-balance text-[clamp(2rem,7vw,3.5rem)] font-bold leading-tight'>
                    {question}
                </p>

                {accepted ? (
                    <p
                        className='flex h-full items-start justify-center pt-8 text-[clamp(1.25rem,4vw,1.75rem)] font-semibold'
                        style={{ color: palette.muted }}>
                        Call me on my cell baby girl 🫦
                    </p>
                ) : (
                    <div className='relative flex h-full w-full items-start justify-center pt-10'>
                        <div className='flex min-h-[96px] w-full max-w-[420px] flex-row items-center justify-center gap-5 sm:gap-8'>
                            {yesButton}
                            {hasMovedNoButton ? <div className={sharedButtonClassName} aria-hidden='true' /> : noButton}
                        </div>
                        {hasMovedNoButton ? noButton : null}
                    </div>
                )}
            </div>
        </section>
    );
}
