import { useEffect, useMemo, useState } from 'react';
import { useTime } from '@/hooks';
import { TextType } from './textType';
import { cn } from '@/utils';

export function DateTime({ containerClassNames }: { containerClassNames?: string }) {
    const { time, date } = useTime();
    const [showDate, setShowDate] = useState<boolean>(false);
    const [showTime, setShowTime] = useState<boolean>(false);

    const title = 'MAGI SYSTEM UNIT 2.0';
    const dateLine = useMemo(() => `${date}`, [date]);
    const timeLine = useMemo(() => `${time}`, [time]);

    const showCursor1 = !showDate && !showTime;
    const showCursor2 = showDate && !showTime;
    const isOver = showDate && showTime;

    return (
        <div className={cn('flex text-left align- flex-col leading-none w-[608px]', containerClassNames)}>
            <div className='flex flex-row items-start text-left align-start'>
                <TextType
                    text={[title]}
                    typingSpeed={100}
                    pauseDuration={0}
                    loop={false}
                    showCursor={showCursor1}
                    cursorClassName='text-accent'
                    className='pr-6 text-5xl tracking-wider text-left'
                    onTypingComplete={() => {
                        setShowDate(true);
                    }}
                />
            </div>
            <div className={cn('flex h-12 flex-row justify-between')}>
                {showDate && (
                    <TextType
                        text={[dateLine]}
                        typingSpeed={100}
                        pauseDuration={0}
                        loop={false}
                        showCursor={showCursor2}
                        cursorClassName='text-accent'
                        className='text-5xl tracking-wider w-fit'
                        onTypingComplete={() => {
                            setShowTime(true);
                        }}
                    />
                )}
                {showTime && (
                    <div className='w-[233px]'>
                        <TextType
                            text={[timeLine]}
                            typingSpeed={100}
                            pauseDuration={0}
                            loop={false}
                            showCursor={true}
                            cursorClassName='text-accent'
                            className='text-5xl tracking-wider w-fit'
                        />
                    </div>
                )}
            </div>
        </div>
    );
}
