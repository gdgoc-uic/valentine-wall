import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.extend(relativeTime);

export function formatDateTime(dt: Date, format: string): string {
  return dayjs(dt).format(format);
}

export function toNow(dt: Date): string {
  return dayjs(dt).toNow(true);
}

export function fromNow(dt: Date): string {
  return dayjs(dt).fromNow();
}

export function prettifyDateTime(dt: Date | string): string {
  return formatDateTime(typeof dt === 'string' ? new Date(dt) : dt, 'MMMM D, YYYY h:mm A');
}