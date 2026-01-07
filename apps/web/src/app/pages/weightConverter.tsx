import { ReactNode } from 'react';
import { Description, Input } from '@/components/ui';
import useWeightConverter from '@/hooks/useWeightConverter';
import { usePageMeta } from '@/hooks';

function Label({ children, htmlFor }: { children: ReactNode; htmlFor: string }) {
    return (
        <label htmlFor={htmlFor} className='grid w-full grid-cols-2 content-center'>
            {children}
        </label>
    );
}

export function WeightConverter() {
    usePageMeta({
        title: 'Weight Converter | Lesi',
        description: 'Convert between weight measurements',
    })
    const {
        pounds,
        updateWeightSourcePounds,
        kilograms,
        updateWeightSourceKilograms,
        ounces,
        updateWeightSourceOunces,
        grams,
        updateWeightSourceGrams,
        stones,
        updateWeightSourceStones,
    } = useWeightConverter();

    return (
        <>
            <Description
                title='Weight Converter'
                subtitle='Type a value in any of the fields to convert between weight measurements'
            />
            <div className='flex flex-col gap-y-4 lg:px-20'>
                <div className='flex flex-col justify-between gap-4'>
                    <Label htmlFor='kilograms'>
                        <span>Kilograms (kg)</span>
                        <Input
                            value={kilograms}
                            onChange={updateWeightSourceKilograms}
                            id='kilograms'
                            type='number'
                            className='w-full !text-left'
                            variant='outline'
                        />
                    </Label>

                    <Label htmlFor='pounds'>
                        <span>Pounds (lbs)</span>
                        <Input
                            value={pounds}
                            onChange={updateWeightSourcePounds}
                            id='pounds'
                            type='number'
                            variant='outline'
                            className='w-full !text-left'
                        />
                    </Label>
                    <Label htmlFor='ounces'>
                        <span>Ounces (oz)</span>
                        <Input
                            value={ounces}
                            onChange={updateWeightSourceOunces}
                            id='ounces'
                            type='number'
                            variant='outline'
                            className='w-full !text-left'
                        />
                    </Label>
                    <Label htmlFor='grams'>
                        <span>Grams (g)</span>
                        <Input
                            value={grams}
                            onChange={updateWeightSourceGrams}
                            id='grams'
                            type='number'
                            variant='outline'
                            className='w-full !text-left'
                        />
                    </Label>
                    <Label htmlFor='stones'>
                        <span>Stones (st)</span>
                        <Input
                            value={stones}
                            onChange={updateWeightSourceStones}
                            id='stones'
                            type='number'
                            variant='outline'
                            className='w-full !text-left'
                        />
                    </Label>
                </div>
            </div>
        </>
    );
}
