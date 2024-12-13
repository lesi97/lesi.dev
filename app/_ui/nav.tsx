import Link from 'next/link';
import lesi from '@/app/_assets/_images/lesi.webp';
import lesiValentines from '@/app/_assets/_images/lesi-valentines.webp';
import lesiHalloween from '@/app/_assets/_images/lesi-halloween.webp';
import lesiBirthday from '@/app/_assets/_images/lesi-birthday.webp';
import lesiChristmas from '@/app/_assets/_images/lesi-christmas.webp';
import lesiNewYear from '@/app/_assets/_images/lesi-newyear.webp';
import { ReactElement } from 'react';
import Image from 'next/image';

const links = [
    { href: '/', label: 'Home' },
    { href: '/aspect-ratio-calculator', label: 'Aspect Ratio Calculator' },
    { href: '/pdf-to-png', label: 'PDF Converter' },
    { href: '/ico-converter', label: 'Icon Converter' },
    { href: '/password-generator', label: 'Password Generator' },
    { href: '/weight-converter', label: 'Weight Converter' },
    { href: '/video-to-mp3', label: 'Video To MP3' },
    { href: '/video-cropper', label: 'Video Cropper' },
    { href: '/minifier', label: 'Minifier' },
    { href: '/settings', label: 'Settings' },
];

const Nav = ({ season }: { season: string | null }) => {
    function getSeasonImage(season: string | null): ReactElement {
        switch (season) {
            case 'Valentines':
                return <Image src={lesiValentines.src} alt="Lesi-Valentines" height={47} width={47} />;
            case 'Halloween':
                return <Image src={lesiHalloween.src} alt="Lesi-Halloween" height={47} width={47} />;
            case 'Lesi-Birthday':
                return <Image src={lesiBirthday.src} alt="Lesi-Birthday" height={47} width={47} />;
            case 'Christmas':
                return <Image src={lesiChristmas.src} alt="Lesi-Christmas" height={47} width={47} />;
            case 'New-Years':
                return <Image src={lesiNewYear.src} alt="Lesi-New-Years" height={47} width={47} />;
            default:
                return <Image src={lesi.src} alt="Lesi" height={47} width={47} />;
        }
    }

    return (
        <header>
            <nav className="w-full flex flex-row h-47px max-h-47px relative bg-primary">
                <Link className="w-47px h-47px relative overflow-hidden" href="/">
                    {getSeasonImage(season)}
                </Link>

                <ul className="m-0 p-0 flex flex-row h-full max-h-47px">
                    {links.map((obj) => {
                        return (
                            <li key={obj.href} className="list-none flex items-center max-h-47px">
                                <Link
                                    href={obj.href}
                                    className="text-[var(--nav-font-colour)] max-h-47px text-center flex items-center px-4 h-full w-full no-underline text-[17px] focus-visible:text-[var(--accent-colour)] hover:bg-[var(--nav-active-background-colour)] hover:text-[var(--nav-active-font-colour)]"
                                >
                                    {obj.label}
                                </Link>
                            </li>
                        );
                    })}
                </ul>
            </nav>
        </header>
    );
};

export default Nav;
