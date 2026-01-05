import { ReactElement, useState } from 'react';
import { useSeason } from '@/context/SeasonContext';
import { Icons } from '../ui';
import { mergeClassNames } from '@/utils';
import { Link } from 'react-router-dom';

const links = [
    {
        href: '/aspect-ratio-calculator',
        label: 'Tools',
        description: '',
        icon: <Icons.AllSeeingEye />,
    },
    { href: '/aim-trainer', label: 'Aim Trainer', description: '', icon: <Icons.CrossHair /> },
    {
        href: '/file-converters',
        label: 'File Converters',
        children: [
            { href: '/pdf-to-png', label: 'PDF To PNG' },
            { href: '/ico-converter', label: 'Image to Icon' },
            { href: '/video-to-mp3', label: 'Video To MP3' },
            { href: '/video-cropper', label: 'Video Cropper' },
        ],
    },
];

export function Nav() {
    const [navExpanded, setNavExpanded] = useState(false);
    return (
        <nav className='flex flex-row w-full h-[47px] items-center'>
            <button className='px-4 xl:hidden' onClick={() => setNavExpanded(!navExpanded)}>
                <Icons.Burger />
            </button>

            <ul className={mergeClassNames('w-full h-full gap-4 flex flex-row items-center justify-center')}>
                {links.map((link) => {
                    return (
                        <Link to={link.href}>
                            <div className='h-full w-10 text-primary'>{link.icon}</div>
                        </Link>
                    );
                })}
            </ul>
        </nav>
    );
}

// export function Nav() {
//     const { season } = useSeason();
//     const [navExpanded, setNavExpanded] = useState(false);
//     function getSeasonImage(season: string | null): ReactElement {
//         switch (season) {
//             case 'Valentines':
//                 return (
//                     <img
//                         src='/_static/images/lesi-valentines.webp'
//                         alt='Lesi-Valentines'
//                         height={47}
//                         width={47}
//                         className='h-full w-auto'
//                     />
//                 );
//             case 'Halloween':
//                 return (
//                     <img
//                         src='/_static/images/lesi-halloween.webp'
//                         alt='Lesi-Halloween'
//                         height={47}
//                         width={47}
//                         className='h-full w-auto'
//                     />
//                 );
//             case 'Lesi-Birthday':
//                 return (
//                     <img
//                         src='/_static/images/lesi-birthday.webp'
//                         alt='Lesi-Birthday'
//                         height={47}
//                         width={47}
//                         className='h-full w-auto'
//                     />
//                 );
//             case 'Christmas':
//                 return (
//                     <img
//                         src='/_static/images/lesi-christmas.webp'
//                         alt='Lesi-Christmas'
//                         height={47}
//                         width={47}
//                         className='h-full w-auto'
//                     />
//                 );
//             case 'New-Years':
//                 return (
//                     <img
//                         src='/_static/images/lesi-newyear.webp'
//                         alt='Lesi-New-Years'
//                         height={47}
//                         width={47}
//                         className='h-full w-auto'
//                     />
//                 );
//             default:
//                 return (
//                     <img src='/_static/images/lesi.webp' alt='Lesi' height={47} width={47} className='h-full w-auto' />
//                 );
//         }
//     }

//     return (
//         <nav className='relative flex h-47px w-full flex-row items-center bg-primary'>
//             <Link className='relative h-full w-fit overflow-hidden' to='/'>
//                 {getSeasonImage(season)}
//             </Link>

//             <button className='px-4 xl:hidden' onClick={() => setNavExpanded(!navExpanded)}>
//                 <Icons.Burger />
//             </button>

//             <ul
//                 className={mergeClassNames(
//                     'h-fit w-full bg-inherit p-0 xl:relative xl:top-0 xl:!flex xl:h-full xl:max-h-47px xl:flex-row items-center',
//                     navExpanded ? 'flex flex-col' : 'hidden'
//                 )}>
//                 {links.map((obj) => {
//                     return (
//                         <li
//                             key={obj.href}
//                             className='group flex max-h-47px list-none items-center text-primary-content hover:bg-secondary/30'>
//                             <div className='relative h-full w-full'>
//                                 <Link
//                                     to={obj.href}
//                                     onClick={() => setNavExpanded(false)}
//                                     className='relative flex h-full max-h-[47px] w-full flex-row items-center gap-2 px-4 py-4 text-center text-[17px] text-primary-content no-underline xl:py-0'>
//                                     <span>{obj.label}</span>
//                                     {obj.children && obj.children.length !== 0 ? (
//                                         <Icons.Chevron className='hidden h-4 w-4 rotate-180 text-accent transition-transform duration-500 ease-in-out group-hover:rotate-0 xl:inline-flex' />
//                                     ) : null}
//                                 </Link>
//                                 {obj.children && obj.children.length !== 0 ? (
//                                     <div className='min-w-parent absolute left-0 top-[47px] hidden w-max min-w-full transform divide-y divide-solid divide-neutral shadow-xl xl:flex-col xl:group-hover:flex'>
//                                         {obj.children.map((child, idx) => {
//                                             return (
//                                                 <Link
//                                                     key={idx}
//                                                     to={child.href}
//                                                     className='outline-l-0 flex h-[47px] w-full items-center bg-secondary px-4 py-4 text-center text-[17px] text-secondary-content no-underline last:rounded-b-md hover:bg-secondary/90 xl:py-0'>
//                                                     {child.label}
//                                                 </Link>
//                                             );
//                                         })}
//                                     </div>
//                                 ) : null}
//                             </div>
//                         </li>
//                     );
//                 })}
//             </ul>
//         </nav>
//     );
// }
