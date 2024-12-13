'use client';

// import './aspect-ratio.scss';
// import './aspect-ratio-mobile.scss';
import { useState, useEffect } from 'react';
import Description from '@/app/_components/description';
import RadioInput from '@/app/_components/inputs/radio';
import NumberInput from '@/app/_components/inputs/number';

export default function AspectRatio() {
    const [originalWidth, setOriginalWidth] = useState<number | ''>('');
    const [originalHeight, setOriginalHeight] = useState<number | ''>('');
    const [newWidth, setNewWidth] = useState<number | ''>('');
    const [newHeight, setNewHeight] = useState<number | ''>('');

    const [aspectRatio, setAspectRatio] = useState(0);
    const [selectedRadio, setSelectedRadio] = useState('width');

    useEffect(() => {
        if (!originalWidth || !originalHeight) return;
        const newAspectRatio = originalWidth / originalHeight;
        setAspectRatio(newAspectRatio);
        getSelectedRadioValue() === 'width' ? calculateNewHeight(newWidth) : calculateNewWidth(newHeight);
    }, [originalWidth, originalHeight]);

    useEffect(() => {
        if (!aspectRatio) return;
        if (getSelectedRadioValue() === 'width' && newWidth) {
            calculateNewHeight(newWidth);
        } else if (getSelectedRadioValue() === 'height' && newHeight) {
            calculateNewWidth(newHeight);
        }
    }, [aspectRatio, newWidth, newHeight]);

    function getSelectedRadioValue() {
        if (originalWidth === 0 || originalHeight === 0) return;
        const radios = Array.from(document.getElementsByName('keepValue'));
        const selectedRadio = radios.filter((radio) => radio.checked);
        return selectedRadio.length ? selectedRadio[0].value : 'width';
    }

    function calculateNewWidth(value) {
        if (parseInt(value) === 0) return;
        setNewHeight(parseInt(value));
        if (!newHeight) return;
        if (!aspectRatio) return;
        const width = Math.round(parseInt(value) * aspectRatio);
        setNewWidth(width);
    }

    function calculateNewHeight(value: number) {
        if (value === 0) return;
        setNewWidth(value);
        if (!newWidth || !aspectRatio) return;
        const height = Math.round(value / aspectRatio);
        setNewHeight(height);
    }

    const exceptThisSymbols = ['e', 'E', '+', '-'];

    return (
        <>
            <Description
                title="Aspect Ratio Calculator"
                subtitle="Input your old width and height then input your new width or height to automatically calculate the former or the latter"
            />

            <form className="grid grid-cols-2 gap-6">
                <label className="flex flex-col gap-2 text-center user-select-none">
                    Original Width
                    <NumberInput
                        value={originalWidth}
                        onChange={(e) => {
                            setOriginalWidth(parseInt(e.target.value));
                        }}
                    />
                </label>

                <label className="flex flex-col gap-2 text-center user-select-none">
                    Original Height
                    <NumberInput
                        value={originalHeight}
                        onChange={(e) => {
                            setOriginalHeight(parseInt(e.target.value));
                        }}
                    />
                </label>

                <label className="flex flex-col gap-2 relative text-center user-select-none">
                    New Width
                    {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                    <label htmlFor="width" className="absolute right-8 top-1 bg-primary">
                        <RadioInput
                            name="keepValue"
                            value="width"
                            checked={selectedRadio === 'width'}
                            onChange={() => setSelectedRadio('width')}
                        />
                    </label>
                    <NumberInput
                        value={newWidth}
                        onChange={(e) => {
                            calculateNewHeight(parseInt(e.target.value));
                        }}
                    />
                </label>

                <label className="flex flex-col gap-2 relative text-center user-select-none">
                    New Height
                    {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                    <label htmlFor="height" className="absolute right-8 top-1 bg-primary">
                        <RadioInput
                            name="keepValue"
                            value="height"
                            checked={selectedRadio === 'height'}
                            onChange={() => setSelectedRadio('height')}
                        />
                    </label>
                    <NumberInput
                        value={newHeight}
                        onChange={(e) => {
                            calculateNewWidth(parseInt(e.target.value));
                        }}
                    />
                </label>
            </form>
        </>
    );
}
