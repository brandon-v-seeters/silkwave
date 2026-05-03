// Client-only type - icons don't exist in Go backend
import duocolor from '../../assets/duocolor.json';

export type IconKey = keyof typeof duocolor;
export type IconVariant = 'duocolor' | 'filled' | 'line' | 'duotone';
