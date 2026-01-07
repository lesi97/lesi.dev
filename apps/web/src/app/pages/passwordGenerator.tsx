import { Description } from '@/components/ui';
import { useRef } from 'react';
import { usePassword } from '@/hooks/usePassword';
import { Input, Icons, Slider, Checkbox } from '@/components/ui';
import { usePageMeta } from '@/hooks';

export function PasswordGenerator() {
    usePageMeta({
        title: 'Password Generator | Lesi',
        description: 'Create a random password',
    })
    const ref = useRef<HTMLInputElement>(null);
    const {
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
    } = usePassword();

    function handleTyping(e: React.ChangeEvent<HTMLInputElement>) {
        setIsUserTyping(true);
        setPassword(e.target.value);
        setSliderVal(e.target.value.length);
    }

    function handleSliderChange(value: number) {
        setSliderVal(value);
        setIsUserTyping(false);
    }

    function handleIncludeSymbolsChange(e: React.ChangeEvent<HTMLInputElement>) {
        setIncludeSymbols(e.target.checked);
    }

    function handleIncludeNumbersChange(e: React.ChangeEvent<HTMLInputElement>) {
        setIncludeNum(e.target.checked);
    }

    return (
        <>
            <Description
                title='Password Generator'
                subtitle='Create a random password, adjust the slider to increase the password length'
            />
            <div className='passwordField relative w-full pb-8'>
                <Input
                    id='password'
                    value={password}
                    onChange={(e) => handleTyping(e)}
                    spellCheck={false}
                    type='text'
                    variant='underline'
                    className='w-full px-20 text-center'
                    ref={ref}
                />

                <div
                    className='absolute right-10 top-2 focus-within:!outline-0 focus-within:!ring-0 focus:!outline-0 focus:!ring-0 focus-visible:!outline-0 focus-visible:!ring-1'
                    tabIndex={0}
                    role='button'
                    onClick={() => copyPassword(ref)}
                    onKeyDown={(e) => {
                        if (e.code === 'Enter') copyPassword(ref);
                    }}>
                    <Icons.Copy width={20} height={20} />
                </div>
                <div
                    className='absolute right-2 top-2 focus-within:!outline-0 focus-within:!ring-0 focus:!outline-0 focus:!ring-0 focus-visible:!outline-0 focus-visible:!ring-1'
                    tabIndex={0}
                    role='button'
                    onClick={() => {
                        const newPassword = generatePassword(includeNum, includeSymbols, sliderVal);
                        setPassword(newPassword);
                        ref.current?.focus();
                    }}
                    onKeyDown={(e) => {
                        if (e.code === 'Enter') {
                            const newPassword = generatePassword(includeNum, includeSymbols, sliderVal);
                            setPassword(newPassword);
                            ref.current?.focus();
                        }
                    }}>
                    <Icons.Refresh width={20} height={20} />
                </div>
            </div>

            <div className='options flex flex-col gap-6 lg:flex-row lg:gap-10'>
                <div className='sliderSection flex w-full flex-col items-center justify-center gap-y-6 focus-visible:outline-none lg:w-3/5 lg:justify-between lg:pb-2'>
                    <p className='w-full text-center'>Password Length: {sliderVal}</p>
                    <Slider min='8' max='128' value={sliderVal} onChange={(value) => handleSliderChange(value)} />
                </div>

                <div className='flex h-full flex-row justify-between gap-4 lg:flex-col'>
                    <label className='flex flex-row items-center gap-2'>
                        <Checkbox
                            checked={includeNum}
                            onChange={handleIncludeNumbersChange}
                            size='sm'
                            variant='secondary'
                        />
                        <span>Include Numbers</span>
                    </label>
                    <label className='flex flex-row items-center gap-2'>
                        <Checkbox
                            checked={includeSymbols}
                            onChange={handleIncludeSymbolsChange}
                            variant='secondary'
                            size='sm'
                        />
                        <span>Include Symbols</span>
                    </label>
                </div>
            </div>
        </>
    );
}
