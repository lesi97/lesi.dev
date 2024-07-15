import { DayNightToggle } from "../../components/dayNightToggle";
import "./settings.scss";
import { useEffect } from "react";

export const Settings = ({ toggleNightMode, nigthModeState, setError }) => {

    setError(false);

    useEffect(() => {
        document.title = `Lesi | Settings`;
    }, []);

    return (
        <main>
            <div className="settings">
                <div className="settingsOptions" onClick={toggleNightMode}>
                    <h3>Night Mode</h3> <DayNightToggle isChecked={nigthModeState} />
                </div>
                <p>
                    Emote by <a href="https://x.com/Tenrinmaru1" target="_blank" rel="noreferrer" >Tenrinmaru</a>
                </p>
                <p>
                    Illustrations by <a href="https://undraw.co" target="_blank" rel="noreferrer" >undraw.co</a>
                </p>
                <p>
                    Background by <a href="https://heropatterns.com" target="blank" rel="noreferrer" >Hero Patterns</a>
                </p>
            </div>
        </main>
    )
}