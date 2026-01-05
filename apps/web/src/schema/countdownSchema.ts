import z from 'zod';
import { genericErrorRequired, type Config } from '.';
import { generateSchema, generateForm } from '@/utils';

const countdownConfig = {
    target_date: {
        label: 'Target Date',
        type: 'date',
        error: genericErrorRequired,
    },
    message: { label: 'Message', type: 'text', error: genericErrorRequired },
    fallback_message: { label: 'Fallback Message', type: 'text', error: genericErrorRequired },
} as const satisfies Config;

type CountdownDataKeysType = keyof typeof countdownConfig;
const CountdownSchema = generateSchema(countdownConfig);
type CountdownSchemaType = z.infer<typeof countdownConfig>;
const { form: CountdownForm, initialState: CountdownInitialState } =
    generateForm<CountdownDataKeysType>(countdownConfig);

export { type CountdownDataKeysType, type CountdownSchemaType, CountdownSchema, CountdownForm, CountdownInitialState };
