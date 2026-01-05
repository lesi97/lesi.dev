'use client';
import { useEffect, useRef, useState } from 'react';
import { Button, Icons, Loaders } from '@/components/ui';
// import signInOauth from '@/lib/auth/signIn-OAuth';
// import createClient from '@/lib/supabase/client';
// import { Session } from '@supabase/supabase-js';
import { mergeClassNames } from '@/utils';
import { Link } from 'react-router-dom';

declare global {
    interface Window {
        createUnityInstance?: (
            canvas: HTMLCanvasElement,
            config: any,
            onProgress?: (progress: number) => void
        ) => Promise<any>;
        loadGameSettings: (key: string) => string;
        saveGameSettings: (key: string, value: string) => void;
        loadScores: (key: string) => string;
        saveTopScore: (key: string, value: string) => void;
        getUsername: () => string;
        loadSensModifier: () => string;
    }
}

const isDev = process.env.NODE_ENV === 'development';

export function AimTrainer() {
    const canvasRef = useRef<HTMLCanvasElement | null>(null);
    const [loadingProgress, setLoadingProgress] = useState<number>(0);
    const [unityInstance, setUnityInstance] = useState<any>(null);
    const [isFullscreen, setIsFullscreen] = useState<boolean>(false);
    // const [session, setSession] = useState<Session | null>(null);
    const [windowWidth, setWindowWidth] = useState(window?.innerWidth);

    // const supabase = createClient();

    const screenWidthMin = 1080;

    useEffect(() => {
        function handleResize() {
            setWindowWidth(window.innerWidth);
        }
        window.addEventListener('resize', handleResize);

        return () => window.removeEventListener('resize', handleResize);
    }, []);

    useEffect(() => {
        if (window && window.innerWidth < screenWidthMin) {
            return;
        }
        // supabase.auth.getSession().then(({ data }) => {
        //     setSession(data.session);
        //     if (data.session) {
        //         window.getUsername = () => {
        //             return data.session?.user.user_metadata.name || '';
        //         };
        //     } else {
        //         window.getUsername = () => {
        //             return '';
        //         };
        //     }
        // });

        window.loadGameSettings = (key) => {
            const settings = JSON.parse(localStorage.getItem('gameSettings') || '{}');
            return settings[key] || '';
        };

        window.saveGameSettings = (key, value) => {
            const settings = JSON.parse(localStorage.getItem('gameSettings') || '{}');
            settings[key] = value;
            localStorage.setItem('gameSettings', JSON.stringify(settings));
        };

        window.loadScores = (key) => {
            const scores = JSON.parse(localStorage.getItem('gameScores') || '{}');
            return scores[key] || '';
        };

        window.saveTopScore = (key, value) => {
            const scores = JSON.parse(localStorage.getItem('gameScores') || '{}');
            scores[key] = value;
            localStorage.setItem('gameScores', JSON.stringify(scores));
        };

        window.loadSensModifier = () => {
            // Higher number lower overall sens
            return '8';
        };

        async function loadUnity() {
            if (!canvasRef.current) return;
            const buildPath = '/_static/aim-trainer/Build';
            const loaderUrl = `${buildPath}/aim-trainer.loader.js?v=2`;

            const config = {
                arguments: [],
                dataUrl: `${buildPath}/aim-trainer.data.br`,
                frameworkUrl: `${buildPath}/aim-trainer.framework.js.br`,
                codeUrl: `${buildPath}/aim-trainer.wasm.br`,
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
                if (window.createUnityInstance) {
                    window
                        .createUnityInstance(canvasRef.current as HTMLCanvasElement, config, (progress) => {
                            setLoadingProgress(progress);
                        })
                        .then((instance) => {
                            setUnityInstance(instance);
                        })
                        .catch((message) => {
                            console.error('Unity Loading Error:', message);
                        });
                } else {
                    console.error('Unity loader script did not define createUnityInstance');
                }
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

    // async function signOut() {
    //     const { error } = await supabase.auth.signOut();
    //     location.reload();
    // }

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
            <div className='flex-1 overflow-hidden w-fit place-self-center'>
                {/* The canvas ref must exist for unityInstance to become populated so rather than an early return, hide elements depending on the instance state */}
                <div
                    className={mergeClassNames(
                        'group relative flex h-full w-full min-w-[1180px] items-center justify-center rounded bg-base-300 motion-safe:animate-pulse xl:h-full xl:w-full',
                        unityInstance ? 'hidden' : 'visible'
                    )}>
                    <div tabIndex={0} className='peer flex h-full w-full items-center justify-center'>
                        <Loaders.Trio />
                    </div>
                    {/* <ToolTip session={session} /> */}
                    <ToolTip session={null} />
                    <GameButtons
                        // session={session}
                        session={null}
                        unityInstance={unityInstance}
                        // signOut={() => signOut()}
                        signOut={() => {}}
                        handleEnterFullscreen={() => handleEnterFullscreen()}
                    />
                </div>

                <div
                    className={mergeClassNames(
                        'group relative h-full w-full items-center place-content-center place-items-center',
                        unityInstance ? 'visible' : 'hidden'
                    )}>
                    <canvas
                        ref={canvasRef}
                        id='unity-canvas'
                        className={mergeClassNames(
                            'max-h-[calc(100vh-4rem)] w-[min(100%,calc((100vh-4rem)*1.4))] min-w-[1180px] aspect-[9/16] peer max-w-full rounded h-full bg-base-100',
                            unityInstance ? 'visible' : 'hidden'
                        )}
                        tabIndex={-1}
                    />
                    <ToolTip session={null} />
                    {/* <ToolTip session={session} /> */}
                    <GameButtons
                        session={null}
                        // session={session}
                        unityInstance={unityInstance}
                        // signOut={() => signOut()}
                        signOut={() => {}}
                        handleEnterFullscreen={() => handleEnterFullscreen()}
                    />
                </div>
            </div>
        </>
    );
}

// function ToolTip({ session }: { session: Session | null }) {
function ToolTip({ session }: { session: null }) {
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
                        {session ? (
                            <span className='w-full text-pretty'>
                                To update the leaderboard, you must login with Twitch
                            </span>
                        ) : null}
                        <span>
                            If you are getting low frame rates, you may need to enable{' '}
                            <Link target='_blank' to='/h/browser-gpu-acceleration'>
                                <Button variant='link' size='link' className='hover:text-accent'>
                                    Hardware Acceleration
                                </Button>
                            </Link>
                        </span>
                        <span className='w-full'>
                            <Link target='_blank' to='/g/aim-trainer/release-notes' className='text-left'>
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
    session,
    unityInstance,
    signOut,
    handleEnterFullscreen,
}: {
    // session: Session | null;
    session: null;
    unityInstance: string;
    signOut: () => void;
    handleEnterFullscreen: () => void;
}) {
    return (
        <div className='absolute bottom-0 hidden h-1/3 w-full items-end justify-between rounded-b bg-gradient-to-t from-black/70 to-transparent p-4 group-hover:inline-flex peer-focus-within:hidden'>
            {session ? (
                <Button size='icon' variant='accent' onClick={signOut}>
                    <Icons.LogOut className='text-white' />
                </Button>
            ) : (
                <Button
                    className={mergeClassNames('w-fit gap-4 bg-twitchPurple text-white hover:bg-twitchPurple/80')}
                    // onClick={() => {
                    //     signInOauth(
                    //         'twitch',
                    //         isDev
                    //             ? `http://localhost:3050/api/v1/twitch/callback?next=${encodeURIComponent('/g/aim-trainer')}`
                    //             : `https://lesi.dev/api/v1/twitch/callback?next=${encodeURIComponent('/g/aim-trainer')}`
                    //     );
                    // }}
                >
                    <Icons.Socials.Twitch width={20} height={20} /> Login With Twitch
                </Button>
            )}
            {unityInstance ? (
                <Button
                    variant='accent'
                    size='icon'
                    className={mergeClassNames('p-2', session ? 'null' : 'bg-twitchPurple hover:bg-twitchPurple/80')}
                    onClick={handleEnterFullscreen}>
                    <Icons.FullScreen width={20} height={20} className='text-white' />
                </Button>
            ) : null}
        </div>
    );
}
