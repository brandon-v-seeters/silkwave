<script lang="ts">
	import duocolor from '../../../assets/duocolor.json';
	import filled from '../../../assets/filled.json';
	import line from '../../../assets/line.json';
	import duotone from '../../../assets/duotone.json';

	import type { IconKey } from '../../types/Icon';

	const svgs = {
		duocolor,
		filled,
		line,
		duotone
	};

	export let icon: IconKey;
	export let variant: 'duocolor' | 'filled' | 'line' | 'duotone' = 'duocolor';

	let icons: Record<string, string> = svgs[variant];
</script>

<svg
	class="icon {$$props.class}"
	viewBox={icon === 'framer' ? '0 0 33 50' : '0 0 24 24'}
	xmlns="http://www.w3.org/2000/svg"
>
	{@html icons[icon]}
</svg>
