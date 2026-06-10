<script lang="ts" module>
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from 'svelte/elements';
	import { type VariantProps, tv } from 'tailwind-variants';

	export const buttonVariants = tv({
		base: 'cursor-pointer ring-offset-background focus-visible:ring-ring inline-flex items-center justify-center gap-2 overflow-hidden whitespace-nowrap text-base font-medium transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0',
		variants: {
			variant: {
				default:
					'bg-foreground text-background hover:bg-foreground/85 shadow-inset-top rounded-full',
				primary:
					'bg-primary text-primary-foreground hover:bg-primary/85 shadow-inset-top border-b-2 border-background/40 border-t-1 border-t-foreground/20',
				destructive:
					'bg-destructive text-foreground hover:bg-destructive/90 shadow-inset-top',
				outline:
					'border-input bg-background hover:bg-accent hover:text-accent-foreground border',
				secondary: 'bg-foreground/10 text-foreground hover:opacity-80 shadow-inset-top-sm',
				ghost: 'hover:bg-accent hover:text-accent-foreground',
				link: 'text-primary underline-offset-4 hover:underline',
				gradient: 'bg-foreground text-background hover:bg-foreground/85 shadow-inset-top'
			},
			size: {
				default: 'h-12 px-4 py-2',
				sm: 'h-9 px-3',
				lg: 'h-14 px-8',
				icon: '!h-10 !w-10'
			}
		},
		defaultVariants: {
			variant: 'default',
			size: 'default'
		}
	});

	export type ButtonVariant = VariantProps<typeof buttonVariants>['variant'];
	export type ButtonSize = VariantProps<typeof buttonVariants>['size'];

	export type ButtonProps = WithElementRef<HTMLButtonAttributes> &
		WithElementRef<HTMLAnchorAttributes> & {
			variant?: ButtonVariant;
			size?: ButtonSize;
		};
</script>

<script lang="ts">
	import { cn } from '$lib/utils/utils.js';

	let {
		class: className,
		variant = 'default',
		size = 'default',
		ref = $bindable(null),
		href = undefined,
		type = 'button',
		children,
		...restProps
	}: ButtonProps = $props();
</script>

{#if href}
	<a
		bind:this={ref}
		class={cn('sw-button', buttonVariants({ variant, size }), className)}
		data-variant={variant}
		{href}
		{...restProps}
	>
		<span class="sw-button__fill" aria-hidden="true"></span>
		<span class="sw-button__content">
			{@render children?.()}
		</span>
	</a>
{:else}
	<button
		bind:this={ref}
		class={cn('sw-button', buttonVariants({ variant, size }), className)}
		data-variant={variant}
		{type}
		{...restProps}
	>
		<span class="sw-button__fill" aria-hidden="true"></span>
		<span class="sw-button__content">
			{@render children?.()}
		</span>
	</button>
{/if}

<style>
	:global(.sw-button) {
		--sw-button-fill: color-mix(in oklch, var(--foreground) 12%, transparent);
		--sw-button-fill-opacity: 1;
		--sw-button-ease: cubic-bezier(0.16, 1, 0.3, 1);
		isolation: isolate;
		position: relative;
		transform: translateZ(0);
		transition:
			color 220ms ease,
			background-color 220ms ease,
			border-color 220ms ease,
			box-shadow 220ms ease,
			transform 180ms var(--sw-button-ease);
	}

	:global(.sw-button[data-variant='default']),
	:global(.sw-button[data-variant='gradient']) {
		--sw-button-fill: color-mix(in oklch, var(--foreground) 86%, var(--primary));
	}

	:global(.sw-button[data-variant='primary']) {
		--sw-button-fill: color-mix(in oklch, var(--primary) 82%, var(--foreground));
	}

	:global(.sw-button[data-variant='destructive']) {
		--sw-button-fill: color-mix(in oklch, var(--destructive) 82%, var(--foreground));
	}

	:global(.sw-button[data-variant='outline']),
	:global(.sw-button[data-variant='secondary']),
	:global(.sw-button[data-variant='ghost']) {
		--sw-button-fill: color-mix(in oklch, var(--foreground) 10%, transparent);
	}

	:global(.sw-button[data-variant='link']) {
		--sw-button-fill-opacity: 0;
	}

	:global(.sw-button:not(:disabled):not([aria-disabled='true']):active) {
		transform: translateY(1px);
	}

	:global(.sw-button__fill) {
		background: var(--sw-button-fill);
		border-radius: inherit;
		inset: -1px;
		opacity: var(--sw-button-fill-opacity);
		pointer-events: none;
		position: absolute;
		transform: translate3d(0, 108%, 0) scaleY(0.42);
		transform-origin: 50% 100%;
		transition:
			transform 520ms var(--sw-button-ease),
			opacity 260ms ease;
		z-index: 0;
	}

	:global(.sw-button__content) {
		align-items: center;
		display: inline-flex;
		gap: inherit;
		justify-content: center;
		position: relative;
		transform: translateY(0);
		transition: transform 420ms var(--sw-button-ease);
		z-index: 1;
	}

	:global(.sw-button:not(:disabled):not([aria-disabled='true']):hover .sw-button__fill),
	:global(.sw-button:not(:disabled):not([aria-disabled='true']):focus-visible .sw-button__fill) {
		transform: translate3d(0, 0, 0) scaleY(1);
	}

	:global(.sw-button:not(:disabled):not([aria-disabled='true']):hover .sw-button__content),
	:global(.sw-button:not(:disabled):not([aria-disabled='true']):focus-visible .sw-button__content) {
		transform: translateY(-1px);
	}

	:global(.sw-button:not(:disabled):not([aria-disabled='true']):active .sw-button__content) {
		transform: translateY(0);
	}

	@media (prefers-reduced-motion: reduce) {
		:global(.sw-button),
		:global(.sw-button__fill),
		:global(.sw-button__content) {
			transition-duration: 0.01ms;
		}

		:global(.sw-button__fill) {
			opacity: 0;
			transform: none;
		}

		:global(.sw-button:not(:disabled):not([aria-disabled='true']):hover .sw-button__fill),
		:global(.sw-button:not(:disabled):not([aria-disabled='true']):focus-visible .sw-button__fill) {
			opacity: var(--sw-button-fill-opacity);
		}

		:global(.sw-button:not(:disabled):not([aria-disabled='true']):hover .sw-button__content),
		:global(.sw-button:not(:disabled):not([aria-disabled='true']):focus-visible .sw-button__content),
		:global(.sw-button:not(:disabled):not([aria-disabled='true']):active .sw-button__content) {
			transform: none;
		}
	}
</style>
