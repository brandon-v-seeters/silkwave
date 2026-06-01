<script lang="ts" module>
	export type BreadcrumbItem = {
		label: string;
		href?: string;
	};

	const routeTrails: Record<string, BreadcrumbItem[]> = {
		'/': [{ label: 'Artists' }, { label: 'Featured Releases' }],
		'/settings': [{ label: 'Account' }, { label: 'Settings' }],
		'/upload': [{ label: 'Artist Studio' }, { label: 'Upload Music' }],
		'/upload/drafts': [
			{ label: 'Artist Studio' },
			{ label: 'Upload Music', href: '/upload' },
			{ label: 'Drafts' }
		],
		'/upload/release': [
			{ label: 'Artist Studio' },
			{ label: 'Upload Music', href: '/upload' },
			{ label: 'New Release' }
		]
	};

	const discoverLabels = [
		{ key: 'shuffle', value: 'true', label: 'Shuffle Play' },
		{ key: 'type', value: 'tracks', label: 'Tracks' },
		{ key: 'view', value: 'genres', label: 'Genres' },
		{ key: 'view', value: 'charts', label: 'Charts' },
		{ key: 'view', value: 'radio', label: 'Radio' },
		{ key: 'sort', value: 'recent', label: 'Recently Played' },
		{ key: 'filter', value: 'favorites', label: 'Favorite Tracks' }
	];

	function labelForDiscover(url: URL) {
		if (url.searchParams.has('q')) return 'Search Results';

		return (
			discoverLabels.find(({ key, value }) => url.searchParams.get(key) === value)?.label ??
			'Discover'
		);
	}

	function titleCase(value: string) {
		return decodeURIComponent(value)
			.replace(/[-_]+/g, ' ')
			.replace(/\b\w/g, (letter) => letter.toUpperCase());
	}

	export function breadcrumbsForUrl(url: URL): BreadcrumbItem[] {
		if (url.pathname === '/discover') {
			return [{ label: 'Artists', href: '/' }, { label: labelForDiscover(url) }];
		}

		const configuredTrail = routeTrails[url.pathname];
		if (configuredTrail) return configuredTrail;

		const segments = url.pathname.split('/').filter(Boolean);
		if (segments.length === 0) return routeTrails['/'];

		return [
			{ label: 'Artists', href: '/' },
			...segments.map((segment, index) => {
				const href = `/${segments.slice(0, index + 1).join('/')}`;
				const isLast = index === segments.length - 1;

				return {
					label: titleCase(segment),
					href: isLast ? undefined : href
				};
			})
		];
	}
</script>

<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { cn } from '$lib/utils/utils.js';

	type BreadcrumbsProps = {
		items?: BreadcrumbItem[];
		class?: string;
	};

	let { items, class: className }: BreadcrumbsProps = $props();

	const breadcrumbs = $derived.by(() => items ?? breadcrumbsForUrl(page.url));
</script>

{#if breadcrumbs.length > 0}
	<nav aria-label="Breadcrumb" class={cn('min-w-0 text-muted-foreground flex', className)}>
		<ol class="flex min-w-0 items-center gap-3">
			{#each breadcrumbs as item, index (`${item.label}-${item.href ?? index}`)}
				{@const isCurrent = index === breadcrumbs.length - 1}

				{#if index > 0}
					<li aria-hidden="true" class="hidden text-border sm:block">›</li>
				{/if}

				<li class={isCurrent ? 'min-w-0' : 'hidden shrink-0 sm:block'}>
					{#if item.href && !isCurrent}
						<a
							href={resolve(item.href as '/')}
							class="transition hover:text-foreground focus-visible:text-foreground focus-visible:outline-none"
						>
							{item.label}
						</a>
					{:else}
						<span
							class={isCurrent
								? 'block truncate font-medium text-foreground'
								: undefined}
							aria-current={isCurrent ? 'page' : undefined}
						>
							{item.label}
						</span>
					{/if}
				</li>
			{/each}
		</ol>
	</nav>
{/if}
