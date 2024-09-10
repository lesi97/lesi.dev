import "./dayNightToggle.scss";
import { Alien } from "../icons";
const DayNightToggle = ({ isChecked }) => {
    return (
        <>
            <input
                type="checkbox"
                id="toggle"
                className="toggle--checkbox"
                checked={isChecked}
                readOnly={true}
            />

            {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
            <label className="toggle--label">
                <span className="toggle--label-background">
                    <Alien />
                </span>
            </label>
        </>
    );
};

export { DayNightToggle };
