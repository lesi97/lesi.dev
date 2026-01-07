import { useState, useRef } from 'react';
import { Description, Button} from '@/components/ui';
import { useMinifier } from '@/hooks/useMinifier';

export function Minifier() {
    const {minify, unminify} = useMinifier()
    const [code, setCode] = useState('');
    const textAreaRef = useRef<HTMLTextAreaElement | null>(null);

    function handleMinify() {
        const newCode = minify(code);
        if (typeof newCode !== 'string') {
            return;
        }
        if (newCode) {
            setCode(newCode);
        }
    }

    function handleUnminify() {
        const newCode = unminify(code);
        if (typeof newCode !== 'string') {
            return;
        }
        if (newCode) {
            setCode(newCode);
        }
    }

    function copyCode() {
        if (!textAreaRef.current) {
            return;
        }
        textAreaRef.current.select();
        document.execCommand('copy');
    }

    return (
        <>
            <Description title='Minifier' subtitle='Minify or unminfiy CSS, JS, XML & JSON' />
            <textarea
                ref={textAreaRef}
                className='w-full rounded border border-accent bg-base-200 shadow-sm'
                rows={10}
                spellCheck='false'
                value={code}
                onChange={(e) => setCode(e.target.value)}></textarea>

            <div className='flex flex-col justify-between gap-4 pt-4 sm:flex-row lg:px-4'>
                <Button variant='secondary' onClick={handleMinify}>
                    MINIFY
                </Button>
                <Button variant='secondary' onClick={handleUnminify}>
                    UNMINIFY
                </Button>
                <Button variant='secondary' onClick={copyCode}>
                    COPY
                </Button>
            </div>
        </>
    );
}
