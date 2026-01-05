import { ReactNode } from 'react';

export function Highlight({ children }: { children: ReactNode }): React.ReactNode {
    return <span className='text-accent'>{children}</span>;
}
