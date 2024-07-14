import "./checkbox.scss";

export const Checkbox = ({ checked, onChange }) => {
    return (
        <label className="checkBoxContainer">
            <input
                type="checkbox"
                checked={checked}
                onChange={onChange}></input>
            <span className="checkmark"></span>
        </label>
    );
};
