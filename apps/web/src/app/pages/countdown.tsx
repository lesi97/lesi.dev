// import { useState, useEffect, useRef } from 'react';
// import { Description, TwitchMessageContainer } from '@/components/layout';
// import { Link } from 'react-router-dom';
// import { Button, Radio, Input } from '@/components/ui';
// import { ZodError } from 'zod';
// import { parseError } from '@/utils';
// import { useCountdown } from '@/hooks';

// type ErrorsType = {
//     target_date: null | string;
//     message: null | string;
//     fallback_message: null | string;
// };

// export default function Countdown() {
//     const {
//         transformCountdown,
//         countdownPlaceholders,
//         CountdownSchema,
//         commandName,
//         form,
//         placeholders,
//         handleSubmit,
//         setEditCommand,
//         editCommand,
//         setCommandName,
//     } = useCountdown();

//     return (
//         <>
//             <Description
//                 title='Countdown Generator'
//                 subtitle='Generate a countdown to use as a command'
//                 className='mb-0'
//             />
//             <TwitchMessageContainer user='C_Lesi'>!{commandName}</TwitchMessageContainer>
//             <TwitchMessageContainer user='Nightbot' className='mb-4'>
//                 {transformCountdown(form.data.target_date) !== 'Passed' ? (
//                     <span className='inline-block w-fit'>
//                         {transformCountdown(form.data.target_date)}{' '}
//                         {form.data.message !== '' ? form.data.message : placeholders.message}{' '}
//                         <Link to='https://twitch.tv/c_lesi' target='_blank' className='inline-block'>
//                             <span className='font-bold text-accent'>@C_Lesi</span>
//                         </Link>
//                     </span>
//                 ) : (
//                     <>
//                         {form.data.fallback_message !== '' ? form.data.fallback_message : placeholders.fallback_message}
//                     </>
//                 )}
//             </TwitchMessageContainer>

//             <form className='flex w-full flex-col gap-8' onSubmit={handleSubmit}>
//                 <div className='grid w-full grid-cols-2 items-end gap-6'>
//                     <Radio
//                         id='addCommand'
//                         name='commandAction'
//                         onChange={() => setEditCommand(false)}
//                         label='Add Command'
//                         className='w-full py-6'
//                         variant='default'
//                         checked={editCommand === false}
//                     />
//                     <Radio
//                         id='editCommand'
//                         name='commandAction'
//                         onChange={() => setEditCommand(true)}
//                         variant='default'
//                         label='Edit Command'
//                         className='w-full py-6'
//                         checked={editCommand === true}
//                     />
//                     <label htmlFor='commandName' className='flex flex-col justify-center gap-1'>
//                         <span className='w-full text-sm'>Command Name:</span>
//                         <Input
//                             id='commandName'
//                             className='w-full text-center'
//                             value={commandName}
//                             onChange={(e) => {
//                                 setCommandName(e.target.value);
//                             }}
//                             variant='outline'
//                         />
//                     </label>
//                     <label htmlFor='date' className='flex flex-col justify-center gap-1'>
//                         <span className='w-full text-sm'>Choose your target date:</span>
//                         <Input
//                             id='date'
//                             type='datetime-local'
//                             className='w-full text-center'
//                             value={data.target_date}
//                             onChange={(e) => {
//                                 setData({ ...data, target_date: e.target.value });
//                             }}
//                             readOnly={!!command}
//                             variant='outline'
//                             error={errors.target_date}
//                         />
//                     </label>

//                     <label htmlFor='event' className='flex flex-col justify-center gap-1'>
//                         <span className='w-full text-sm'>Type the event that the countdown is for:</span>
//                         <Input
//                             className='w-full text-center'
//                             id='event'
//                             value={data.message}
//                             onChange={(e) => {
//                                 setData({ ...data, message: e.target.value });
//                             }}
//                             readOnly={!!command}
//                             variant='outline'
//                             error={errors.message}
//                             placeholder={placeholders.message || ''}
//                         />
//                     </label>
//                     <label htmlFor='fallback' className='flex flex-col justify-center gap-1'>
//                         <span className='w-full text-sm'>Type a message to appear if the date has passed:</span>
//                         <Input
//                             className='w-full text-center'
//                             id='fallback'
//                             value={data.fallback_message}
//                             onChange={(e) => {
//                                 setData({ ...data, fallback_message: e.target.value });
//                             }}
//                             readOnly={!!command}
//                             variant='outline'
//                             error={errors.fallback_message}
//                             placeholder={placeholders.fallback_message || ''}
//                         />
//                     </label>
//                 </div>
//                 <Button type='submit' variant='gradient' size='xl' disabled={!!command}>
//                     Generate & Copy Nightbot Command
//                 </Button>
//                 {command && (
//                     <div className='flex flex-col'>
//                         <textarea
//                             readOnly
//                             id='generatedCommand'
//                             className='h-fit min-h-[100px] resize-none text-pretty rounded bg-primary/80 p-2 text-primary-content focus:outline-none'
//                             ref={commandRef}
//                             value={`${editCommand ? '!editcom' : '!addcom'} !${commandName} $(urlfetch ${command}) @$(touser)`}
//                         />
//                     </div>
//                 )}
//             </form>
//         </>
//     );
// }
