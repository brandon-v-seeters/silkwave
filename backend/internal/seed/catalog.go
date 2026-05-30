package seed

import (
	"fmt"
	"strings"
	"time"
)

const (
	DefaultNamespace = "silkwave-dev"
	DefaultVersion   = "2026-05-dev-universe"
	DefaultPassword  = "Silkwave123!"
)

type seedMeta struct {
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
}

type catalog struct {
	Users         []document
	Artists       []document
	UserArtists   []document
	Releases      []document
	Tracks        []document
	Subscriptions []document
	Subscribers   []document
}

type document map[string]interface{}

type releaseFixture struct {
	Key         string
	ID          string
	ArtistKey   string
	Slug        string
	Title       string
	Description string
	Type        string
	Status      string
	PriceCents  int
	PublishAt   time.Time
	Genres      []string
	Tags        []string
	Moods       []string
	Tracks      []trackFixture
}

type trackFixture struct {
	Title    string
	Duration int
}

func buildCatalog(namespace, passwordHash string, now time.Time) catalog {
	meta := seedMeta{Namespace: namespace, Version: DefaultVersion}
	users := seedUsers(meta, passwordHash, now)
	artists := seedArtists(meta, now)
	releases := seedReleases(now)

	var docs catalog
	docs.Users = users
	docs.Artists = artists
	docs.UserArtists = []document{
		edge("user_artist_framer", "Users/user_artist", "Artists/artist_framer", meta, now),
		edge("user_artist_mira_vale", "Users/user_artist", "Artists/artist_mira_vale", meta, now),
		edge("user_artist_north_index", "Users/user_artist", "Artists/artist_north_index", meta, now),
	}

	for _, release := range releases {
		docs.Releases = append(docs.Releases, releaseDoc(release, meta, now))
		for index, track := range release.Tracks {
			docs.Tracks = append(docs.Tracks, trackDoc(release, track, index+1, meta, now))
		}
	}

	docs.Subscriptions = []document{
		subscriptionDoc("subscription_framer_supporter", "artist_framer", "Supporter", "Back Framer and stream the current catalog.", 5, "EUR", 10, meta, now),
		subscriptionDoc("subscription_framer_patron", "artist_framer", "Patron", "Higher monthly support for fans who want to keep the catalog moving.", 12, "EUR", 20, meta, now),
	}
	docs.Subscribers = []document{
		subscriberDoc("subscriber_framer_supporter_user_subscriber", "artist_framer", "user_subscriber", "subscription_framer_supporter", meta, now),
	}

	return docs
}

func seedUsers(meta seedMeta, passwordHash string, now time.Time) []document {
	return []document{
		userDoc("user_artist", "artist@silkwave.test", "artist", passwordHash, meta, now),
		userDoc("user_fan", "fan@silkwave.test", "fan", passwordHash, meta, now),
		userDoc("user_subscriber", "subscriber@silkwave.test", "subscriber", passwordHash, meta, now),
	}
}

func seedArtists(meta seedMeta, now time.Time) []document {
	return []document{
		artistDoc("artist_framer", "11111111-1111-4111-8111-111111111111", "Framer", "framer", "Framer builds precise electronic Releases with warm edges, night-drive bass, and hooks that arrive like signal through fog.", meta, now.AddDate(-3, 0, 0), now),
		artistDoc("artist_mira_vale", "11111111-1111-4111-8111-111111111112", "Mira Vale", "mira-vale", "Mira Vale makes ambient pop and glassy club sketches for quiet rooms with large windows.", meta, now.AddDate(-2, -4, 0), now),
		artistDoc("artist_north_index", "11111111-1111-4111-8111-111111111113", "North Index", "north-index", "North Index moves through heavier instrumental work, pressure-system drums, and low-lit melodic debris.", meta, now.AddDate(-2, -1, 0), now),
	}
}

