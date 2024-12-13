const { styleText } = require('util');

/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ['./pages/**/*.{js,ts,jsx,tsx}', './components/**/*.{js,ts,jsx,tsx}', './app/**/*.{js,ts,jsx,tsx}'],
    darkMode: 'selector',
    daisyui: {
        themes: [
            'light',
            'dark',
            'cupcake',
            'forest',
            {
                'lesi-dark': {
                    primary: '#1f2429',
                    'primary-content': '#000916',
                    secondary: '#0094c3',
                    'secondary-content': '#00080e',
                    accent: '#0000ff',
                    'accent-content': '#c6dbff',
                    neutral: '#181818',
                    'neutral-content': '#cbcbcb',
                    'base-100': '#151a1f',
                    'base-200': '#13171c',
                    'base-300': '#111519',
                    'base-content': '#161615',
                    info: '#0088a6',
                    'info-content': '#00060a',
                    success: '#29c840',
                    'success-content': '#001603',
                    warning: '#febc2e',
                    'warning-content': '#160200',
                    error: '#ff5f57',
                    'error-content': '#160004',
                },
            },
        ],
        styled: false,
    },
    theme: {
        screens: {
            sm: '480px',
            md: '768px',
            lg: '976px',
            xl: '1440px',
        },
        extend: {
            colors: {
                light: {
                    background: {
                        nav: '#1265be',
                    },
                },
                dark: {
                    background: {
                        nav: '#1f2429',
                        primary: '#1f2429',
                        secondary: '#333d57',
                        tertiary: '#684656',
                    },
                },
            },
            spacing: {
                '50%': '50%',
                '90%': '90%',
                '47px': '47px',
            },
        },
    },
    plugins: [require('daisyui')],
};
