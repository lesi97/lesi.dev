import "./password-generator.scss";
import "./password-generator-mobile.scss";
import { useState, useEffect } from "react";
import { RefreshSvg, Copy } from "../../components/icons";
import { Checkbox } from "../../components/checkbox/checkbox";

export const PasswordGenerator = () => {
    const [password, setPassword] = useState("");
    const [sliderVal, setSliderVal] = useState(16);
    const [includeNum, setIncludeNum] = useState(true);
    const [includeSymbols, setIncludeSymbols] = useState(true);
    const [isUserTyping, setIsUserTyping] = useState(false);

    useEffect(() => {
        document.title = `Lesi | Password Generator`;
    }, []);

    useEffect(() => {
        if (!isUserTyping) {
            generatePassword();
        }
    }, [sliderVal, includeNum, includeSymbols]);


    function generatePassword() {
        const numbers = "0123456789";
        const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
        const symbols = "[!@#$£%^&*()+=-.,]";

        let chars = letters;
        if (includeNum) chars += numbers;
        if (includeSymbols) chars += symbols;

        const array = new Uint32Array(sliderVal);
        window.crypto.getRandomValues(array);

        let tempPassword = "";

        for (let i = 0; i < sliderVal; i++) {
            tempPassword += chars[array[i] % chars.length];
        }

        if (includeNum && !/\d/.test(tempPassword)) {
            const randomIndex = Math.floor(Math.random() * tempPassword.length);
            const randomNumber = numbers[Math.floor(Math.random() * numbers.length)];
            tempPassword = tempPassword.substring(0, randomIndex) + randomNumber + tempPassword.substring(randomIndex + 1);
        }

        if (includeSymbols && !/!@#$£%^&*()+=-.,/.test(tempPassword)) {
            const randomIndex = Math.floor(Math.random() * tempPassword.length);
            const randomSymbol = symbols[Math.floor(Math.random() * symbols.length)];
            tempPassword = tempPassword.substring(0, randomIndex) + randomSymbol + tempPassword.substring(randomIndex + 1);
        }

        setPassword(tempPassword);
    }

    function handleTyping(e) {
        setIsUserTyping(true);
        setPassword(e.target.value);
        setSliderVal(e.target.value.length);
    }

    function handleSliderChange(e) {
        setSliderVal(e.target.value);
        setIsUserTyping(false);
    }

    function copyPassword() {
        const copyText = document.getElementById("password");
        copyText.select();
        document.execCommand("copy");
    }

    const handleIncludeSymbolsChange = (e) => {
        setIncludeSymbols(e.target.checked);
    };

    const handleIncludeNumbersChange = (e) => {
        setIncludeNum(e.target.checked);
    };



    return (
        <main>
            <div className="passwordGenerator">
                <div className="description">
                    <h1>Password Generator</h1>
                    <h2>
                        Create a random password, adjust the slider to increase the password length
                    </h2>
                </div>

                <div className="passwordField">
                    <input type="text" value={password} id="password" onChange={e => handleTyping(e)}></input>
                    <div className="reloadSvg" onClick={generatePassword}>
                        <RefreshSvg />
                    </div>
                    <div className="copySvg" onClick={copyPassword}>
                        <Copy />
                    </div>
                </div>

                <div className="options">

                    <div className="sliderSection">
                        <p>Password Length: {sliderVal}</p>
                        <input
                            className="slider"
                            type="range"
                            min="8"
                            max="128"
                            value={sliderVal}
                            onChange={e => handleSliderChange(e)}></input>
                    </div>
                    <div className="additionalOptions">
                        <div>
                            <label>
                                <Checkbox checked={includeNum} onChange={handleIncludeNumbersChange} />
                                {/* <input
                                    type="checkbox"
                                    checked={includeNum}
                                    onChange={e => setIncludeNum(e.target.checked)}></input> */}
                                <span>Include Numbers</span>
                            </label>
                        </div>
                        <div>
                            <label>
                                <Checkbox checked={includeSymbols} onChange={handleIncludeSymbolsChange} />
                                {/* <input
                                    type="checkbox"
                                    checked={includeSymbols}
                                    onChange={e => setIncludeSymbols(e.target.checked)}></input> */}
                                <span>Include Symbols</span>
                            </label>
                        </div>



                        {/* <div>
                            <button onClick={copyPassword}>COPY TO CLIPBOARD</button>
                        </div> */}
                    </div>
                </div>


            </div>
        </main>
    )
}