import { useEffect, useState, useRef } from 'react';

export function useTime() {
    const [time, setTime] = useState<string>('');
    const [serverTime, setServerTime] = useState<Date | null>(null);
    const [intervalId, setIntervalId] = useState<number | null>(null);
    const [date, setDate] = useState<string | null>(null);

    function formatDate(d: Date): string {
        const day = d.getDate().toString().padStart(2, '0');
        const months = ['JAN', 'FEB', 'MAR', 'APR', 'MAY', 'JUN', 'JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'];
        const month = months[d.getMonth()];
        const year = d.getFullYear().toString();
        return `${day} ${month} ${year}`;
    }

    async function fetchInitialTime() {
        try {const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
        const res = await fetch(`/api/time?zone=${timeZone}`);
        const data = await res.json();
        const [day, month, year] = data.message.date.split('/');
        const [hours, minutes, seconds] = data.message.time.split(':');
        const initialServerTime = new Date(year, month - 1, day, hours, minutes, seconds);
        setServerTime(initialServerTime);
        setTime(data.message.time);} catch (error) {
            console.error(error)
        }
        
    }

    useEffect(() => {
        if (!serverTime) {
            return;
        }
        setDate(formatDate(serverTime));
    }, [serverTime]);

    useEffect(() => {
        const id = setInterval(() => {
            setServerTime((prevTime) => {
                if (!prevTime) return null;
                const newTime = new Date(prevTime.getTime() + 1000);
                const hours = newTime.getHours().toString().padStart(2, '0');
                const minutes = newTime.getMinutes().toString().padStart(2, '0');
                const seconds = newTime.getSeconds().toString().padStart(2, '0');
                setTime(`${hours}:${minutes}:${seconds}`);

                if (prevTime.getDate() !== newTime.getDate()) {
                    updateDate(newTime);
                }

                return newTime;
            });
        }, 1000);
        setIntervalId(id);

        const updateDate = (date: Date) => {
            const day = date.getDate();
            const suffix =
                ['th', 'st', 'nd', 'rd'][(day % 10 > 3 ? 0 : day % 10) * (day < 10 || day > 20 ? 1 : 0)] || 'th';
            const formattedDate = date.toLocaleDateString('en-US', {
                weekday: 'long',
                day: 'numeric',
                month: 'long',
                year: 'numeric',
            });
            const [weekday, month, dayNum, year] = formattedDate.replace(',', '').split(' ');
            setDate(`${weekday} ${day}${suffix} ${month} ${year}`);
        };

        fetchInitialTime();

        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, []);

    return { time, date };
}
