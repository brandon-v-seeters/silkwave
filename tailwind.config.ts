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
			fontSize: {
				'xs': ['0.75rem', { lineHeight: '1rem' }],
				'sm': ['0.875rem', { lineHeight: '1.25rem' }],
				'base': ['1rem', { lineHeight: '1.5rem' }],
				'lg': ['1.125rem', { lineHeight: '1.75rem' }],
				'xl': ['1.25rem', { lineHeight: '1.75rem' }],
				'2xl': ['1.5rem', { lineHeight: '2rem' }],
				'3xl': ['1.875rem', { lineHeight: '2.25rem' }],
				'4xl': ['2.25rem', { lineHeight: '2.5rem' }],
				'5xl': ['3rem', { lineHeight: '1.1' }],
				'6xl': ['3.75rem', { lineHeight: '1.1' }],
				'7xl': ['4.5rem', { lineHeight: '1.1' }],
			},
			colors: {
				border: 'hsl(var(--border) / <alpha-value>)',
				input: 'hsl(var(--input) / <alpha-value>)',
				ring: 'hsl(var(--ring) / <alpha-value>)',
				background: 'hsl(var(--background) / <alpha-value>)',
				foreground: {
					DEFAULT: 'hsl(var(--foreground) / <alpha-value>)',
					muted: 'hsl(var(--foreground-muted) / <alpha-value>)'
				},
				primary: {
					DEFAULT: 'hsl(var(--primary) / <alpha-value>)',
					foreground: 'hsl(var(--primary-foreground) / <alpha-value>)'
				},
				secondary: {
					DEFAULT: 'hsl(var(--secondary) / <alpha-value>)',
					foreground: 'hsl(var(--secondary-foreground) / <alpha-value>)'
				},
				destructive: {
					DEFAULT: 'hsl(var(--destructive) / <alpha-value>)',
					foreground: 'hsl(var(--destructive-foreground) / <alpha-value>)',
					background: 'var(--destructive-background)'
				},
				success: {
					DEFAULT: 'hsl(var(--success) / <alpha-value>)',
					background: 'hsl(var(--success-background) / <alpha-value>)'
				},
				info: {
					DEFAULT: 'hsl(var(--info) / <alpha-value>)',
					foreground: 'hsl(var(--info-foreground) / <alpha-value>)'
				},
				warning: {
					DEFAULT: 'hsl(var(--warning) / <alpha-value>)',
					foreground: 'hsl(var(--warning-foreground) / <alpha-value>)'
				},
				muted: {
					DEFAULT: 'hsl(var(--muted) / <alpha-value>)',
					foreground: 'hsl(var(--muted-foreground) / <alpha-value>)',
					background: 'hsl(var(--muted-background) / <alpha-value>)'
				},
				accent: {
					DEFAULT: 'hsl(var(--accent) / <alpha-value>)',
					foreground: 'hsl(var(--accent-foreground) / <alpha-value>)'
				},
				popover: {
					DEFAULT: 'hsl(var(--popover) / <alpha-value>)',
					foreground: 'hsl(var(--popover-foreground) / <alpha-value>)'
				},
				card: {
					DEFAULT: 'hsl(var(--card) / <alpha-value>)',
					foreground: 'hsl(var(--card-foreground) / <alpha-value>)'
				},
				sidebar: {
					DEFAULT: 'hsl(var(--sidebar-background) / <alpha-value>)',
					foreground: 'hsl(var(--sidebar-foreground) / <alpha-value>)',
					primary: 'hsl(var(--sidebar-primary) / <alpha-value>)',
					'primary-foreground': 'hsl(var(--sidebar-primary-foreground) / <alpha-value>)',
					accent: 'hsl(var(--sidebar-accent) / <alpha-value>)',
					'accent-foreground': 'hsl(var(--sidebar-accent-foreground) / <alpha-value>)',
					border: 'hsl(var(--sidebar-border) / <alpha-value>)',
					ring: 'hsl(var(--sidebar-ring) / <alpha-value>)'
				},
				link: {
					DEFAULT: 'hsl(var(--link) / <alpha-value>)',
					hover: 'hsl(var(--link-hover) / <alpha-value>)'
				}
			},
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
