import { useEffect, useRef, useState } from 'react';
import { Button, Icons, Loaders } from '@/components/ui';
import { cn } from '@/utils';
import { Link } from 'react-router-dom';
import { usePageMeta } from '@/hooks';

declare global {
    interface Window {
        createUnityInstance?: (
            canvas: HTMLCanvasElement,
            config: {
                arguments: string[];
                dataUrl: string;
                frameworkUrl: string;
                codeUrl: string;
                streamingAssetsUrl: string;
                companyName: string;
                productName: string;
                productVersion: string;
                showBanner: (msg: string, type: 'error' | 'warning') => void;
            },
            onProgress?: (progress: number) => void
        ) => Promise<{
            SetFullscreen: (value: 0 | 1) => void;
        }>;
        loadGameSettings: (key: string) => string;
        saveGameSettings: (key: string, value: string) => void;
        loadScores: (key: string) => string;
        saveTopScore: (key: string, value: string) => void;
        getUsername: () => string;
        loadSensModifier: () => string;
    }
}

type SessionResponse = {
    message: Session;
    error: null;
};

type Session = {
    id: string;
    twitchUserId: string;
    twitchLogin: string;
    twitchDisplayName: string | null;
    twitchAvatarUrl: string | null;
};

export function AimTrainer() {
    usePageMeta({ title: 'Aim Trainer | Lesi', description: 'Aim Trainer Game' });
    const canvasRef = useRef<HTMLCanvasElement | null>(null);
    const [loadingProgress, setLoadingProgress] = useState<number>(0);
    const [unityInstance, setUnityInstance] = useState<{ SetFullscreen: (value: 0 | 1) => void } | null>(null);
    const [isFullscreen, setIsFullscreen] = useState<boolean>(false);
    const [me, setMe] = useState<Session | null>(null);
    const [windowWidth, setWindowWidth] = useState<number>(window.innerWidth);

    const screenWidthMin = 1080;

    useEffect(() => {
        function handleResize() {
            setWindowWidth(window.innerWidth);
        }
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    useEffect(() => {
        if (window.innerWidth < screenWidthMin) {
            return;
        }

        async function loadMe() {
            try {
                const res = await fetch('/api/auth/twitch/me', { credentials: 'include' });
                if (!res.ok) {
                    setMe(null);
                    window.getUsername = () => '';
                    return;
                }
                const session = (await res.json()) as SessionResponse;
                const data = session.message;
                setMe(data);
                window.getUsername = () => data.twitchDisplayName ?? data.twitchLogin ?? '';
            } catch (error) {
                console.error('error', error);
                setMe(null);
                window.getUsername = () => '';
            }
        }

        loadMe();

        window.loadGameSettings = (key) => {
            const settings = JSON.parse(localStorage.getItem('gameSettings') ?? '{}') as Record<string, string>;
            return settings[key] ?? '';
        };

        window.saveGameSettings = (key, value) => {
            const settings = JSON.parse(localStorage.getItem('gameSettings') ?? '{}') as Record<string, string>;
            settings[key] = value;
            localStorage.setItem('gameSettings', JSON.stringify(settings));
        };

        window.loadScores = (key) => {
            const scores = JSON.parse(localStorage.getItem('gameScores') ?? '{}') as Record<string, string>;
            return scores[key] ?? '';
        };

        window.saveTopScore = (key, value) => {
            const scores = JSON.parse(localStorage.getItem('gameScores') ?? '{}') as Record<string, string>;
            scores[key] = value;
            localStorage.setItem('gameScores', JSON.stringify(scores));
        };

        window.loadSensModifier = () => {
            return '8';
        };

        async function loadUnity() {
            if (!canvasRef.current) {
                return;
            }

            const buildPath = '/aim-trainer/Build';
            const loaderUrl = `${buildPath}/aim-trainer.loader.js?v=2`;

            const config = {
                arguments: [],
                dataUrl: `${buildPath}/aim-trainer.data`,
                frameworkUrl: `${buildPath}/aim-trainer.framework.js`,
                codeUrl: `${buildPath}/aim-trainer.wasm`,
                streamingAssetsUrl: '/StreamingAssets',
                companyName: 'comp',
                productName: 'aim-trainer',
                productVersion: '0.1.0',
                showBanner: unityShowBanner,
            };

            const script = document.createElement('script');
            script.src = loaderUrl;
            script.async = true;

            script.onload = () => {
                if (!window.createUnityInstance) {
                    return;
                }
                window
                    .createUnityInstance(canvasRef.current as HTMLCanvasElement, config, (progress) => {
                        setLoadingProgress(progress);
                    })
                    .then((instance) => {
                        setUnityInstance(instance);
                    })
                    .catch(() => {});
            };

            document.body.appendChild(script);
        }

        loadUnity();
    }, []);

    function unityShowBanner(msg: string, type: 'error' | 'warning') {
        alert(`${type.toUpperCase()}: ${msg}`);
    }

    function handleEnterFullscreen() {
        if (canvasRef.current) {
            canvasRef.current.focus();
        }
        if (unityInstance) {
            unityInstance.SetFullscreen(1);
            setIsFullscreen(true);
        }
    }

    async function signOut() {
        await fetch('/api/auth/twitch/logout', { method: 'POST', credentials: 'include' });
        location.reload();
    }

    function startTwitchLogin() {
        const returnTo = window.location.pathname;
        window.location.href = `/auth/twitch?returnTo=${encodeURIComponent(returnTo)}`;
    }

    if (windowWidth < screenWidthMin) {
        return (
            <div className='flex flex-col gap-2'>
                <h1 className='text-lg font-bold text-error'>Sorry!</h1>
                <span>The game is currently only supported on window sizes larger than 1080 pixels</span>
                <span>
                    Please try increasing the size of your browser window or if you&#39;re on mobile, please access this
                    page on a desktop/laptop
                </span>
            </div>
        );
    }

    return (
        <>
            <h2 className='flex items-center justify-end gap-4 text-pretty text-lg'></h2>

            <div
                className={cn(
                    'group relative flex h-full w-full items-center justify-center rounded bg-base-300 motion-safe:animate-pulse xl:h-[540px] xl:w-[960px]',
                    unityInstance ? 'hidden' : 'visible'
                )}>
                <div tabIndex={0} className='peer flex h-full w-full items-center justify-center'>
                    <Loaders.Trio />
                </div>
                <ToolTip isLoggedIn={me !== null} />
                <GameButtons
                    isLoggedIn={me !== null}
                    unityInstance={unityInstance}
                    signOut={() => {
                        void signOut();
                    }}
                    startLogin={startTwitchLogin}
                    handleEnterFullscreen={handleEnterFullscreen}
                />
            </div>

            <div className={cn('group relative h-[540px] w-[960px]', unityInstance ? 'visible' : 'hidden')}>
                <canvas
                    ref={canvasRef}
                    id='unity-canvas'
                    className={cn(
                        'peer max-h-full w-full min-w-full max-w-full rounded xl:h-[540px]',
                        unityInstance ? 'visible' : 'hidden'
                    )}
                    tabIndex={-1}
                />
                <ToolTip isLoggedIn={me !== null} />
                <GameButtons
                    isLoggedIn={me !== null}
                    unityInstance={unityInstance}
                    signOut={() => {
                        void signOut();
                    }}
                    startLogin={startTwitchLogin}
                    handleEnterFullscreen={handleEnterFullscreen}
                />
            </div>
        </>
    );
}

function ToolTip({ isLoggedIn }: { isLoggedIn: boolean }) {
    return (
        <div className='absolute top-0 hidden h-1/3 w-full items-start justify-end rounded-t bg-gradient-to-b from-black/70 to-transparent p-4 group-hover:inline-flex peer-focus-within:hidden'>
            <div className='relative flex w-[34%] items-start justify-end'>
                <div className='peer flex h-20 w-10 justify-end'>
                    <Icons.Info className='peer text-info' />
                </div>
                <div className='peer absolute left-0 top-0 z-20 hidden w-[280px] flex-col gap-2 rounded bg-base-100 p-4 text-sm hover:flex peer-hover:flex'>
                    <h1 className='w-full pt-4 text-center text-base underline decoration-accent decoration-2 underline-offset-2'>
                        Keybinds
                    </h1>
                    <div className='flex flex-row items-center justify-between'>
                        <span className='w-[40%]'>Menu</span>
                        <span className='flex w-1/4 items-center justify-center rounded border-2 border-accent py-2'>
                            TAB
                        </span>
                    </div>
                    <div className='flex flex-row items-center justify-between'>
                        <span className='w-[40%]'>Restart</span>
                        <span className='flex w-1/4 items-center justify-center rounded border-2 border-accent py-2'>
                            G
                        </span>
                        <span className='flex w-1/4 items-center justify-center rounded border-2 border-accent py-2'>
                            R
                        </span>
                    </div>
                    <h1 className='w-full pt-4 text-center text-base underline decoration-accent decoration-2 underline-offset-2'>
                        Info
                    </h1>
                    <div className='flex w-full flex-col items-center justify-between gap-4'>
                        <span>You can full screen the game by clicking on the icon at the bottom right</span>
                        {!isLoggedIn ? (
                            <span className='w-full text-pretty'>
                                To update the leaderboard, you must login with Twitch
                            </span>
                        ) : null}
                        <span>
                            If you are getting low frame rates, you may need to enable{' '}
                            <Link target='_blank' to='/docs/browser-gpu-acceleration'>
                                <Button variant='link' size='link' className='hover:text-accent'>
                                    Hardware Acceleration
                                </Button>
                            </Link>
                        </span>
                        <span className='w-full'>
                            <Link target='_blank' to='/aim-trainer/release-notes' className='text-left'>
                                <Button variant='link' size='link' className='hover:text-accent'>
                                    New Release Notes
                                </Button>
                            </Link>
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
}

function GameButtons({
    isLoggedIn,
    unityInstance,
    signOut,
    startLogin,
    handleEnterFullscreen,
}: {
    isLoggedIn: boolean;
    unityInstance: { SetFullscreen: (value: 0 | 1) => void } | null;
    signOut: () => void;
    startLogin: () => void;
    handleEnterFullscreen: () => void;
}) {
    return (
        <div className='absolute bottom-0 hidden h-1/3 w-full items-end justify-between rounded-b bg-gradient-to-t from-black/70 to-transparent p-4 group-hover:inline-flex peer-focus-within:hidden'>
            {isLoggedIn ? (
                <Button size='icon' variant='accent' onClick={signOut}>
                    <Icons.LogOut className='text-white' />
                </Button>
            ) : (
                <Button
                    className={cn('w-fit gap-4 bg-twitchPurple text-white hover:bg-twitchPurple/80')}
                    onClick={startLogin}>
                    <Icons.Socials.Twitch width={20} height={20} /> Login With Twitch
                </Button>
            )}
            {unityInstance ? (
                <Button
                    variant='accent'
                    size='icon'
                    className={cn('p-2', isLoggedIn ? 'null' : 'bg-twitchPurple hover:bg-twitchPurple/80')}
                    onClick={handleEnterFullscreen}>
                    <Icons.FullScreen width={20} height={20} className='text-white' />
                </Button>
            ) : null}
        </div>
    );
}