func seedReleases(now time.Time) []releaseFixture {
	return []releaseFixture{
		{
			Key: "release_framer_moon_in_my_sky", ID: "22222222-2222-4222-8222-222222222201", ArtistKey: "artist_framer",
			Slug: "moon-in-my-sky", Title: "Moon in my Sky", Type: "album", Status: "published", PriceCents: 890,
			PublishAt: now.AddDate(-1, -2, 0), Genres: []string{"Electronic", "Downtempo"}, Tags: []string{"night drive", "melodic", "bass"}, Moods: []string{"luminous", "focused"},
			Description: "A patient electronic Release built around soft pressure, long shadows, and the small bright thing that refuses to leave the sky.",
			Tracks:      []trackFixture{{"Low Orbit Letters", 178}, {"Moon in my Sky", 214}, {"Afterimage Traffic", 196}, {"Small Hours Engine", 232}, {"Halogen Rain", 205}, {"Where the Signal Sleeps", 248}},
		},
		{
			Key: "release_framer_soft_static", ID: "22222222-2222-4222-8222-222222222202", ArtistKey: "artist_framer",
			Slug: "soft-static", Title: "Soft Static", Type: "ep", Status: "published", PriceCents: 490,
			PublishAt: now.AddDate(-1, -7, 0), Genres: []string{"Electronic", "Ambient"}, Tags: []string{"tape", "haze", "warm"}, Moods: []string{"soft", "restless"},
			Description: "Four pieces of warm interference and half-remembered melody.",
			Tracks:      []trackFixture{{"Soft Static", 184}, {"Bright Dust", 201}, {"Carrier Wave", 173}, {"Noon Ghost", 219}},
		},
		{
			Key: "release_framer_glass_weather", ID: "22222222-2222-4222-8222-222222222203", ArtistKey: "artist_framer",
			Slug: "glass-weather", Title: "Glass Weather", Type: "album", Status: "published", PriceCents: 790,
			PublishAt: now.AddDate(-2, -1, 0), Genres: []string{"Electronic", "IDM"}, Tags: []string{"fractured", "crystalline", "rhythmic"}, Moods: []string{"clear", "nervous"},
			Description: "Sharp little storms, tuned percussion, and synth lines that behave until they don't.",
			Tracks:      []trackFixture{{"Glass Weather", 206}, {"Facet Drift", 194}, {"Thin Ice Delay", 221}, {"Prism Debt", 188}, {"Rain With Edges", 243}, {"A Clear Error", 211}, {"Last Window Open", 267}},
		},
		{
			Key: "release_framer_signal_bloom", ID: "22222222-2222-4222-8222-222222222204", ArtistKey: "artist_framer",
			Slug: "signal-bloom", Title: "Signal Bloom", Type: "single", Status: "published", PriceCents: 190,
			PublishAt: now.AddDate(0, -5, 0), Genres: []string{"Electronic"}, Tags: []string{"single", "bright", "pulse"}, Moods: []string{"direct", "uplifted"},
			Description: "A concise burst of bright synth pressure and clean percussion.",
			Tracks:      []trackFixture{{"Signal Bloom", 226}},
		},
		{
			Key: "release_framer_untitled_night_work", ID: "22222222-2222-4222-8222-222222222205", ArtistKey: "artist_framer",
			Slug: "untitled-night-work", Title: "Untitled Night Work", Type: "ep", Status: "draft", PriceCents: 490,
			PublishAt: now, Genres: []string{"Electronic"}, Tags: []string{"draft", "sketches"}, Moods: []string{"unfinished", "late"},
			Description: "A private draft Release for upload and lifecycle testing.",
			Tracks:      []trackFixture{{"Working Title One", 192}, {"Blue Screen Morning", 207}, {"Temporary Stars", 181}},
		},
		{
			Key: "release_framer_old_light_archive", ID: "22222222-2222-4222-8222-222222222206", ArtistKey: "artist_framer",
			Slug: "old-light-archive", Title: "Old Light Archive", Type: "compilation", Status: "archived", PriceCents: 590,
			PublishAt: now.AddDate(-3, 0, 0), Genres: []string{"Electronic"}, Tags: []string{"archive", "early"}, Moods: []string{"dusty", "patient"},
			Description: "Older Framer fragments preserved for lifecycle testing, hidden from public listings by Archived status.",
			Tracks:      []trackFixture{{"Old Light", 199}, {"Basement Satellite", 223}, {"Former Signal", 187}, {"Pale Machine", 216}},
		},
		{
			Key: "release_mira_vale_clear_room", ID: "22222222-2222-4222-8222-222222222301", ArtistKey: "artist_mira_vale",
			Slug: "clear-room", Title: "Clear Room", Type: "ep", Status: "published", PriceCents: 390,
			PublishAt: now.AddDate(-1, -3, 0), Genres: []string{"Ambient", "Electronic"}, Tags: []string{"space", "minimal"}, Moods: []string{"calm", "open"},
			Description: "Minimal ambient pieces with enough pulse to keep the room awake.",
			Tracks:      []trackFixture{{"Clear Room", 240}, {"Window Tone", 218}, {"A Quiet Map", 252}},
		},
		{
			Key: "release_mira_vale_blue_hours", ID: "22222222-2222-4222-8222-222222222302", ArtistKey: "artist_mira_vale",
			Slug: "blue-hours", Title: "Blue Hours", Type: "album", Status: "published", PriceCents: 690,
			PublishAt: now.AddDate(-2, 0, 0), Genres: []string{"Ambient Pop"}, Tags: []string{"blue", "vocal", "soft"}, Moods: []string{"tender", "slow"},
			Description: "Soft vocal fragments and wide synths for the hour before the city restarts.",
			Tracks:      []trackFixture{{"Blue Hours", 214}, {"Half Awake", 198}, {"Old Apartment Air", 229}, {"Every Lamp Low", 241}, {"Close Enough to Rain", 260}},
		},
		{
			Key: "release_north_index_pressure_map", ID: "22222222-2222-4222-8222-222222222401", ArtistKey: "artist_north_index",
			Slug: "pressure-map", Title: "Pressure Map", Type: "album", Status: "published", PriceCents: 790,
			PublishAt: now.AddDate(-1, -9, 0), Genres: []string{"Experimental", "Bass"}, Tags: []string{"heavy", "instrumental"}, Moods: []string{"tense", "wide"},
			Description: "Low-end cartography for concrete rooms and stubborn speakers.",
			Tracks:      []trackFixture{{"Pressure Map", 233}, {"Tectonic Sleep", 249}, {"Blackline Survey", 221}, {"Signal Under Stone", 278}, {"Load Bearing Sky", 260}, {"Fault Memory", 245}},
		},
		{
			Key: "release_north_index_dead_reckoner", ID: "22222222-2222-4222-8222-222222222402", ArtistKey: "artist_north_index",
			Slug: "dead-reckoner", Title: "Dead Reckoner", Type: "ep", Status: "published", PriceCents: 490,
			PublishAt: now.AddDate(-2, -6, 0), Genres: []string{"Experimental"}, Tags: []string{"dark", "mechanical"}, Moods: []string{"stark", "driving"},
			Description: "Short mechanical pieces that navigate by damaged instruments.",
			Tracks:      []trackFixture{{"Dead Reckoner", 210}, {"Northing Error", 188}, {"Iron Meridian", 226}, {"Last Known Point", 239}},
		},
	}
}

