'use client'

import { selectSql } from '@/lib/server/data/database/helpers';
import Image from 'next/image';
import Description from '@/app/_components/description';
import { useState, useEffect } from 'react';

export default function Home() {
    const [time, setTime] = useState<string>('');
    const [serverTime, setServerTime] = useState<Date | null>(null);
    const [intervalId, setIntervalId] = useState<NodeJS.Timeout | null>(null);
    const [date, setDate] = useState<string | null>(null);

    useEffect(() => {
        const fetchInitialTime = async () => {
            const res = await fetch('/api/v1/time');
            const data = await res.json();
            const [day, month, year] = data.message.date.split('-');
            const [hours, minutes, seconds] = data.message.time.split(':');
            const initialServerTime = new Date(year, month - 1, day, hours, minutes, seconds);
            setServerTime(initialServerTime);
            setTime(data.message.time);

            const id = setInterval(() => {
                setServerTime((prevTime) => {
                    if (!prevTime) return null;
                    const newTime = new Date(prevTime.getTime() + 1000);
                    const hours = newTime.getHours().toString().padStart(2, '0');                    const minutes = newTime.getMinutes().toString().padStart(2, '0');
                    const seconds = newTime.getSeconds().toString().padStart(2, '0');
                    setTime(`${hours}:${minutes}:${seconds}`);
                    
                    if (prevTime.getDate() !== newTime.getDate()) {
                        updateDate(newTime);
                    }
                    
                    return newTime;
                });
            }, 1000);
            setIntervalId(id);
        };

        const updateDate = (date: Date) => {
            const day = date.getDate();
            const suffix = ['th', 'st', 'nd', 'rd'][(day % 10 > 3 ? 0 : day % 10) * (day < 10 || day > 20 ? 1 : 0)] || 'th';
            const formattedDate = date.toLocaleDateString('en-US', {
                weekday: 'long',
                day: 'numeric',
                month: 'long',
                year: 'numeric'
            });
            const [weekday, month, dayNum, year] = formattedDate.replace(',', '').split(' ');
            setDate(`${weekday} ${day}${suffix} ${month} ${year}`);
        };

        fetchInitialTime();
        
        const initialDate = new Date();
        updateDate(initialDate);

        return () => {
            if (intervalId) {
                clearInterval(intervalId);
            }
        };
    }, []);

    return (
        <>
            <Description title={time} subtitle={date} />
        </>
    );
}
