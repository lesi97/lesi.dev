import { countdownPlaceholders } from '@/lib';
import { useState, useEffect, useRef } from 'react';
import {
    type CountdownDataKeysType,
    type CountdownSchemaType,
    CountdownSchema,
    CountdownForm,
    CountdownInitialState,
} from '@/schema/countdownSchema';
import { ZodError } from 'zod';
import { parseError } from '@/utils';

export function useCountdown() {
    const [today, setToday] = useState(new Date());
    const [form, setForm] = useState(CountdownInitialState);
    const [editCommand, setEditCommand] = useState(false);
    const [commandName, setCommandName] = useState('countdown');
    const [command, setCommand] = useState<string | undefined>();
    const [placeholders, setPlaceholders] = useState({ message: '', fallback_message: '' });
    const commandRef = useRef(null);

    useEffect(() => {
        const sevenDaysLater = new Date(today);
        sevenDaysLater.setDate(sevenDaysLater.getDate() + 7);

        const year = sevenDaysLater.getFullYear();
        const month = String(sevenDaysLater.getMonth() + 1).padStart(2, '0');
        const day = String(sevenDaysLater.getDate()).padStart(2, '0');
        const hours = String(sevenDaysLater.getHours()).padStart(2, '0');
        const minutes = String(sevenDaysLater.getMinutes()).padStart(2, '0');

        const formattedDateTime = `${year}-${month}-${day}T${hours}:${minutes}`;
        setForm({ ...form, data: { ...form.data, target_date: formattedDateTime } });

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

    function updateState(key: CountdownDataKeysType, value: string | number, section: 'data' | 'errors' = 'data') {
        const newValue = section === 'data' ? value : { type: 'unique', value: value };
        setForm((prev) => {
            const updatedForm = {
                ...prev,
                [section]: {
                    ...prev[section],
                    [key]: newValue,
                },
            };

            if (section === 'data') {
                updatedForm.errors[key] = { type: 'unique', value: undefined };
            }

            return updatedForm;
        });
    }

    async function handleSubmit(e: React.FormEvent, form: typeof CountdownInitialState.data) {
        e.preventDefault();
        if (command) {
            return;
        }
        const localDate = new Date(form.target_date);
        const body = { ...form, target_date: localDate.toISOString() };
        try {
            CountdownSchema.parse(form);
            const response = await fetch('/api/v1/countdown', {
                method: 'POST',
                headers: { Accept: 'Application/Json' },
                body: JSON.stringify(body),
            });
            if (!response.ok) {
                return;
            }
            const uuid = await response.text();
            const url = `https://lesi.dev/api/countdown/${uuid}`;
            setCommandName(url);
        } catch (error) {
            if (error instanceof ZodError) {
                const parsed = parseError(error);
                Object.entries(parsed).forEach(([key, value]) => {
                    updateState(key as CountdownDataKeysType, value, 'errors');
                });
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
        transformCountdown,
        countdownPlaceholders,
        CountdownSchema,
        form,
        updateState,
        commandName,
        placeholders,
        handleSubmit,
        setEditCommand,
        editCommand,
        setCommandName,
    };
}
