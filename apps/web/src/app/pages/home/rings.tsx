import { mergeClassNames } from '@/utils';
import { useId } from 'react';

export function Rings() {
    const gradientStops = 'from-error via-primary via-[60%] to-success';
    const gradient = 'bg-gradient-to-t from-error via-primary via-[60%] to-success';
    const colourRing = 'rounded-full absolute animate-spin-slow';
    const blackRing = 'rounded-full absolute bg-base-100 flex text-xl';
    const textAnim =
        'bg-gradient-to-b bg-clip-text text-transparent [-webkit-text-fill-color:transparent] inline-block bg-[length:100%_200%] bg-[position:50%_0%] animate-[text-gradient-y_3s_linear_infinite_alternate]';
    const enChars = 'absolute left-2 top-1/2 -translate-y-1/2 text-xl';
    const jpChars = 'absolute left-1/2 -bottom-[1px] -translate-x-1/2 text-lg font-bold';

    const rings = [
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[350px] h-[345px] z-[20]', blackRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[380px] h-[370px] z-[19]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '2',
            textJP: '二',
            classNames: mergeClassNames('w-[430px] h-[425px] z-[18]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[460px] h-[450px] z-[17]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '3',
            textJP: '三',
            classNames: mergeClassNames('w-[510px] h-[505px] z-[16]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[540px] h-[530px] z-[15]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '4',
            textJP: '四',
            classNames: mergeClassNames('w-[590px] h-[585px] z-[14]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[620px] h-[610px] z-[13]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '5',
            textJP: '五',
            classNames: mergeClassNames('w-[670px] h-[665px] z-[12]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[700px] h-[690px] z-[11]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '6',
            textJP: '六',
            classNames: mergeClassNames('w-[750px] h-[745px] z-[10]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[780px] h-[770px] z-[9]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '7',
            textJP: '七',
            classNames: mergeClassNames('w-[830px] h-[825px] z-[8]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[860px] h-[850px] z-[7]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '8',
            textJP: '八',
            classNames: mergeClassNames('w-[910px] h-[905px] z-[6]', blackRing),
            hasText: true,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[940px] h-[930px] z-[5]', gradient, colourRing),
            hasText: false,
        },
        {
            id: useId(),
            textEN: '',
            textJP: '',
            classNames: mergeClassNames('w-[990px] h-[985px] z-[4]', blackRing),
            hasText: false,
        },
    ];

    return (
        <div className='w-fit h-fit absolute z-10 flex items-center justify-center opacity-70 pointer-events-none'>
            {rings.map((ring) => {
                return (
                    <div className={ring.classNames} key={ring.id}>
                        {ring.hasText && (
                            <div className='absolute inset-0 pointer-events-none'>
                                <span className={mergeClassNames(enChars, gradientStops, textAnim)}>{ring.textEN}</span>
                                <span className={mergeClassNames(jpChars, gradientStops, textAnim)}>{ring.textJP}</span>
                            </div>
                        )}
                    </div>
                );
            })}
        </div>
    );
}