func userDoc(key, email, username, passwordHash string, meta seedMeta, now time.Time) document {
	return document{
		"_key": key, "username": username, "email": email, "password": passwordHash,
		"emailUnsubscribeToken": fmt.Sprintf("seed-unsubscribe-%s", key),
		"socials":               map[string]interface{}{}, "validatedEmail": true, "role": "user",
		"settings": map[string]interface{}{"receiveEmails": true}, "createdAt": now, "seed": meta,
	}
}

func artistDoc(key, id, name, slug, bio string, meta seedMeta, createdAt, now time.Time) document {
	return document{"_key": key, "id": id, "name": name, "slug": slug, "bio": bio, "createdAt": createdAt, "publishedAt": now, "seed": meta}
}

func releaseDoc(release releaseFixture, meta seedMeta, now time.Time) document {
	basePath := fmt.Sprintf("artist_content/%s/releases/%s/", release.ArtistKey, release.ID)
	trackCount := len(release.Tracks)
	total := 0
	for _, track := range release.Tracks {
		total += track.Duration
	}
	return document{
		"_key": release.Key, "id": release.ID, "slug": release.Slug, "title": release.Title,
		"description": release.Description, "releaseType": release.Type, "status": release.Status,
		"artistKey": release.ArtistKey, "pricing": map[string]interface{}{"basePrice": release.PriceCents, "minimumPrice": release.PriceCents, "trackPrice": 139},
		"metadata":      map[string]interface{}{"genres": release.Genres, "tags": release.Tags, "moods": release.Moods},
		"credits":       map[string]interface{}{"artists": []interface{}{}, "credits": []interface{}{}, "writers": []string{}, "copyrightHolder": "Silkwave Seed", "copyrightYear": release.PublishAt.Year()},
		"assets":        map[string]interface{}{"basePath": basePath, "coverArt": coverArt(release.ArtistKey, release.ID, release.Slug)},
		"schedule":      map[string]interface{}{"createdAt": release.PublishAt, "updatedAt": now, "publishAt": release.PublishAt, "publishedAt": release.PublishAt},
		"totalDuration": total, "trackCount": trackCount, "availableFormats": []string{"mp3_320", "flac", "wav"},
		"trackList": trackTitles(release.Tracks), "license": "all_rights_reserved", "isExplicit": false,
		"isFeatured": release.ArtistKey == "artist_framer" && release.Status == "published", "downloadEnabled": true,
		"streamingEnabled": true, "publishAt": release.PublishAt, "cover": coverArt(release.ArtistKey, release.ID, release.Slug)["medium"],
		"genres": release.Genres, "createdAt": release.PublishAt, "updatedAt": now, "seed": meta,
	}
}

