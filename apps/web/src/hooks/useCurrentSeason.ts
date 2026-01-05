import { useState, useEffect } from 'react';

export default function useCurrentSeason() {
    const [season, setSeason] = useState<string | null>(null);
    const [effectsEnabled, setEffectsEnabled] = useState<boolean>(() => {
        if (typeof window !== 'undefined') {
            const storedEffects = localStorage.getItem('season');
            return storedEffects !== null ? storedEffects === 'true' : true;
        }
        return true;
    });

    useEffect(() => {
        const currentSeason = getSeason();
        setSeason(currentSeason);
        const theme = localStorage.getItem('theme');

        // if ((!theme || !theme['user-selected']) && currentSeason) {
        //     if (currentSeason === 'Valentines') {
        //         document.documentElement.setAttribute('data-theme', 'valentine');
        //     } else if (currentSeason === 'Halloween') {
        //         document.documentElement.setAttribute('data-theme', 'halloween');
        //     } else if (currentSeason === 'Christmas') {
        //         document.documentElement.setAttribute('data-theme', 'winter');
        //     }
        // } else if (!theme) {
        //     if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        //         const themeData = {
        //             theme: 'lesi-default',
        //             'user-selected': false,
        //         };
        //         document.documentElement.setAttribute('data-theme', 'lesi-default');
        //         localStorage.setItem('theme', JSON.stringify(themeData));
        //         document.cookie = `theme=${JSON.stringify(themeData)}; path=/;`;
        //     } else {
        //         const themeData = {
        //             theme: 'gear5',
        //             'user-selected': false,
        //         };
        //         document.documentElement.setAttribute('data-theme', 'gear5');
        //         localStorage.setItem('theme', JSON.stringify(themeData));
        //         document.cookie = `theme=${JSON.stringify(themeData)}; path=/;`;
        //     }
        // }
    }, []);

    useEffect(() => {
        localStorage.setItem('season', effectsEnabled.toString());
    }, [effectsEnabled]);

    function getSeason() {
        const date = new Date();
        const month = date.getMonth() + 1;
        const day = date.getDate();

        switch (true) {
            case month === 2 && day >= 7 && day <= 14:
                return 'Valentines';
            case month === 10 && day >= 15 && day <= 31:
                return 'Halloween';
            case month === 11 && day === 1:
                return 'Lesi-Birthday';
            case month === 12 && day >= 1 && day <= 26:
                return 'Christmas';
            case (month === 12 && day >= 27 && day <= 31) || (month === 1 && day === 1):
                return 'New-Years';
            default:
                return null;
        }
    }

    return { season, effectsEnabled, setEffectsEnabled };
}
