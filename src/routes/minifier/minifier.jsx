import { useEffect, useState } from "react";
import "./minifier.scss";
import "./minifier-mobile.scss";
import { minify, unminify } from "./minifier-fns";

export const Minifier = ({ setError }) => {
    const [code, setCode] = useState("");

    useEffect(() => {
        document.title = `Lesi | Minifier`;
        setError(false);
    }, []);

    function handleMinify() {
        const newCode = minify(code);
        if (newCode) {
            setCode(newCode);
        }
    }

    function handleUnminify() {
        const newCode = unminify(code);
        if (newCode) {
            setCode(newCode);
        }
    }

    function copyCode() {
        const copyText = document.getElementById("code");
        copyText.select();
        document.execCommand("copy");
    }

    return (
        <main>
            <div className="minifier">
                <div className="description">
                    <h1>Minifier</h1>
                    <h2>
                        Minify or unminfiy CSS, JS, XML &amp; JSON
                        <br />
                        &nbsp;
                    </h2>
                </div>
                <textarea
                    id="code"
                    className="codeArea"
                    rows="10"
                    spellCheck="false"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                ></textarea>
                <div className="buttons">
                    <button onClick={handleMinify}>MINIFY</button>
                    <button onClick={handleUnminify}>UNMINIFY</button>
                    <button onClick={copyCode}>COPY</button>
                </div>
            </div>
        </main>
    );
};
