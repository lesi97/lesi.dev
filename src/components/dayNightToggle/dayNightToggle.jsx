import "./dayNightToggle.scss";
import { Alien } from "../icons";
const DayNightToggle = ({ isChecked, toggle }) => {
  return (
    <>
      <input
        type="checkbox"
        id="toggle"
        className="toggle--checkbox"
        checked={isChecked}
        readOnly={true}
      />

      <label className="toggle--label">
        <span className="toggle--label-background"><Alien /></span>
      </label>
    </>
  );
};

export { DayNightToggle };
