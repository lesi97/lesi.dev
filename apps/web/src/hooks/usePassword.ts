import { useState, useEffect } from 'react';
// import generatePassword from '@/src/app/(routes)/(default)/password-generator/helpers/generatePassword';
type usePasswordType = {
    password: string;
    setPassword: (password: string) => void;
    sliderVal: number;
    setSliderVal: (sliderVal: number) => void;
    includeNum: boolean;
    setIncludeNum: (includeNum: boolean) => void;
    includeSymbols: boolean;
    setIncludeSymbols: (includeSymbols: boolean) => void;
    setIsUserTyping: (isUserTyping: boolean) => void;
    generatePassword: (includeNum: boolean, includeSymbols: boolean, sliderVal: number) => string;
    copyPassword: (ref: React.RefObject<HTMLInputElement | null>) => void;
};

export function usePassword(): usePasswordType {
    const [password, setPassword] = useState('');
    const [sliderVal, setSliderVal] = useState(16);
    const [includeNum, setIncludeNum] = useState(true);
    const [includeSymbols, setIncludeSymbols] = useState(true);
    const [isUserTyping, setIsUserTyping] = useState(false);

    useEffect(() => {
        if (!isUserTyping) {
            setPassword(' ');
            setPassword(generatePassword(includeNum, includeSymbols, sliderVal));
        }
    }, [sliderVal, includeNum, includeSymbols]);

    function generatePassword(includeNum: boolean, includeSymbols: boolean, sliderVal: number): string {
        const numbers = '0123456789';
        const letters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        const symbols = '[!@#$£%^&*()+=-.,]';

        let chars = letters;
        if (includeNum) chars += numbers;
        if (includeSymbols) chars += symbols;

        const array = new Uint32Array(sliderVal);
        window.crypto.getRandomValues(array);

        let password = '';

        for (let i = 0; i < sliderVal; i++) {
            password += chars[array[i] % chars.length];
        }

        if (includeNum && !/\d/.test(password)) {
            const randomIndex = Math.floor(Math.random() * password.length);
            const randomNumber = numbers[Math.floor(Math.random() * numbers.length)];
            password = password.substring(0, randomIndex) + randomNumber + password.substring(randomIndex + 1);
        }

        if (includeSymbols && !/!@#$£%^&*()+=-.,/.test(password)) {
            const randomIndex = Math.floor(Math.random() * password.length);
            const randomSymbol = symbols[Math.floor(Math.random() * symbols.length)];
            password = password.substring(0, randomIndex) + randomSymbol + password.substring(randomIndex + 1);
        }

        return password;
    }

    function copyPassword(ref: React.RefObject<HTMLInputElement | null>) {
        const copyText = ref.current;
        if (copyText) {
            copyText.select();
            document.execCommand('copy');
        }
    }

    return {
        password,
        setPassword,
        sliderVal,
        setSliderVal,
        includeNum,
        setIncludeNum,
        includeSymbols,
        setIncludeSymbols,
        setIsUserTyping,
        generatePassword,
        copyPassword,
    };
}
