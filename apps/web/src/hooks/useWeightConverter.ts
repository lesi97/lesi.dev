import { useState } from 'react';

type WeightConverterFormType = {
    pounds: string | number | readonly string[] | undefined;
    updateWeightSourcePounds: (e: React.ChangeEvent<HTMLInputElement>) => void;
    kilograms: string | number | readonly string[] | undefined;
    updateWeightSourceKilograms: (e: React.ChangeEvent<HTMLInputElement>) => void;
    ounces: string | number | readonly string[] | undefined;
    updateWeightSourceOunces: (e: React.ChangeEvent<HTMLInputElement>) => void;
    grams: string | number | readonly string[] | undefined;
    updateWeightSourceGrams: (e: React.ChangeEvent<HTMLInputElement>) => void;
    stones: string | number | readonly string[] | undefined;
    updateWeightSourceStones: (e: React.ChangeEvent<HTMLInputElement>) => void;
};

export default function useWeightConverter(): WeightConverterFormType {
    const [pounds, setPounds] = useState<WeightConverterFormType['pounds']>();
    const [kilograms, setKilograms] = useState<WeightConverterFormType['kilograms']>();
    const [ounces, setOunces] = useState<WeightConverterFormType['ounces']>();
    const [grams, setGrams] = useState<WeightConverterFormType['grams']>();
    const [stones, setStones] = useState<WeightConverterFormType['stones']>();

    function updateWeightSourcePounds(e: React.ChangeEvent<HTMLInputElement>) {
        const value = parseFloat(e.target.value);
        setPounds(value);
        setKilograms(parseFloat((value / 2.2046).toFixed(4)));
        setOunces(parseFloat((value * 16).toFixed(0)));
        setGrams(parseFloat((value / 0.0022046).toFixed(4)));
        setStones(parseFloat((value * 0.071429).toFixed(4)));
    }

    function updateWeightSourceKilograms(e: React.ChangeEvent<HTMLInputElement>) {
        const value = parseFloat(e.target.value);
        setKilograms(value);
        setPounds(parseFloat((value * 2.2046).toFixed(4)));
        setOunces(parseFloat((value * 35.274).toFixed(4)));
        setGrams(parseFloat((value * 1000).toFixed(0)));
        setStones(parseFloat((value * 0.1574).toFixed(4)));
    }

    function updateWeightSourceOunces(e: React.ChangeEvent<HTMLInputElement>) {
        const value = parseFloat(e.target.value);
        setOunces(value);
        setPounds(parseFloat((value * 0.0625).toFixed(4)));
        setKilograms(parseFloat((value / 35.274).toFixed(4)));
        setGrams(parseFloat((value / 0.035274).toFixed(4)));
        setStones(parseFloat((value * 0.0044643).toFixed(4)));
    }

    function updateWeightSourceGrams(e: React.ChangeEvent<HTMLInputElement>) {
        const value = parseFloat(e.target.value);
        setGrams(value);
        setPounds(parseFloat((value * 0.0022046).toFixed(4)));
        setKilograms(parseFloat((value / 1000).toFixed(2)));
        setOunces(parseFloat((value * 0.035274).toFixed(4)));
        setStones(parseFloat((value * 0.00015747).toFixed(4)));
    }

    function updateWeightSourceStones(e: React.ChangeEvent<HTMLInputElement>) {
        const value = parseFloat(e.target.value);
        setStones(value);
        setPounds(parseFloat((value * 14).toFixed(0)));
        setKilograms(parseFloat((value / 0.15747).toFixed(4)));
        setOunces(parseFloat((value * 224).toFixed(0)));
        setGrams(parseFloat((value / 0.00015747).toFixed(4)));
    }

    return {
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
    };
}
