import fs from 'fs';
import path from 'path';

// Base directory containing icon folders (e.g., "filled", "outline", etc.)
const baseDir = process.argv[2] || 'C:/Users/brand/Pictures/svgs';
// Output directory for JSON files
const outputDir = process.argv[3] || 'src/assets';

/**
 * Extract SVG paths from a single SVG file
 */
const extractPathsFromSvg = (filePath: string): string | null => {
    const content = fs.readFileSync(filePath, 'utf-8');

    // Extract all path elements (handles both self-closing and regular paths)
    const pathMatches = content.match(/<path[^>]*(?:\/>|>.*?<\/path>)/gs);

    if (pathMatches) {
        // Clean up the paths: remove fill color placeholders, normalize whitespace
        const cleanedPaths = pathMatches
            .map((p: string) => {
                return p
                    .replace(/fill="#[0-9A-Fa-f]{6}"/g, '')
                    .replace(/fill="currentColor"/g, '')
                    .replace(/\s+/g, ' ')
                    .trim();
            })
            .join('');

        return cleanedPaths;
    }

    return null;
};

/**
 * Process a single folder of SVGs and create a JSON file
 */
const processSvgFolder = (folderPath: string, folderName: string): void => {
    const icons: Record<string, string> = {};

    // Read all SVG files in the folder
    const files = fs.readdirSync(folderPath).filter((f: string) => f.endsWith('.svg'));

    if (files.length === 0) {
        console.log(`⚠️  No SVG files found in: ${folderName}`);
        return;
    }

    files.forEach((file: string) => {
        const filePath = path.join(folderPath, file);

        // Extract icon name from filename (remove .svg, keep lowercase)
        const iconName = file.replace('.svg', '').toLowerCase();

        const paths = extractPathsFromSvg(filePath);
        if (paths) {
            icons[iconName] = paths;
        } else {
            console.log(`⚠️  Could not extract paths from: ${file}`);
        }
    });

    // Sort icons alphabetically
    const sortedIcons: Record<string, string> = {};
    Object.keys(icons)
        .sort()
        .forEach((key) => {
            sortedIcons[key] = icons[key];
        });

    // Ensure output directory exists
    if (!fs.existsSync(outputDir)) {
        fs.mkdirSync(outputDir, { recursive: true });
    }

    // Write to JSON file named after the folder
    const outputFile = path.join(outputDir, `${folderName}.json`);
    fs.writeFileSync(outputFile, JSON.stringify(sortedIcons, null, '\t'));

    console.log(`✅ Created ${folderName}.json with ${Object.keys(sortedIcons).length} icons`);
};

/**
 * Main function - processes all subfolders in the base directory
 */
const main = (): void => {
    console.log(`\n📁 Scanning: ${baseDir}`);
    console.log(`📤 Output to: ${outputDir}\n`);

    if (!fs.existsSync(baseDir)) {
        console.error(`❌ Directory not found: ${baseDir}`);
        process.exit(1);
    }

    // Get all items in the base directory
    const items = fs.readdirSync(baseDir);

    // Filter for directories only
    const folders = items.filter((item: string) => {
        const itemPath = path.join(baseDir, item);
        return fs.statSync(itemPath).isDirectory();
    });

    if (folders.length === 0) {
        // No subfolders - treat the baseDir itself as the icon folder
        const folderName = path.basename(baseDir);
        console.log(`📂 Processing single folder: ${folderName}`);
        processSvgFolder(baseDir, folderName);
    } else {
        // Process each subfolder
        console.log(`Found ${folders.length} folder(s) to process:\n`);

        folders.forEach((folder: string) => {
            const folderPath = path.join(baseDir, folder);
            console.log(`📂 Processing: ${folder}`);
            processSvgFolder(folderPath, folder);
        });
    }

    console.log('\n✨ Done!\n');
};

main();
