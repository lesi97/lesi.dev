import "./toggle.scss";

const Toggle = ({ isChecked }) => {
    return (
        <input type="checkbox"
            className="customToggle"
            checked={isChecked}
            readOnly={true} />
    )
}

export { Toggle }