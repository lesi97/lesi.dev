import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface SeasonContextType {
    season: string | null;
    effectsEnabled: boolean;
    setEffectsEnabled: (enabled: boolean) => void;
}

export const SeasonContext = createContext<SeasonContextType | undefined>(undefined);

export function SeasonProvider({ children }: { children: ReactNode }) {
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
        let parsedTheme = null;

        try {
            parsedTheme = theme ? JSON.parse(theme) : null;
        } catch (error) {
            console.error('Invalid theme in localStorage, resetting...', error);
            localStorage.removeItem('theme');
        }

        function applyTheme(themeName: string) {
            document.documentElement.setAttribute('data-theme', themeName);
            return;
        }

        // If user has selected a theme, respect their choice and apply no matter
        if (parsedTheme?.['user-selected']) {
            return applyTheme(parsedTheme.theme);
        }

        // If there is a season with corresponding theme and the user has not selected a theme, apply the seasonal theme
        const seasonalThemes = {
            Valentines: 'valentine',
            Halloween: 'halloween',
            Christmas: 'winter',
        } as const;
        type Season = keyof typeof seasonalThemes;
        const themeForSeason = seasonalThemes[currentSeason as Season];

        if (themeForSeason) {
            return applyTheme(themeForSeason);
        }

        // If user has not selected a theme and there is no seasonal theme, take into consideration their OS theme and apply a light/dark theme
        const systemTheme = window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'lesi-default' : 'gear5';

        const themeData = {
            theme: systemTheme,
            'user-selected': false,
        };

        applyTheme(systemTheme);
        localStorage.setItem('theme', JSON.stringify(themeData));
        // document.cookie = `theme=${JSON.stringify(themeData)}; path=/;`;
    }, []);

    useEffect(() => {
        localStorage.setItem('season', effectsEnabled.toString());
    }, [effectsEnabled]);

    function getSeason(): 'Valentines' | 'Halloween' | 'Lesi-Birthday' | 'Christmas' | 'New-Years' | null {
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

    return (
        <SeasonContext.Provider value={{ season, effectsEnabled, setEffectsEnabled }}>
            {children}
        </SeasonContext.Provider>
    );
}

export function useSeason() {
    const context = useContext(SeasonContext);
    if (context === undefined) {
        throw new Error('useSeason must be used within a SeasonProvider');
    }
    return context;
}
