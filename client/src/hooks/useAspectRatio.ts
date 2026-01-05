import { useState, useEffect } from 'react';

export function useAspectRatio() {
    const [originalWidth, setOriginalWidth] = useState<number | undefined>();
    const [originalHeight, setOriginalHeight] = useState<number | undefined>();
    const [newWidth, setNewWidth] = useState<number | undefined>();
    const [newHeight, setNewHeight] = useState<number | undefined>();

    const [aspectRatio, setAspectRatio] = useState<string>();
    const [selectedRadio, setSelectedRadio] = useState('width');

    useEffect(() => {
        if (!originalWidth || !originalHeight) return;
        const newAspectRatio = originalWidth / originalHeight;
        setAspectRatio(newAspectRatio.toString());
        getSelectedRadioValue() === 'width'
            ? calculateNewHeight(newWidth?.toString() ?? '0')
            : calculateNewWidth(newHeight?.toString() ?? '0');
    }, [originalWidth, originalHeight]);

    useEffect(() => {
        if (!aspectRatio) return;
        if (getSelectedRadioValue() === 'width' && newWidth) {
            calculateNewHeight(newWidth.toString());
        } else if (getSelectedRadioValue() === 'height' && newHeight) {
            calculateNewWidth(newHeight.toString());
        }
    }, [aspectRatio, newWidth, newHeight]);

    function getSelectedRadioValue() {
        if (originalWidth === 0 || originalHeight === 0) return;
        const radios: HTMLInputElement[] = Array.from(
            document.getElementsByName('keepValue') as NodeListOf<HTMLInputElement>
        );
        const selectedRadio = radios.filter((radio) => radio.checked);
        return selectedRadio.length ? selectedRadio[0].value : 'width';
    }

    function calculateNewWidth(value: string) {
        const newHeightValue = parseInt(value);
        if (newHeightValue === 0 || !aspectRatio) return;
        setNewHeight(newHeightValue);
        const computedWidth = Math.round(newHeightValue * parseFloat(aspectRatio));
        setNewWidth(computedWidth);
    }

    function calculateNewHeight(value: string) {
        const newWidthValue = parseInt(value);
        if (newWidthValue === 0 || !aspectRatio) return;
        setNewWidth(newWidthValue);
        const computedHeight = Math.round(newWidthValue / parseFloat(aspectRatio));
        setNewHeight(computedHeight);
    }
    return {
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
    };
}
