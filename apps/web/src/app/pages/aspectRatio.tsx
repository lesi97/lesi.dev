import { Description, Input, Radio } from '@/components/ui';
import { KeyboardEvent } from 'react';
import { useAspectRatio } from '@/hooks';

export function AspectRatio() {
    const {
        originalWidth,
        originalHeight,
        newWidth,
        newHeight,
        selectedRadio,
        calculateNewWidth,
        calculateNewHeight,
        setOriginalWidth,
        setOriginalHeight,
        setSelectedRadio,
    } = useAspectRatio();

    function preventLetters(e: KeyboardEvent<HTMLInputElement>) {
        const chars = ['e', '+', '-'];
        if (chars.includes(e.key.toLocaleLowerCase())) {
            e.preventDefault();
        }
    }
    return (
        <>
            <Description
                title='Aspect Ratio Calculator'
                subtitle='Input your old width and height then input your new width or height to automatically calculate the former or the latter'
            />
            <form className='flex flex-col gap-6 md:grid md:grid-cols-2'>
                <label className='user-select-none flex flex-col gap-2 text-center'>
                    Original Width
                    <Input
                        id='original-width'
                        value={originalWidth ?? ''}
                        onKeyDown={preventLetters}
                        onChange={(e) => {
                            const newVal = parseInt(e.target.value) || 0;
                            setOriginalWidth(newVal);
                        }}
                        data-testid='originalWidth'
                        variant='outline'
                        type='number'
                    />
                </label>

                <label className='user-select-none flex flex-col gap-2 text-center'>
                    Original Height
                    <Input
                        id='original-height'
                        value={originalHeight ?? ''}
                        onKeyDown={preventLetters}
                        onChange={(e) => {
                            const newVal = parseInt(e.target.value) || 0;
                            setOriginalHeight(newVal);
                        }}
                        data-testid='originalHeight'
                        variant='outline'
                        type='number'
                    />
                </label>

                <label className='user-select-none relative flex flex-col gap-2 text-center'>
                    New Width
                    {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                    <Radio
                        name='keepValue'
                        checked={selectedRadio === 'width'}
                        size='circle'
                        variant='circle-primary'
                        onChange={() => setSelectedRadio('width')}
                        className='absolute right-4 top-0'
                        value='width'
                    />
                    <Input
                        id='new-width'
                        value={newWidth ?? ''}
                        onKeyDown={preventLetters}
                        onChange={(e) => {
                            const newVal = e.target.value || '0';
                            calculateNewHeight(newVal);
                        }}
                        data-testid='newWidth'
                        variant='outline'
                        type='number'
                    />
                </label>

                <label className='user-select-none relative flex flex-col gap-2 text-center'>
                    New Height
                    {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                    <Radio
                        name='keepValue'
                        size='circle'
                        variant='circle-primary'
                        checked={selectedRadio === 'height'}
                        onChange={() => setSelectedRadio('height')}
                        className='absolute right-4 top-0'
                        value='height'
                    />
                    <Input
                        id='new-height'
                        value={newHeight ?? ''}
                        onKeyDown={preventLetters}
                        onChange={(e) => {
                            const newVal = e.target.value || '0';
                            calculateNewWidth(newVal);
                        }}
                        data-testid='newHeight'
                        variant='outline'
                        type='number'
                    />
                </label>
            </form>
        </>
    );
}
