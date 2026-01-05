import { useState, useEffect } from 'react';
import generatePassword from '@/src/app/(routes)/(default)/password-generator/helpers/generatePassword';
type usePasswordType = {
    password: string;
    setPassword: (password: string) => void;
    sliderVal: number;
    setSliderVal: (sliderVal: number) => void;
    includeNum: boolean;
    setIncludeNum: (includeNum: boolean) => void;
    includeSymbols: boolean;
    setIncludeSymbols: (includeSymbols: boolean) => void;
    isUserTyping: boolean;
    setIsUserTyping: (isUserTyping: boolean) => void;
};

export default function usePassword(): usePasswordType {
    const [password, setPassword] = useState('');
    const [sliderVal, setSliderVal] = useState(16);
    const [includeNum, setIncludeNum] = useState(true);
    const [includeSymbols, setIncludeSymbols] = useState(true);
    const [isUserTyping, setIsUserTyping] = useState(false);

    useEffect(() => {
        if (!isUserTyping) {
            setPassword(generatePassword(includeNum, includeSymbols, sliderVal));
        }
    }, [sliderVal, includeNum, includeSymbols]);

    return {
        password,
        setPassword,
        sliderVal,
        setSliderVal,
        includeNum,
        setIncludeNum,
        includeSymbols,
        setIncludeSymbols,
        isUserTyping,
        setIsUserTyping,
    };
}
