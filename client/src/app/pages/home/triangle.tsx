export function Triangle() {
    return (
        <div data-id-name='Triangle' className='relative'>
            <div className='absolute w-0 h-0 border-r-[120px] border-l-transparent border-l-[120px] border-r-transparent border-t-[203px] border-t-primary'></div>
            <div className='absolute left-[20px] top-[10px] w-0 h-0 border-r-[100px] border-l-transparent border-l-[100px] border-r-transparent border-t-[173.2px] border-base-300 items-center flex flex-col'>
                <div className='flex flex-col items-center absolute -top-44 text-primary'>
                    <span className='text-4xl'>MAGI</span>
                    <span className='text-3xl'>02</span>
                    <span className='tracking-tight'>ORIGINAL</span>
                </div>
            </div>
        </div>
    );
}
