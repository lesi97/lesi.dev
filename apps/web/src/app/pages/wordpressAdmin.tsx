import { ChangeEvent, MouseEvent, useState } from 'react';
import { z, ZodError } from 'zod';

const fakeLoginSchema = z.object({
    username: z.string().trim().min(1, 'Username is required'),
    password: z.string().trim().min(1, 'Password is required'),
});

type FakeLoginSchemaType = z.infer<typeof fakeLoginSchema>;

export function WordpressAdmin() {
    const [loading, setLoading] = useState(false);
    const [data, setData] = useState<FakeLoginSchemaType>({ username: '', password: '' });
    const [errors, setErrors] = useState<FakeLoginSchemaType>({ username: '', password: '' });

    function onClick(e: MouseEvent<HTMLButtonElement>) {
        try {
            if (e.currentTarget.innerText === 'Log In') {
                fakeLoginSchema.parse(data);
            }
            setLoading(true);
            const rick = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
            setTimeout(() => {
                window.location.replace(rick);
            }, 4000);
        } catch (error) {
            if (error instanceof ZodError) {
                const errs = { username: '', password: '' };
                error.issues.forEach((err) => {
                    const path = err.path[0] as keyof FakeLoginSchemaType;
                    errs[path] = err.message;
                });
                setErrors(errs);
            }
        }
    }

    if (loading) {
        return (
            <>
                <style>
                    {`
                  .container {
                        --uib-size: 100px;
                        --uib-color: #2271b1;
                        --uib-speed: 2s;
                        --uib-bg-opacity: 0;
                        height: var(--uib-size);
                        width: var(--uib-size);
                        transform-origin: center;
                        animation: rotate var(--uib-speed) linear infinite;
                        will-change: transform;
                        overflow: visible;
                    }

                    .car {
                        fill: none;
                        stroke: var(--uib-color);
                        stroke-dasharray: 1, 200;
                        stroke-dashoffset: 0;
                        stroke-linecap: round;
                        animation: stretch calc(var(--uib-speed) * 0.75) ease-in-out infinite;
                        will-change: stroke-dasharray, stroke-dashoffset;
                        transition: stroke 0.5s ease;
                    }

                    .track {
                        fill: none;
                        stroke: var(--uib-color);
                        opacity: var(--uib-bg-opacity);
                        transition: stroke 0.5s ease;
                    }

                    @keyframes rotate {
                        100% {
                        transform: rotate(360deg);
                        }
                    }

                    @keyframes stretch {
                        0% {
                        stroke-dasharray: 0, 150;
                        stroke-dashoffset: 0;
                        }
                        50% {
                        stroke-dasharray: 75, 150;
                        stroke-dashoffset: -25;
                        }
                        100% {
                        stroke-dashoffset: -100;
                        }
                    }
                `}
                </style>
                <div className='bg-[#f0f0f1] w-full h-full flex-1 flex items-center justify-center'>
                    <div className='flex flex-col items-center gap-6'>
                        <svg className='container' viewBox='0 0 40 40' height='40' width='40'>
                            <circle
                                className='track'
                                cx='20'
                                cy='20'
                                r='17.5'
                                pathLength='100'
                                stroke-width='2px'
                                fill='none'
                            />
                            <circle
                                className='car'
                                cx='20'
                                cy='20'
                                r='17.5'
                                pathLength='100'
                                stroke-width='2px'
                                fill='none'
                            />
                        </svg>
                    </div>
                </div>
            </>
        );
    }

    function updateData(e: ChangeEvent<HTMLInputElement>) {
        const { id, value } = e.currentTarget;

        setData((prev) => {
            return { ...prev, [id]: value };
        });
    }

    return (
        <>
            <div className='bg-[#f0f0f1] w-full h-full flex-1 flex items-center justify-center'>
                <div className='flex flex-col items-center gap-6'>
                    <span>
                        <img src='/images/wordpress.png' width={64} />
                    </span>
                    <div className='bg-white h-fit py-8 px-6 border gap-4 flex flex-col rounded-sm border-[#abaeb1] text-[#767b80] text-sm min-w-80'>
                        <WordpressInput
                            inputLabel='Username or Email Address'
                            type='text'
                            id='username'
                            value={data.username}
                            error={errors.username}
                            onChange={updateData}
                        />
                        <WordpressInput
                            inputLabel='Password'
                            type='password'
                            id='password'
                            value={data.password}
                            error={errors.password}
                            onChange={updateData}
                        />
                        <div className='flex flex-row justify-between items-center'>
                            <label className='flex flex-row items-center justify-center gap-2'>
                                <input type='checkbox' className='accent-[#2271b1]' />
                                <span>Remember Me</span>
                            </label>
                            <button className='bg-[#2271b1] rounded py-1.5 px-3 text-white' onClick={onClick}>
                                Log In
                            </button>
                        </div>
                    </div>
                    <div className='flex place-content-start justify-start w-full text-[#767b80] text-sm pl-6 font-light flex-col gap-4'>
                        <button className='hover:underline w-fit' onClick={onClick}>
                            Lost your password?
                        </button>
                        <button className='hover:underline w-fit' onClick={onClick}>
                            ← Go to Lesi
                        </button>
                    </div>
                </div>
            </div>
        </>
    );
}

type WordpressInputProps = {
    inputLabel: string;
    type: string;
    id: string;
    value: string;
    error: string;
    onChange: (e: ChangeEvent<HTMLInputElement>) => void;
};

function WordpressInput(props: WordpressInputProps) {
    const { inputLabel, type, id, value, error, onChange } = props;
    return (
        <div>
            <label className='flex flex-col w-full' htmlFor={id}>
                <span>{inputLabel}</span>
                <input
                    value={value}
                    id={id}
                    className='h-6 p-4 rounded border border-[#abaeb1] focus:outline-[#2271b1]'
                    type={type}
                    onChange={onChange}
                />
            </label>
            {error !== '' ? <span className='text-xs text-error'>{error}</span> : null}
        </div>
    );
}
