import { DayNightToggle } from "../../components/dayNightToggle";
import "./settings.scss";
import { useEffect } from "react";
import { Link } from "react-router-dom";

export const Settings = ({ toggleNightMode, nigthModeState, setError }) => {

    useEffect(() => {
        document.title = `Lesi | Settings`;
        setError(false);
    }, []);

    return (
        <main>
            <div className="settings">
                <div className="settingsOptions"
                    role="button"
                    tabIndex="0"
                    onClick={toggleNightMode}
                    onKeyDown={(e) => {
                        if (e.code === "Enter" || e.code === "Space") toggleNightMode();
                    }}>
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
                <p>
                    <Link to="/privacy"
                        tabIndex={0}>
                        Privacy Policy
                    </Link>
                </p>
            </div>
        </main>
    )
}