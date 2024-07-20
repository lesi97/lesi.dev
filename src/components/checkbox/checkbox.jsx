import "./checkbox.scss";

export const Checkbox = ({ checked, onChange }) => {
    return (
        // eslint-disable-next-line jsx-a11y/label-has-associated-control
        <label className="checkBoxContainer">
            <input
                type="checkbox"
                checked={checked}
                onChange={onChange}
                tabIndex="0"></input>
            <span className="checkmark"></span>
        </label>
    );
};
