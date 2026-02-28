import { fontFamily } from 'tailwindcss/defaultTheme';
import type { Config } from 'tailwindcss';
import tailwindcssAnimate from 'tailwindcss-animate';

const config: Config = {
	darkMode: ['class'],
	content: ['./src/**/*.{html,js,svelte,ts}'],
	safelist: ['dark'],
	theme: {
		container: {
			center: true,
			padding: '2rem',
			screens: {
				'2xl': '1400px'
			}
		},
		extend: {
			// Using Tailwind's default 4px grid for spacing
			// Custom font sizes with proper line heights (minor third scale ~1.2x)
			borderRadius: {
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)'
			},
			fontFamily: {
				sans: ['Inter', ...fontFamily.sans],
				serif: 'Silk Serif'
			},
			backgroundImage: {
				'primary-gradient': 'var(--primary-gradient)'
			},
			boxShadow: {
				'inset-top': 'inset 0 1px 1px 0 rgba(255,255,255,0.15)',
				'inset-top-sm': 'inset 0 1px 1px 0 rgba(255,255,255,0.08)',
				'inset-top-lg': 'inset 0 2px 2px 0 rgba(255,255,255,0.2)',
				'focus-primary': '0 0 0px 0px hsl(var(--background)), 0 0 4px 0px hsl(var(--primary) / 0.5)',
				'focus-primary-glow': '0 0 4px 0px hsl(var(--background)), 0 0 0 4px hsl(var(--primary) / 0.6), 0 0 16px hsl(var(--primary) / 0.3)'
			},
			keyframes: {
				'dialog-in': {
					from: { opacity: '0', transform: 'scale(0.95) translate(-50%, -50%)' },
					to: { opacity: '1', transform: 'scale(1) translate(-50%, -50%)' }
				},
				'dialog-out': {
					from: { opacity: '1', transform: 'scale(1) translate(-50%, -50%)' },
					to: { opacity: '0', transform: 'scale(0.95) translate(-50%, -50%)' }
				}
			},
			animation: {
				'dialog-in': 'dialog-in 0.2s ease-out',
				'dialog-out': 'dialog-out 0.15s ease-in'
			}
		}
	},
	plugins: [tailwindcssAnimate]
};

export default config;