func trackDoc(release releaseFixture, track trackFixture, order int, meta seedMeta, now time.Time) document {
	slug := slugifyTitle(track.Title)
	id := fmt.Sprintf("33333333-3333-4333-8333-%012d", trackIDNumber(release.ID, order))
	originalPath := fmt.Sprintf("artist_content/%s/releases/%s/flacs/%02d-%s.flac", release.ArtistKey, release.ID, order, slug)
	previewPath := fmt.Sprintf("artist_content/%s/releases/%s/previews/%02d-%s.mp3", release.ArtistKey, release.ID, order, slug)
	return document{
		"_key": fmt.Sprintf("track_%s_%02d", strings.TrimPrefix(release.Key, "release_"), order),
		"id":   id, "releaseKey": release.Key, "title": track.Title, "order": order,
		"duration": track.Duration, "durationDisplay": durationDisplay(track.Duration),
		"files": map[string]interface{}{
			"original":    map[string]interface{}{"path": originalPath, "format": "flac", "fileSize": 42000000 + order*137000, "bitrate": 1411, "sampleRate": 44100},
			"transcoded":  []interface{}{map[string]interface{}{"format": "mp3_320", "path": fmt.Sprintf("artist_content/%s/releases/%s/mp3s/%02d-%s.mp3", release.ArtistKey, release.ID, order, slug), "fileSize": 9000000 + order*81000, "bitrate": 320}},
			"previewPath": previewPath,
		},
		"metadata":        map[string]interface{}{"genres": release.Genres, "tags": release.Tags},
		"featuredArtists": []interface{}{}, "writers": []string{}, "uploaded": true, "transcodingComplete": true,
		"createdAt": release.PublishAt, "updatedAt": now, "isExplicit": false, "streamingEnabled": true,
		"downloadEnabled": true, "playCount": 0, "downloadCount": 0, "seed": meta,
	}
}

func edge(key, from, to string, meta seedMeta, now time.Time) document {
	return document{"_key": key, "_from": from, "_to": to, "createdAt": now, "seed": meta}
}

func subscriptionDoc(key, artistKey, name, description string, price float64, currency string, discount float64, meta seedMeta, now time.Time) document {
	return document{"_key": key, "artistKey": artistKey, "name": name, "description": description, "price": price, "currency": currency, "subscriberDiscountPercent": discount, "createdAt": now, "seed": meta}
}

func subscriberDoc(key, artistKey, subscriberKey, subscriptionKey string, meta seedMeta, now time.Time) document {
	return document{"_key": key, "artistKey": artistKey, "subscriberKey": subscriberKey, "subscriptionKey": subscriptionKey, "status": "active", "createdAt": now, "seed": meta}
}

func coverArt(artistKey, releaseID, slug string) map[string]string {
	base := fmt.Sprintf("artist_content/%s/releases/%s/cover", artistKey, releaseID)
	return map[string]string{"original": fmt.Sprintf("%s/%s-original.jpg", base, slug), "thumbnail": fmt.Sprintf("%s/%s-300.jpg", base, slug), "medium": fmt.Sprintf("%s/%s-600.jpg", base, slug)}
}

func trackTitles(tracks []trackFixture) []string {
	titles := make([]string, 0, len(tracks))
	for _, track := range tracks {
		titles = append(titles, track.Title)
	}
	return titles
}

func durationDisplay(seconds int) string {
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}

func slugifyTitle(value string) string {
	value = strings.ToLower(value)
	var out strings.Builder
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			out.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}

func trackIDNumber(releaseID string, order int) int {
	parts := strings.Split(releaseID, "222222222")
	if len(parts) != 2 {
		return order
	}
	var base int
	_, _ = fmt.Sscanf(parts[1], "%d", &base)
	return base*100 + order
}
