import { countdownPlaceholders } from '@/lib';
import { useState, useEffect, useRef, type FormEvent } from 'react';
import { countdownSchema } from '@/schema/countdownSchema';
import { ZodError } from 'zod';
import { parseError } from '@/utils';

const formIntialState = {
    data: {
        command: '',
        target_date: '',
        message: '',
        fallback_message: '',
    }, errors: {
        command: {type: 'unique', value: undefined},
        target_date: {type: 'unique', value: undefined},
        message: {type: 'unique', value: undefined},
        fallback_message: {type: 'unique', value: undefined},
    }
}

type InitialStateType = typeof formIntialState.data;
type CountdownDataKeysType = keyof InitialStateType;


type ErrorsType = {
    target_date: string | null;
    message: string | null;
    fallback_message:string | null;
};


export function useCountdown() {
const [today, setToday] = useState(new Date());
    const [data, setData] = useState({
        target_date: '',
        message: '',
        fallback_message: '',
    });
    const [errors, setErrors] = useState<ErrorsType>({
        target_date: null,
        message: null,
        fallback_message: null,
    });
    const [editCommand, setEditCommand] = useState(false);
    const [commandName, setCommandName] = useState('countdown');
    const [command, setCommand] = useState<string | undefined>();
    const [placeholders, setPlaceholders] = useState({ message: '', fallback_message: '' });
    const commandRef = useRef(null);
    const commandRef2 = useRef(null);

    useEffect(() => {
        const sevenDaysLater = new Date(today);
        sevenDaysLater.setDate(sevenDaysLater.getDate() + 7);

        const year = sevenDaysLater.getFullYear();
        const month = String(sevenDaysLater.getMonth() + 1).padStart(2, '0');
        const day = String(sevenDaysLater.getDate()).padStart(2, '0');
        const hours = String(sevenDaysLater.getHours()).padStart(2, '0');
        const minutes = String(sevenDaysLater.getMinutes()).padStart(2, '0');

        const formattedDateTime = `${year}-${month}-${day}T${hours}:${minutes}`;
        setData({ ...data, target_date: formattedDateTime });

        const random = Math.floor(Math.random() * countdownPlaceholders.length);
        const newPlaceholder = countdownPlaceholders[random];
        setPlaceholders(newPlaceholder);

        const interval = setInterval(() => {
            setToday(new Date());
        }, 1000);

        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        if (commandRef.current) {
            (commandRef.current as HTMLTextAreaElement).select();
            document.execCommand('copy');
            (commandRef.current as HTMLTextAreaElement).scrollIntoView({
                block: 'center',
                inline: 'center',
            });
        }
    }, [command]);

    function handleZodErrors(zodErrors: ZodError) {
        const newErrors: Record<string, string> = {};
        zodErrors.errors.forEach((error) => {
            const path = error.path[error.path.length - 1];
            newErrors[path] = error.message;
        });
        setErrors(newErrors as ErrorsType);
    }

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        if (command) {
            return;
        }
        const localDate = new Date(data.target_date);
        const body = { ...data, target_date: localDate.toISOString() };
        try {
            countdownSchema.parse(data);
            const response = await fetch('/api/v1/countdown', {
                method: 'POST',
                headers: { Accept: 'Application/Json' },
                body: JSON.stringify(body),
            });
            if (!response.ok) {
                return;
            }
            const uuid = await response.text();
            const url = `https://lesi.dev/api/v1/countdown/${uuid}`;
            setCommand(url);
        } catch (error) {
            if (error instanceof ZodError) {
                handleZodErrors(error);
            }
            console.error(error);
        }
    }


    function transformCountdown(targetDate: string) {
        const futureDate = new Date(targetDate);
        const currentDate = new Date();

        const diffMs = futureDate.getTime() - currentDate.getTime();

        if (diffMs < 0) {
            return 'Passed';
        }

        let diffSeconds = Math.floor(diffMs / 1000);
        const days = Math.floor(diffSeconds / (3600 * 24));
        diffSeconds %= 3600 * 24;
        const hours = Math.floor(diffSeconds / 3600);
        diffSeconds %= 3600;
        const minutes = Math.floor(diffSeconds / 60);
        const seconds = diffSeconds % 60;

        const countdownParts = [];
        if (days > 0) {
            countdownParts.push(`${days} days`);
        }
        countdownParts.push(`${hours} ${hours === 1 ? 'hour' : 'hours'}`);
        countdownParts.push(`${minutes} ${minutes === 1 ? 'minute' : 'minutes'}`);
        countdownParts.push(`${seconds} ${seconds === 1 ? 'second' : 'seconds'}`);

        return countdownParts.join(', ');
    }

    return {
        command,
        commandRef,
        commandRef2,
        setData,
        setErrors,
        transformCountdown,
        countdownPlaceholders,
        countdownSchema,
        data,
        errors,
        commandName,
        placeholders,
        handleSubmit,
        setEditCommand,
        editCommand,
        setCommandName,
    };
}
