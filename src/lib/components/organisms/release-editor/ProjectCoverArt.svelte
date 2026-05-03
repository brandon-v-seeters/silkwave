<script lang="ts">
	import Icon from '$lib/components/atoms/Icon.svelte';
	import { generateId } from '$lib/utils/utils';

	interface Props {
		preview?: string | null;
		onchange?: (file: File) => void;
	}

	let { preview = $bindable(null), onchange }: Props = $props();

	const inputId = `cover-art-${generateId()}`;

	const processFile = (file: File) => {
		const reader = new FileReader();
		reader.onload = (e) => {
			preview = e.target?.result as string;
		};
		reader.readAsDataURL(file);
		onchange?.(file);
	};

	const handleDrop = (e: DragEvent) => {
		e.preventDefault();
		const files = e.dataTransfer?.files;
		if (files && files.length > 0) {
			const file = files[0];
			if (file.type.startsWith('image/')) {
				processFile(file);
			}
		}
	};

	const handleSelect = (e: Event) => {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			const file = input.files[0];
			if (file.type.startsWith('image/')) {
				processFile(file);
			}
		}
	};

	const openFilePicker = () => {
		document.getElementById(inputId)?.click();
	};
</script>

<div class="relative w-full">
	<!-- Rotating disc accent, peeks from behind the cover -->
	<div
		class="pointer-events-none absolute top-[2%] aspect-square w-[96%] opacity-0 transition-all duration-700 ease-out {preview
			? '-right-[10%] opacity-100'
			: 'right-[2%]'}"
		aria-hidden="true"
	>
		<div class="disc-spin relative h-full w-full rounded-full">
			<div
				class="absolute inset-0 rounded-full bg-foreground/10"
				style="background: radial-gradient(circle at 30% 30%, oklch(from var(--color-foreground) l c h / 0.12), oklch(from var(--color-foreground) l c h / 0.06) 60%, oklch(from var(--color-foreground) l c h / 0.1));"
			></div>
			<div
				class="absolute inset-[6%] rounded-full"
				style="background: repeating-radial-gradient(circle at center, transparent 0, transparent 3%, oklch(from var(--color-background) l c h / 0.06) 3%, oklch(from var(--color-background) l c h / 0.06) 3.4%);"
			></div>
			<div class="absolute inset-[42%] rounded-full bg-primary/40"></div>
			<div class="absolute inset-[48%] rounded-full bg-background"></div>
		</div>
	</div>

	<!-- Cover art -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="group relative z-10 aspect-square w-full cursor-pointer overflow-hidden rounded-2xl shadow-lg transition-transform duration-300"
		style="will-change: transform;"
		ondragover={(e) => e.preventDefault()}
		ondrop={handleDrop}
		onclick={openFilePicker}
		onkeydown={(e) => e.key === 'Enter' && openFilePicker()}
		role="button"
		tabindex="0"
	>
		{#if preview}
			<img
				src={preview}
				alt="Cover art"
				class="h-full w-full select-none object-cover transition-all duration-300"
				draggable="false"
			/>
		{:else}
			<div class="flex h-full w-full items-center justify-center bg-secondary">
				<Icon icon="image" variant="filled" class="h-16 w-16 fill-foreground-muted/30" />
			</div>
		{/if}

		<!-- Hover overlay -->
		<div
			class="absolute inset-0 z-10 flex cursor-pointer items-end justify-center rounded-2xl bg-transparent transition-all duration-200 ease-out sm:bg-black/0 sm:group-hover:bg-black/50"
		>
			<span
				class="relative mb-10 hidden translate-y-2 text-sm font-medium text-white opacity-0 transition-all duration-200 sm:block sm:group-hover:translate-y-0 sm:group-hover:opacity-100"
			>
				Change Cover Art
			</span>
		</div>
	</div>
</div>
<input id={inputId} type="file" accept="image/*" class="hidden" onchange={handleSelect} />

<style>
	.disc-spin {
		animation: disc-rotate 24s linear infinite;
	}

	@keyframes disc-rotate {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.disc-spin {
			animation: none;
		}
	}
</style>
