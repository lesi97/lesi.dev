import "./aspect-ratio.scss";
import "./aspect-ratio-mobile.scss";
import { useState, useEffect } from "react";

export const AspectRatio = () => {
    const [originalWidth, setOriginalWidth] = useState("");
    const [originalHeight, setOriginalHeight] = useState("");
    const [newWidth, setNewWidth] = useState("");
    const [newHeight, setNewHeight] = useState("");
    const [aspectRatio, setAspectRatio] = useState("");
    const [selectedRadio, setSelectedRadio] = useState("width");

    useEffect(() => {
        document.title = `Lesi | Aspect Ratio Calculator`;
    }, []);

    useEffect(() => {
        if (!originalWidth || !originalHeight) return;
        const newAspectRatio = parseInt(originalWidth) / parseInt(originalHeight);
        setAspectRatio(newAspectRatio);
        getSelectedRadioValue() === "width" ? calculateNewHeight(newWidth) : calculateNewWidth(newHeight);
    }, [originalWidth, originalHeight]);

    useEffect(() => {
        if (!aspectRatio) return;
        if (getSelectedRadioValue() === "width" && newWidth) {
            calculateNewHeight(newWidth);
        } else if (getSelectedRadioValue() === "height" && newHeight) {
            calculateNewWidth(newHeight);
        }
    }, [aspectRatio, newWidth, newHeight]);

    function getSelectedRadioValue() {
        if (parseInt(originalWidth) === 0 || parseInt(originalHeight) === 0) return;
        const radios = Array.from(document.getElementsByName("keepValue"));
        const selectedRadio = radios.filter(radio => radio.checked);
        return selectedRadio.length ? selectedRadio[0].value : "width";
    };

    function calculateNewWidth(value) {
        if (parseInt(value) === 0) return;
        setNewHeight(parseInt(value));
        if (!aspectRatio) return;
        const width = Math.round(parseInt(value) * aspectRatio);
        setNewWidth(width);
    }

    function calculateNewHeight(value) {
        if (parseInt(value) === 0) return;
        setNewWidth(parseInt(value))
        if (!aspectRatio) return;
        const height = Math.round(parseInt(value) / aspectRatio);
        setNewHeight(height);
    }

    const exceptThisSymbols = ["e", "E", "+", "-"];


    return (
        <main>
            <div className="aspectRatio">
                <div className="description">
                    <h1>Aspect Ratio Calculator</h1>
                    <h2>
                        Input your old width and height<br />
                        Then input your new width or height to automatically calculate the former or the latter
                    </h2>
                    <h2>
                        You can also select to keep the new width or height value if you change the original values
                    </h2>
                </div>

                <form className="aspectRatioForm">
                    <div>
                        <label>
                            Original Width
                            <input type="number"
                                value={originalWidth}
                                onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                                onChange={(e) => { setOriginalWidth(e.target.value) }} />
                        </label>
                    </div>

                    <div>
                        <label>
                            Original Height
                            <input type="number"
                                value={originalHeight}
                                onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                                onChange={(e) => { setOriginalHeight(e.target.value) }} />
                        </label>
                    </div>

                    <div>
                        <label>
                            New Width
                            <label htmlFor="width" className="customRadio">
                                <input type="radio"
                                    id="width"
                                    name="keepValue"
                                    value="width"
                                    tabIndex="-1"
                                    checked={selectedRadio === "width"}
                                    onChange={() => setSelectedRadio("width")} />
                            </label>

                            <input type="number"
                                value={newWidth}
                                onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                                onChange={(e) => { calculateNewHeight(e.target.value) }} />
                        </label>
                    </div>

                    <div>
                        <label>
                            New Height
                            <label htmlFor="height"
                                className="customRadio">
                                <input type="radio"
                                    id="height"
                                    name="keepValue"
                                    value="height"
                                    tabIndex="-1"
                                    checked={selectedRadio === "height"}
                                    onChange={() => setSelectedRadio("height")} />
                            </label>
                            <input type="number"
                                value={newHeight}
                                onKeyDown={e => exceptThisSymbols.includes(e.key) && e.preventDefault()}
                                onChange={(e) => { calculateNewWidth(e.target.value) }} />
                        </label>
                    </div>
                </form >
            </div>
        </main>
    )
}