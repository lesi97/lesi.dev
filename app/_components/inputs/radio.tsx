interface RadioProps {
    name: string;
    value: string;
    checked: boolean;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export default function RadioInput({ name, value, checked, onChange }: RadioProps) {
    return (
        <input
            type="radio"
            name={name}
            value={value}
            checked={checked}
            onChange={onChange}
            className="appearance-none rounded-full w-4 h-4 border border-gray-300 absolute checked:bg-primary"
        ></input>
    );
}
