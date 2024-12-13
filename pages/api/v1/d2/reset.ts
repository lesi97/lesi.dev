import { NextApiRequest, NextApiResponse } from 'next';
import { Http_Status_Code } from '@/lib/server/enums/enums';

export default async function GET(res: NextApiResponse, req: NextApiRequest) {
    function getNextTuesdayAtFivePM(): Date {
        const currentDateTime = new Date();
        const currentDayOfWeek = currentDateTime.getDay();
        let daysUntilNextTuesday = (9 - currentDayOfWeek) % 7;

        if (
            daysUntilNextTuesday === 0 &&
            (currentDateTime.getHours() > 17 || (currentDateTime.getHours() === 17 && currentDateTime.getMinutes() > 0))
        ) {
            daysUntilNextTuesday = 7;
        }

        const nextTuesdayDateTime = new Date(currentDateTime);
        nextTuesdayDateTime.setDate(currentDateTime.getDate() + daysUntilNextTuesday);
        nextTuesdayDateTime.setHours(17, 0, 0, 0);
        return nextTuesdayDateTime;
    }

    function calculateInterval(futureDateTime: Date): string {
        const currentDateTime = new Date();
        const diffInMilliseconds = futureDateTime.getTime() - currentDateTime.getTime();

        const sixHoursInMilliseconds = 6 * 60 * 60 * 1000;
        if (diffInMilliseconds <= 0 && diffInMilliseconds >= -sixHoursInMilliseconds) {
            return "Reset's here! Maybe gift a sub 👉👈";
        }

        const diffInSeconds = Math.floor(diffInMilliseconds / 1000);
        const days = Math.floor(diffInSeconds / (3600 * 24));
        const hours = Math.floor((diffInSeconds % (3600 * 24)) / 3600);
        const minutes = Math.floor((diffInSeconds % 3600) / 60);
        const seconds = diffInSeconds % 60;

        const countdownParts: string[] = [];
        if (days > 0) {
            countdownParts.push(`${days} days`);
        }
        countdownParts.push(`${hours} hours`);
        countdownParts.push(`${minutes} minutes`);
        countdownParts.push(`${seconds} seconds`);

        const countdownString = countdownParts.join(', ');
        const finalCountdown = `${countdownString} until reset!`;
        return finalCountdown;
    }

    const futureDateTime = getNextTuesdayAtFivePM();
    const time = calculateInterval(futureDateTime);

    // Always return 200, this is used as a Twitch Nightbot command and if the status code is an error
    // The status code is returned instead of a custom error message, e.g a 400 returns SERVER RESPONDED WITH 400 or something similar
    res.setHeader('Content-Type', 'text/plain; charset=utf-8');
    return res.status(Http_Status_Code.OK).send(time);
}