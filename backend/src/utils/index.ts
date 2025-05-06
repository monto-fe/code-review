import dayjs from 'dayjs';
export const getUnixTimestamp = (date?: string | number | Date): number => dayjs(date).unix();

export const DefaultPassword = '12345678';

export const Administrator = 'admin';