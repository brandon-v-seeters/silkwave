// Vendor
import { getContext, setContext } from 'svelte';
import { toast } from 'svelte-sonner';

// Types
import type { WizardTrack } from '$lib/types/WizardTrack';

export type ReleasePricing = {
	basePrice: number;      // In cents (999 = $9.99)
	minimumPrice: number | null;  // For "pay what you want"
	trackPrice: number | null;    // Individual track price
};

// Services
import { releaseUploadService } from '$lib/services/release-upload.svelte';
import type { ReleaseWithTracks } from '$lib/types/generated';

// Helpers
function generateId(): string {
	return Math.random().toString(36).substring(2, 9);
}

function createTrackFromFile(file: File): WizardTrack {
	return {
		id: generateId(),
		title: file.name.replace(/\.[^/.]+$/, '').replace(/_/g, ' '),
		file,
		duration: '0:00',
		order: 0
	};
}

function formatTrackDuration(seconds: number): string {
	const minutes = Math.floor(seconds / 60);
	const remainingSeconds = seconds % 60;
	return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}

const WIZARD_KEY = Symbol('release-wizard');

export type Step = {
	id: string;
	title: string;
	description: string;
};

export function createWizardContext() {
	// Step configuration - add/remove steps here
	const steps: Step[] = [
		{ id: 'details', title: 'Release details', description: 'Name your release and pick genres' },
		{ id: 'artwork', title: 'Cover artwork', description: 'Upload your release cover' },
		{ id: 'tracks', title: 'Add tracks', description: 'Upload and organize your tracks' },
		{ id: 'pricing', title: 'Pricing', description: 'Set your release price' }
	];

	// Dynamic total steps
	const totalSteps = steps.length;

	// Step management
	let currentStep = $state(0);

	// Release state
	let releaseTitle = $state('');
	let description = $state('');
	let genres = $state<string[]>([]);
	let tracks = $state<WizardTrack[]>([]);
	let coverArt = $state<File | null>(null);
	let coverArtPreview = $state<string | null>(null);
	let releaseKey = $state<string | null>(null);

	// Pricing state (all values in USD cents)
	let pricing = $state<ReleasePricing>({
		basePrice: 0,
		minimumPrice: null,
		trackPrice: null
	});

	const userCtx = getContext<{ current: any }>('user');
	let artistKey = $state<string>(userCtx.current?.artist?._key);

	// Available genres
	const suggestedGenres = [
		'House', 'Techno', 'Trance', 'Drum & Bass', 'Dubstep', 'Deep House',
		'Progressive House', 'Tech House', 'Minimal', 'Electro', 'Ambient',
		'Downtempo', 'Breakbeat', 'Garage', 'Hardcore', 'Hardstyle', 'Hip Hop',
		'R&B', 'Pop', 'Rock', 'Indie', 'Jazz', 'Classical', 'Electronic', 'Experimental'
	];

	// Step validations - index matches step index
	const stepValidations = $derived([
		releaseTitle.trim() !== '' && genres.length > 0,  // Step 0: details
		coverArt !== null,                                  // Step 1: artwork
		tracks.length > 0 && tracks.every((t) => t.title && t.file), // Step 2: tracks
		true                                                // Step 3: pricing (always valid, price can be 0 for free)
	]);

	// Dynamic step validation getter
	const isStepValid = (index: number) => stepValidations[index] ?? false;

	// Check if all steps are valid
	const isFormValid = $derived(stepValidations.every(Boolean));

	// Check if current step allows proceeding
	const canProceed = $derived(isStepValid(currentStep));

	// Legacy validators for backwards compatibility
	const isStep1Valid = $derived(stepValidations[0]);
	const isStep2Valid = $derived(stepValidations[1]);
	const isStep3Valid = $derived(stepValidations[2]);
	const isStep4Valid = $derived(stepValidations[3]);

	function initWizard(draft: ReleaseWithTracks) {
		releaseTitle = draft.title;
		genres = draft.genres || [];
		tracks = draft.tracks.map((track, index) => ({
			id: track.id || generateId(),
			title: track.title,
			file: new File([], track.files.original.path.split('/').pop() || track.title),
			duration: track.durationDisplay || formatTrackDuration(track.duration),
			order: track.order || index + 1
		}));
	}
	// Navigation
	function nextStep() {
		if (currentStep < totalSteps - 1 && canProceed) {
			currentStep++;
		}
	}

	function prevStep() {
		if (currentStep > 0) currentStep--;
	}

	function goToStep(step: number) {
		if (step < currentStep) currentStep = step;
	}

	function getStepState(index: number): 'completed' | 'current' | 'upcoming' {
		if (index < currentStep) return 'completed';
		if (index === currentStep) return 'current';
		return 'upcoming';
	}

	// Cover art handler
	function setCoverArt(file: File) {
		coverArt = file;
	}

	// Track handlers
	function updateTrackTitle(id: string, title: string) {
		tracks = tracks.map((t) => (t.id === id ? { ...t, title } : t));
	}

	function updateTrackFile(id: string, file: File) {
		tracks = tracks.map((t) => {
			if (t.id === id) {
				const autoTitle = t.title || file.name.replace(/\.[^/.]+$/, '').replace(/_/g, ' ');
				return { ...t, file, title: autoTitle };
			}
			return t;
		});

		// Get audio duration
		const audio = new Audio();
		audio.src = URL.createObjectURL(file);
		audio.onloadedmetadata = () => {
			const minutes = Math.floor(audio.duration / 60);
			const seconds = Math.floor(audio.duration % 60);
			tracks = tracks.map((t) => {
				if (t.id === id) {
					return { ...t, duration: `${minutes}:${seconds.toString().padStart(2, '0')}` };
				}
				return t;
			});
			URL.revokeObjectURL(audio.src);
		};
	}

	function moveTrack(index: number, direction: 'up' | 'down') {
		const newIndex = direction === 'up' ? index - 1 : index + 1;
		if (newIndex < 0 || newIndex >= tracks.length) return;
		const newTracks = [...tracks];
		[newTracks[index], newTracks[newIndex]] = [newTracks[newIndex], newTracks[index]];
		tracks = newTracks;
	}

	function removeTrack(id: string) {
		tracks = tracks.filter((t) => t.id !== id);
	}

	function addBulkTracks(files: File[]) {
		files.forEach((file) => {
			const newTrack = createTrackFromFile(file);
			tracks = [...tracks, newTrack];

			const audio = new Audio();
			audio.src = URL.createObjectURL(file);
			audio.onloadedmetadata = () => {
				const minutes = Math.floor(audio.duration / 60);
				const seconds = Math.floor(audio.duration % 60);
				tracks = tracks.map((t) => {
					if (t.id === newTrack.id) {
						return { ...t, duration: `${minutes}:${seconds.toString().padStart(2, '0')}` };
					}
					return t;
				});
				URL.revokeObjectURL(audio.src);
			};
		});
	}

	// Save & Publish
	async function saveAsDraft() {
		if (releaseUploadService.isUploading) return;

		try {
			const result = await releaseUploadService.saveAsDraft(releaseTitle, genres, tracks, coverArt, artistKey);
			console.log('releaseKey', releaseKey);
			releaseKey = result.releaseKey;
		} catch (error) {
			console.error('Failed to save as draft:', error);
			toast.error('Failed to save as draft');
		}
	}

	async function publishRelease() {
		if (releaseUploadService.isUploading) return;
		setTimeout(() => {
			toast.success('Release published successfully');
		}, 2000);
	}

	const context = {
		// State (getters for reactivity)
		get currentStep() { return currentStep; },
		get totalSteps() { return totalSteps; },
		get steps() { return steps; },
		get releaseTitle() { return releaseTitle; },
		set releaseTitle(value: string) { releaseTitle = value; },
		get genres() { return genres; },
		set genres(value: string[]) { genres = value; },
		get description() { return description; },
		set description(value: string) { description = value; },
		get tracks() { return tracks; },
		get coverArt() { return coverArt; },
		get coverArtPreview() { return coverArtPreview; },
		set coverArtPreview(value: string | null) { coverArtPreview = value; },
		get pricing() { return pricing; },
		set pricing(value: ReleasePricing) { pricing = value; },
		get suggestedGenres() { return suggestedGenres; },
		get releaseKey() { return releaseKey; },

		// Dynamic validation
		isStepValid,
		get isFormValid() { return isFormValid; },
		get canProceed() { return canProceed; },

		// Legacy validators (for backwards compatibility)
		get isStep1Valid() { return isStep1Valid; },
		get isStep2Valid() { return isStep2Valid; },
		get isStep3Valid() { return isStep3Valid; },
		get isStep4Valid() { return isStep4Valid; },

		// Upload store
		get isUploading() { return releaseUploadService.isUploading; },
		get uploadStatus() { return releaseUploadService.status; },

		// Navigation actions
		nextStep,
		prevStep,
		goToStep,
		getStepState,

		// Release actions
		setCoverArt,
		initWizard,

		// Track actions
		updateTrackTitle,
		updateTrackFile,
		moveTrack,
		removeTrack,
		addBulkTracks,

		// Submit actions
		saveAsDraft,
		publishRelease
	};

	setContext(WIZARD_KEY, context);
	return context;
}

export type WizardContext = ReturnType<typeof createWizardContext>;

export function getWizardContext(): WizardContext {
	return getContext<WizardContext>(WIZARD_KEY);
}
