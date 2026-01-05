const Loaders = {
    Trio: () => (
        <div className='relative inline-block h-10 w-10 animate-spin'>
            <div className='absolute left-1/2 top-0 h-full w-1/4 -translate-x-1/2 rotate-[120deg]'>
                <div className='pb-full animate-wobble absolute left-0 top-0 h-1/4 w-full rounded-full bg-accent'></div>
            </div>
            <div className='absolute left-1/2 top-0 h-full w-1/4 -translate-x-1/2 -rotate-[120deg]'>
                <div className='pb-full animate-wobble absolute left-0 top-0 h-1/4 w-full rounded-full bg-accent'></div>
            </div>
            <div className='absolute left-1/2 top-0 h-full w-1/4 -translate-x-1/2'>
                <div className='pb-full animate-wobble absolute left-0 top-0 h-1/4 w-full rounded-full bg-accent'></div>
            </div>
        </div>
    ),
};

export default Loaders;
