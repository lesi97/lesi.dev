import { Description } from '@/components/ui/description';
import { usePageMeta } from '@/hooks';
import { useTime } from '@/hooks/useTime';

export function Home() {
    usePageMeta({
        title: 'Home | Lesi',
        description: 'The homepage of lesi.dev',
    });
    const { time, date } = useTime();
    return (
        <>
            <Description title={time} subtitle={date} />
        </>
    );
}
