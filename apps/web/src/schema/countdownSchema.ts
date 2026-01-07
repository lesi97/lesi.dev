import z from 'zod';

export const countdownSchema = z.object({
    target_date: z.string(),
    message: z.string(),
    fallback_message: z.string(),
});
