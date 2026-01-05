import * as z from 'zod';
import { genericErrorRequired, type Config, type SchemaFromConfig } from '@/schema';

export function generateSchema<T extends Config>(config: T): z.ZodObject<SchemaFromConfig<T>> {
    const zodSchema = {} as SchemaFromConfig<T>;

    for (const [key, fieldConfig] of Object.entries(config) as [keyof T, T[keyof T]][]) {
        if ('validate' in fieldConfig && fieldConfig.validate === 'email') {
            zodSchema[key] = z.email({
                message: fieldConfig.error || genericErrorRequired,
            }) as unknown as SchemaFromConfig<T>[typeof key];
            continue;
        }

        if ('optional' in fieldConfig && fieldConfig.optional) {
            zodSchema[key] = z.string().optional() as unknown as SchemaFromConfig<T>[typeof key];
            continue;
        }

        if ('dependantKey' in fieldConfig) {
            zodSchema[key] = z.string().optional() as unknown as SchemaFromConfig<T>[typeof key];
            continue;
        }

        zodSchema[key] = z.string().min(1, {
            message: fieldConfig.error || genericErrorRequired,
        }) as unknown as SchemaFromConfig<T>[typeof key];
    }

    return z.object(zodSchema).superRefine((data, ctx) => {
        for (const [key, fieldConfig] of Object.entries(config)) {
            if ('dependantKey' in fieldConfig && 'dependantValue' in fieldConfig) {
                const dependantKey = fieldConfig.dependantKey as keyof typeof data;
                const dependantValue = fieldConfig.dependantValue;

                const dependantFieldValue = data[dependantKey!];
                const currentFieldValue = data[key as keyof typeof data];

                const shouldBeRequired = dependantFieldValue === dependantValue;

                if (shouldBeRequired && !currentFieldValue?.toString().trim()) {
                    ctx.addIssue({
                        path: [key],
                        code: 'custom',
                        message: fieldConfig.error || genericErrorRequired,
                    });
                }
            }
        }
    });
}
