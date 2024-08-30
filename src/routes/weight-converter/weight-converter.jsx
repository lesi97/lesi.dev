import "./weight-converter.scss";
import "./weight-converter-mobile.scss";
import { useEffect, useState } from "react";

export const WeightConverter = ({ setError }) => {
    const [pounds, setPounds] = useState("");
    const [kilograms, setKilograms] = useState("");
    const [ounces, setOunces] = useState("");
    const [grams, setGrams] = useState("");
    const [stones, setStones] = useState("");

    useEffect(() => {
        document.title = `Lesi | Weight Converter`;
        setError(false);
    }, []);

    function sourcePounds(e) {
        setPounds(e.target.value);
        setKilograms((e.target.value / 2.2046).toFixed(4));
        setOunces((e.target.value * 16).toFixed(0));
        setGrams((e.target.value / 0.0022046).toFixed(4));
        setStones((e.target.value * 0.071429).toFixed(4));
    }

    function sourceKilograms(e) {
        setKilograms(e.target.value);
        setPounds((e.target.value * 2.2046).toFixed(4));
        setOunces((e.target.value * 35.274).toFixed(4));
        setGrams((e.target.value * 1000).toFixed(0));
        setStones((e.target.value * 0.1574).toFixed(4));
    }

    function sourceOunces(e) {
        setOunces(e.target.value);
        setPounds((e.target.value * 0.062500).toFixed(4));
        setKilograms((e.target.value / 35.274).toFixed(4));
        setGrams((e.target.value / 0.035274).toFixed(4));
        setStones((e.target.value * 0.0044643).toFixed(4));
    }

    function sourceGrams(e) {
        setGrams(e.target.value);
        setPounds((e.target.value * 0.0022046).toFixed(4));
        setKilograms((e.target.value / 1000).toFixed(2));
        setOunces((e.target.value * 0.035274).toFixed(4));
        setStones((e.target.value * 0.00015747).toFixed(4));
    }

    function sourceStones(e) {
        setStones(e.target.value);
        setPounds((e.target.value * 14).toFixed(0));
        setKilograms((e.target.value / 0.15747).toFixed(4));
        setOunces((e.target.value * 224).toFixed(0));
        setGrams((e.target.value / 0.00015747).toFixed(4));
    }

    const exceptThisSymbols = ["e", "E", "+", "-"];

    return (
        <main>
            <div className="weightConverter">
                <div className="description">
                    <h1>Weight Converter</h1>
                    <h2>
                        Type a value in any of the fields to convert between weight measurements<br />&nbsp;
                    </h2>
                </div>

                <div className="weights">

                    <label>
                        Kilograms (kg)
                        <input
                            type="number"
                            value={kilograms}
                            onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                            onChange={sourceKilograms}></input>
                    </label>

                    <label>
                        Pounds (lbs)
                        <input
                            type="number"
                            value={pounds}
                            onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                            onChange={sourcePounds}></input>
                    </label>

                    <label>
                        Ounces (oz)
                        <input
                            type="number"
                            value={ounces}
                            onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                            onChange={sourceOunces}></input>
                    </label>

                    <label>
                        Grams (g)
                        <input
                            type="number"
                            value={grams}
                            onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                            onChange={sourceGrams}></input>
                    </label>

                    <label>
                        Stones (st.)
                        <input
                            type="number"
                            value={stones}
                            onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                            onChange={sourceStones}></input>
                    </label>
                </div>
            </div>
        </main>
    )
}