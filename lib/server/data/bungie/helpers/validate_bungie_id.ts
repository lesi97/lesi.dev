import { Http_Status_Code } from '@/lib/server/enums/enums';
import { HTTP_Response } from '@/lib/server/types/types';

export function validateId(id: string): [string, number, null] | [null, null, HTTP_Response] {
    const isValid =
        /^[\w !@#$%^&*()_+={}\[\]:;"'<>,.?\/\\|-]+#[0-9]{4}$/.test(id) &&
        !/\b(drop|alter|delete|insert|update|create|select|truncate|exec|union)\b/i.test(id);
    if (!isValid) {
        const error = {
            message: 'Error: Invalid ID',
            status_code: Http_Status_Code.BadRequest,
        };
        return [null, null, error];
    }
    const [username, username_code] = id.split(/#(\d{4})$/).filter(Boolean);
    return [username, parseInt(username_code), null];
}
