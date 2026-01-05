import { Description } from '@/components/layout';
import { useAspectRatio } from '@/hooks';
import { Input, Radio } from '@/components/ui';
import { KeyboardEvent, useCallback } from 'react';
import { preventCharsInInput } from '@/utils';

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

    const disallowedChars = ['e', '+', '-'];

    return (
        <div className='flex w-11/12 flex-col gap-4'>
            <Description
                title='Aspect Ratio Calculator'
                subtitle='Input your old width and height then input your new width or height to automatically calculate the former or the latter'
            />
            <form className='flex flex-col gap-6 md:grid md:grid-cols-2'>
                <label className='user-select-none flex flex-col gap-2 text-center'>
                    Original Width
                    <Input
                        id='originalWidth'
                        value={originalWidth ?? ''}
                        onKeyDown={useCallback(
                            (e: KeyboardEvent<HTMLInputElement>) => preventCharsInInput(e, disallowedChars),
                            []
                        )}
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
                        id='originalHeight'
                        value={originalHeight ?? ''}
                        onKeyDown={useCallback(
                            (e: KeyboardEvent<HTMLInputElement>) => preventCharsInInput(e, disallowedChars),
                            []
                        )}
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
                        id='newWidth'
                        value={newWidth ?? ''}
                        onKeyDown={useCallback(
                            (e: KeyboardEvent<HTMLInputElement>) => preventCharsInInput(e, disallowedChars),
                            []
                        )}
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
                        id='newHeight'
                        value={newHeight ?? ''}
                        onKeyDown={useCallback(
                            (e: KeyboardEvent<HTMLInputElement>) => preventCharsInInput(e, disallowedChars),
                            []
                        )}
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
        </div>
    );
}
