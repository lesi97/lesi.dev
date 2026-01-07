import { Description, Button, Input, Checkbox } from '@/components/ui';
import { useSeason } from '@/context/SeasonContext';
import { type MouseEvent } from 'react';

function changeTheme(e: MouseEvent<HTMLButtonElement>) {
    const theme = e.currentTarget.dataset.theme;
    if (!theme) {
        return;
    }
    document.documentElement.setAttribute('data-theme', theme);
    const themeData = {
        theme,
        'user-selected': true,
    };
    localStorage.setItem('theme', JSON.stringify(themeData));
}

const themes = [
    { name: 'bumblebee', label: 'Bumblebee' },
    { name: 'retro', label: 'Retro' },
    { name: 'valentine', label: 'Valentines' },
    { name: 'halloween', label: 'Halloween' },
    { name: 'garden', label: 'Garden' },
    { name: 'lofi', label: 'Lofi' },
    { name: 'pastel', label: 'Pastel' },
    { name: 'dracula', label: 'Dracula' },
    { name: 'cmyk', label: 'CMYK' },
    { name: 'autumn', label: 'Autumn' },
    { name: 'business', label: 'Business' },
    { name: 'acid', label: 'Acid' },
    { name: 'coffee', label: 'Coffee' },
    { name: 'winter', label: 'Winter' },
    { name: 'dim', label: 'Dim' },
    { name: 'nord', label: 'Nord' },
    { name: 'sunset', label: 'Sunset' },
    { name: 'gear5', label: 'Gear 5' },
    { name: 'naruto', label: 'Naruto' },
    { name: 'lesi-default', label: 'Lesi' },
];

export function Settings() {
    const { season, effectsEnabled, setEffectsEnabled } = useSeason();

    function toggleSeasonEffects() {
        setEffectsEnabled(!effectsEnabled);
    }

    return (
        <>
            <Description title='Settings' subtitle='Manage your settings' />
            {season && (
                <label className='mb-4 flex w-fit cursor-pointer flex-row items-center gap-4'>
                    {season} Background Effects:
                    <Checkbox size='sm' onChange={toggleSeasonEffects} checked={effectsEnabled} />
                </label>
            )}
            <h1 className='mb-1'>Theme:</h1>
            <div className='grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-6'>
                {themes.map((theme) => (
                    <Button data-theme={theme.name} onClick={changeTheme} className='min-w-fit'>
                        {theme.label}
                    </Button>
                ))}
            </div>
        </>
    );
}
