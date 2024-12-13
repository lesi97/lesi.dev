interface NumberProps {
    value: number;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export default function NumberInput({ value, onChange }: NumberProps) {
    return (
        <input
            type="number"
            value={value}
            onChange={onChange}
            className="text-black !appearance-none cursor-text rounded bg-white h-8 w-full text-base px-4 py-2"
        ></input>
    );
}
