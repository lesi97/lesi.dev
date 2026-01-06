import { Description } from '@/components/ui/description';
import { useTime } from '@/hooks/useTime';

export function Home() {
    const { time, date } = useTime();
    return (
        <>
            <Description title={time} subtitle={date} />
        </>
    );
}
