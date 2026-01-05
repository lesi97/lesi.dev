import z from 'zod';

export const genericErrorRequired = 'This field is required';
export const genericErrorSelectAnOption = 'Please select an option';

export type StateErrorType = {
    type: 'unique' | 'general';
    value: string | undefined;
};

export const uniqueStateError: StateErrorType = {
    type: 'unique',
    value: undefined,
} as const;
export const generalStateError: StateErrorType = {
    type: 'general',
    value: undefined,
} as const;

export type FormQuestionsType<K extends string> = {
    label: string;
    id: string;
    className?: string;
    type: string;
    dataKey: K;
    options: string[];
    dependantKey: K | null;
    dependantValue: string | null;
    autoComplete: AutoCompleteType;
    requiredError?: string;
}[];

export type AutoCompleteType =
    | 'off'
    | 'street-address'
    | 'address-level1'
    | 'address-level2'
    | 'county'
    | 'postal-code'
    | 'country'
    | 'given-name'
    | 'family-name'
    | 'email'
    | 'tel'
    | 'organization'
    | 'honorific-prefix'
    | 'current-password'
    | 'rutjfkde';

type Inputs =
    | 'text'
    | 'email'
    | 'number'
    | 'password'
    | 'checkbox'
    | 'radio'
    | 'date'
    | 'file'
    | 'hidden'
    | 'search'
    | 'tel'
    | 'url'
    | 'time'
    | 'range'
    | 'color'
    | 'dropdown'
    | 'combobox'
    | 'multiselect'
    | 'textarea'
    | 'yesNo';

type ZodValidationKey = 'email' | 'url' | 'uuid' | 'min' | 'max' | 'regex' | 'date' | 'length';

type BaseConfig = {
    label: string;
    error: string;
    className?: string;
    validate?: ZodValidationKey;
    autoComplete?: AutoCompleteType;
    dependantKey?: string;
    dependantValue?: string | number | boolean;
    optional?: boolean;
};

interface OptionsConfig extends BaseConfig {
    type: 'dropdown' | 'combobox' | 'yesNo';
    options: string[];
}

interface TextConfig extends BaseConfig {
    type: 'text' | 'email' | 'password' | 'number' | 'money' | 'date';
    options?: never;
}

export type Config = Record<string, OptionsConfig | TextConfig>;

type ZodTypeFromField<C> = C extends { validate: 'email'; optional: true }
    ? z.ZodOptional<z.ZodString>
    : C extends { validate: 'email' }
      ? z.ZodString
      : C extends { optional: true }
        ? z.ZodOptional<z.ZodString>
        : z.ZodString;

export type SchemaFromConfig<T extends Record<string, any>> = {
    [K in keyof T]: ZodTypeFromField<T[K]>;
};

export type MimeType =
    // Images
    | 'image/*'
    | 'image/jpeg'
    | 'image/png'
    | 'image/gif'
    | 'image/webp'
    | 'image/bmp'
    | 'image/svg+xml'
    | 'image/tiff'
    | 'image/x-icon'

    // Videos
    | 'video/*'
    | 'video/mp4'
    | 'video/webm'
    | 'video/ogg'
    | 'video/avi'
    | 'video/mpeg'
    | 'video/quicktime'
    | 'video/x-ms-wmv'

    // Audio
    | 'audio/*'
    | 'audio/mpeg'
    | 'audio/wav'
    | 'audio/ogg'
    | 'audio/aac'
    | 'audio/flac'
    | 'audio/x-midi'

    // Documents
    | 'application/pdf'
    | 'application/msword'
    | 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    | 'application/vnd.ms-excel'
    | 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    | 'application/vnd.ms-powerpoint'
    | 'application/vnd.openxmlformats-officedocument.presentationml.presentation'
    | 'text/plain'
    | 'text/csv'
    | 'text/html'
    | 'application/json'
    | 'application/xml'

    // Archives
    | 'application/zip'
    | 'application/x-tar'
    | 'application/x-7z-compressed'
    | 'application/gzip'

    // Applications
    | 'application/javascript'
    | 'application/x-sh'
    | 'application/java-archive'
    | 'application/x-httpd-php'
    | 'application/octet-stream'
    | 'application/x-dosexec'

    // Fonts
    | 'font/otf'
    | 'font/ttf'
    | 'font/woff'
    | 'font/woff2';
