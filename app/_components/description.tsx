interface DescriptionProps {
    title: string;
    subtitle: string | React.ReactNode;
}

const Description = ({ title, subtitle }: DescriptionProps) => {
    return (
        <div className="text-center mb-4">
            <h1 className="text-2xl mb-4">{title}</h1>
            <h2 className="text-lg">{subtitle}</h2>
        </div>
    );
};

export default Description;
