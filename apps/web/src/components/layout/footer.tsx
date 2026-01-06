import { Link } from 'react-router-dom';
import { Icons } from '../ui';
import { cn } from '@/utils';
import { useState, useEffect } from 'react';

export function Footer() {
    const [theme, setTheme] = useState<string | undefined>();

    useEffect(() => {
        if (typeof document !== 'undefined') {
            const htmlElement = document.documentElement;

            const currentTheme = htmlElement.getAttribute('data-theme');
            if (currentTheme) {
                setTheme(currentTheme);
            }

            const observer = new MutationObserver((mutations) => {
                mutations.forEach((mutation) => {
                    if (mutation.attributeName === 'data-theme') {
                        const updatedTheme = htmlElement.getAttribute('data-theme');
                        setTheme(updatedTheme || undefined);
                    }
                });
            });
            observer.observe(htmlElement, { attributes: true });

            return () => observer.disconnect();
        }
    }, []);

    const getColorClass = (defaultClass: string) => {
        if (!theme) return defaultClass;
        const mappings: Record<string, Record<string, string>> = {
            bumblebee: {
                'text-primary': 'text-neutral',
                'text-secondary': 'text-info',
                'text-accent': 'text-accent',
                'text-neutral': 'text-accent',
                'text-error': 'text-accent',
            },
            retro: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-info/40',
                'text-error': 'text-accent',
            },
            valentine: {
                'text-primary': 'text-secondary',
                'text-secondary': 'text-primary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-neutral',
                'text-error': 'text-neutral',
            },
            halloween: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-warning/70',
                'text-error': 'text-error/80',
            },
            garden: {
                'text-primary': 'text-secondary',
                'text-secondary': 'text-accent',
                'text-accent': 'text-neutral',
                'text-neutral': 'text-primary',
                'text-error': 'text-primary',
            },
            lofi: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-primary',
                'text-accent': 'text-primary',
                'text-neutral': 'text-primary',
                'text-error': 'text-primary',
            },
            pastel: {
                'text-primary': 'text-secondary',
                'text-secondary': 'text-accent',
                'text-accent': 'text-neutral',
                'text-neutral': 'text-primary',
                'text-error': 'text-primary',
            },
            dracula: {
                'text-primary': 'text-secondary',
                'text-secondary': 'text-accent',
                'text-accent': 'text-primary',
                'text-neutral': 'text-primary',
                'text-error': 'text-primary',
            },
            cmyk: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-warning',
                'text-neutral': 'text-neutral',
                'text-error': 'text-error',
            },
            autumn: {
                'text-primary': 'text-secondary',
                'text-secondary': 'text-accent',
                'text-accent': 'text-primary',
                'text-neutral': 'text-secondary',
                'text-error': 'text-accent',
            },
            business: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-error',
                'text-error': 'text-error',
            },
            acid: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-neutral',
                'text-neutral': 'text-neutral',
                'text-error': 'text-error',
            },
            coffee: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-success',
                'text-accent': 'text-warning',
                'text-neutral': 'text-warning',
                'text-error': 'text-error',
            },
            dim: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-info',
                'text-error': 'text-error',
            },
            nord: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-info',
                'text-error': 'text-error',
            },
            sunset: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-secondary',
                'text-accent': 'text-accent',
                'text-neutral': 'text-info',
                'text-error': 'text-error',
            },
            gear5: {
                'text-primary': 'text-primary',
                'text-secondary': 'text-info',
                'text-accent': 'text-warning',
                'text-neutral': 'text-accent',
                'text-error': 'text-accent',
            },
            naruto: {
                'text-primary': 'text-base-100',
                'text-secondary': 'text-base-100',
                'text-accent': 'text-base-100',
                'text-neutral': 'text-base-100',
                'text-error': 'text-base-100',
            },
            'lesi-default': {
                'text-primary': 'text-base-content',
                'text-secondary': 'text-base-content',
                'text-accent': 'text-base-content',
                'text-neutral': 'text-base-content',
                'text-error': 'text-base-content',
            },
        };

        return mappings[theme || '']?.[defaultClass] || defaultClass;
    };

    const socials = [
        {
            href: 'https://www.youtube.com/@C_Lesi',
            icon: () => (
                <Icons.Socials.YouTube
                    width={32}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-error')
                    )}
                />
            ),
        },
        {
            href: 'https://twitch.tv/c_lesi',
            icon: () => (
                <Icons.Socials.Twitch
                    width={25}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-primary')
                    )}
                />
            ),
        },
        {
            href: 'https://discord.gg/RUDhkXT',
            icon: () => (
                <Icons.Socials.Discord
                    width={30}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-secondary')
                    )}
                />
            ),
        },
        {
            href: 'https://x.com/Chris_Lesi',
            icon: () => (
                <Icons.Socials.X
                    width={22}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-accent')
                    )}
                />
            ),
        },
        {
            href: 'https://www.instagram.com/christian_lesi',
            icon: () => (
                <Icons.Socials.Instagram
                    width={25}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-primary')
                    )}
                />
            ),
        },
        {
            href: 'https://www.tiktok.com/@c_lesi',
            icon: () => (
                <Icons.Socials.TikTok
                    width={32}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-secondary')
                    )}
                />
            ),
        },
        {
            href: 'https://www.buymeacoffee.com/lesi',
            icon: () => (
                <Icons.Socials.BuyMeACoffe
                    width={19}
                    className={cn(
                        'drop-shadow-md duration-200 ease-in-out hover:scale-110',
                        getColorClass('text-neutral')
                    )}
                />
            ),
        },
        // { href: 'mailto:chris.lesi001@gmail.com', icon: <Icons.Socials.E /> },
    ];
    return (
        <footer className='my-4 flex h-fit w-fit flex-row items-center gap-5'>
            {socials.map((social, idx) => {
                return (
                    <Link key={idx} to={social.href}>
                        {social.icon()}
                    </Link>
                );
            })}
        </footer>
    );
}
