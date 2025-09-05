import { BaseAssetProvider } from './AssetProvider';
import { AllowedUnitIDs } from '../ColorsAndNames';

/**
 * Player color definitions for SVG templating
 */
interface PlayerColors {
    primary: string;
    secondary: string;
    accent?: string;
}

/**
 * Theme mapping structure from mapping.json
 */
interface ThemeMapping {
    units: Record<string, {
        old: string;
        name: string;
        image: string;
    }>;
    terrains: Record<string, {
        old: string;
        name: string;
        image: string;
    }>;
}

/**
 * Asset provider that uses SVG templates with color replacement
 * Loads base SVG templates and generates color variations for each player
 */
export class TemplateSVGAssetProvider extends BaseAssetProvider {
    private themeName: string;
    private themePath: string;
    private rasterSize: number;
    private fallbackToPNG: boolean;
    private debugMode: boolean;
    private processedAssets: Set<string> = new Set();
    private themeMapping: ThemeMapping | null = null;
    
    // Player color schemes
    private playerColors: PlayerColors[] = [
        { primary: '#808080', secondary: '#606060', accent: '#404040' }, // Player 0 - Neutral gray
        { primary: '#ff4444', secondary: '#cc0000', accent: '#ffaaaa' }, // Player 1 - Red
        { primary: '#4444ff', secondary: '#0000cc', accent: '#aaaaff' }, // Player 2 - Blue
        { primary: '#44ff44', secondary: '#00cc00', accent: '#aaffaa' }, // Player 3 - Green
        { primary: '#ffff44', secondary: '#cccc00', accent: '#ffffaa' }, // Player 4 - Yellow
        { primary: '#ff44ff', secondary: '#cc00cc', accent: '#ffaaff' }, // Player 5 - Magenta
        { primary: '#44ffff', secondary: '#00cccc', accent: '#aaffff' }, // Player 6 - Cyan
        { primary: '#ff8844', secondary: '#cc6600', accent: '#ffccaa' }, // Player 7 - Orange
        { primary: '#8844ff', secondary: '#6600cc', accent: '#ccaaff' }, // Player 8 - Purple
        { primary: '#88ff44', secondary: '#66cc00', accent: '#ccffaa' }, // Player 9 - Lime
        { primary: '#ff4488', secondary: '#cc0066', accent: '#ffaacc' }, // Player 10 - Pink
        { primary: '#44ff88', secondary: '#00cc66', accent: '#aaffcc' }, // Player 11 - Teal
        { primary: '#8888ff', secondary: '#6666cc', accent: '#ccccff' }, // Player 12 - Light Blue
    ];
    
    constructor(themeName: string = 'fantasy', rasterSize: number = 160, fallbackToPNG: boolean = true, debugMode: boolean = true) {
        super();
        this.themeName = themeName;
        this.themePath = `/static/assets/themes/${themeName}/`;
        this.rasterSize = rasterSize;
        this.fallbackToPNG = fallbackToPNG;
        this.debugMode = debugMode;
        this.assetSize = { width: rasterSize, height: rasterSize };
    }
    
    async preloadAssets(): Promise<void> {
        if (!this.loader) {
            console.error('[TemplateSVGAssetProvider] Loader not configured');
            return;
        }
        
        console.log(`[TemplateSVGAssetProvider] Loading theme: ${this.themeName}`);
        
        // Load the mapping.json file using Phaser's loader
        // Add timestamp to force reload during development
        const timestamp = this.debugMode ? `?t=${Date.now()}` : '';
        const mappingPath = `${this.themePath}mapping.json${timestamp}`;
        console.log(`[TemplateSVGAssetProvider] Queuing mapping.json for loading: ${mappingPath}`);
        
        // Use Phaser's JSON loader to load the mapping
        this.loader.json('themeMapping', mappingPath);
        
        // Set up a one-time handler for when the mapping loads
        this.loader.once('filecomplete-json-themeMapping', (key: string, type: string, data: any) => {
            console.log(`[TemplateSVGAssetProvider] Mapping loaded via Phaser`);
            this.themeMapping = data as ThemeMapping;
            
            if (this.themeMapping) {
                console.log(`[TemplateSVGAssetProvider] Mapping contains: ${Object.keys(this.themeMapping.terrains).length} terrains, ${Object.keys(this.themeMapping.units).length} units`);
                
                // Now queue all SVG files based on the mapping
                // These will be added to the loader queue dynamically
                this.loadMappedTerrainTemplates();
                this.loadMappedUnitTemplates();
                
                const totalAssets = Object.keys(this.themeMapping.terrains).length + 
                                   Object.keys(this.themeMapping.units).length;
                console.log(`[TemplateSVGAssetProvider] Queued ${totalAssets} SVG files for loading`);
            }
        });
        
        // Set up error handling
        this.loader.on('loaderror', (file: any) => {
            console.warn(`[TemplateSVGAssetProvider] Failed to load: ${file.key} from ${file.url}`);
            
            if (this.fallbackToPNG && file.key.includes('_template')) {
                // Try to load PNG fallback
                this.loadPNGFallback(file.key);
            }
        });
    }
    
    private loadMappedTerrainTemplates(): void {
        if (!this.themeMapping) {
            console.error('[TemplateSVGAssetProvider] No theme mapping loaded');
            return;
        }
        
        // Load terrain templates based on mapping
        Object.entries(this.themeMapping.terrains).forEach(([id, terrain]) => {
            // The image path in mapping is relative to theme folder (e.g., "Tiles/Castle.svg")
            const svgPath = `${this.themePath}${terrain.image}`;
            const templateKey = `terrain_${id}_template`;
            
            // Load the SVG as text for template processing
            this.loader.text(templateKey, svgPath);
            
            console.log(`[TemplateSVGAssetProvider] Loading terrain ${id}: ${terrain.name} from ${svgPath}`);
        });
    }
    
    private loadMappedUnitTemplates(): void {
        if (!this.themeMapping) {
            console.error('[TemplateSVGAssetProvider] No theme mapping loaded');
            return;
        }
        
        // Load unit templates based on mapping
        Object.entries(this.themeMapping.units).forEach(([id, unit]) => {
            // Skip units with empty image paths
            if (!unit.image || unit.image === 'Units/.svg') {
                console.warn(`[TemplateSVGAssetProvider] Skipping unit ${id}: ${unit.name} (no image)`);
                return;
            }
            
            // The image path in mapping is relative to theme folder (e.g., "Units/Knight.svg")
            const svgPath = `${this.themePath}${unit.image}`;
            const templateKey = `unit_${id}_template`;
            
            // Load the SVG as text for template processing
            this.loader.text(templateKey, svgPath);
            
            console.log(`[TemplateSVGAssetProvider] Loading unit ${id}: ${unit.name} from ${svgPath}`);
        });
    }
    
    private loadPNGFallback(templateKey: string): void {
        // Extract type from template key
        const match = templateKey.match(/(terrain|unit)_(\d+)_template/);
        if (!match) return;
        
        const assetType = match[1];
        const typeId = match[2];
        
        // Load PNG versions for all player colors
        for (let color = 0; color <= this.maxPlayers; color++) {
            const pngPath = assetType === 'terrain' 
                ? `/static/assets/v1/Tiles/${typeId}/${color}.png`
                : `/static/assets/v1/Units/${typeId}/${color}.png`;
            
            const textureKey = `${assetType}_${typeId}_${color}`;
            this.loader.image(textureKey, pngPath);
        }
    }
    
    protected onLoadComplete(): void {
        // Don't mark as ready yet - we need to post-process
        // The ready flag will be set after postProcessAssets() completes
        console.log('[TemplateSVGAssetProvider] Assets loaded, awaiting post-processing');
    }
    
    async postProcessAssets(): Promise<void> {
        console.log('[TemplateSVGAssetProvider] Starting post-processing of SVG templates');
        
        // If mapping wasn't set during loading, try to get it from cache
        if (!this.themeMapping) {
            const cachedMapping = this.scene.cache.json.get('themeMapping');
            if (cachedMapping) {
                this.themeMapping = cachedMapping as ThemeMapping;
                console.log('[TemplateSVGAssetProvider] Retrieved mapping from Phaser cache');
            }
        }
        
        console.log('[TemplateSVGAssetProvider] themeMapping status:', this.themeMapping ? 'loaded' : 'null');
        
        if (!this.themeMapping) {
            console.error('[TemplateSVGAssetProvider] No theme mapping available for post-processing');
            console.error('[TemplateSVGAssetProvider] This should not happen - preloadAssets should have loaded the mapping');
            return;
        }
        
        const promises: Promise<void>[] = [];
        
        // Process terrain templates based on mapping
        Object.entries(this.themeMapping.terrains).forEach(([id, terrain]) => {
            const templateKey = `terrain_${id}_template`;
            const svgTemplate = this.scene.cache.text.get(templateKey);
            
            if (svgTemplate) {
                const typeId = parseInt(id);
                // For nature terrains, only create neutral variant
                const maxColor = this.natureTerrains.includes(typeId) ? 0 : this.maxPlayers;
                
                console.log(`[TemplateSVGAssetProvider] Processing terrain ${id} (${terrain.name}), nature=${this.natureTerrains.includes(typeId)}, maxColor=${maxColor}`);
                
                for (let player = 0; player <= maxColor; player++) {
                    const textureKey = `terrain_${id}_${player}`;
                    promises.push(this.createColorVariant(svgTemplate, textureKey, player));
                    
                    // For nature terrains, create aliases for all players
                    if (this.natureTerrains.includes(typeId) && player === 0) {
                        for (let p = 1; p <= this.maxPlayers; p++) {
                            this.processedAssets.add(`terrain_${id}_${p}`);
                        }
                    }
                }
            } else {
                console.warn(`[TemplateSVGAssetProvider] No SVG template found for terrain ${id} (${terrain.name})`);
            }
        });
        
        // Process unit templates based on mapping
        Object.entries(this.themeMapping.units).forEach(([id, unit]) => {
            // Skip units with no image
            if (!unit.image || unit.image === 'Units/.svg') {
                return;
            }
            
            const templateKey = `unit_${id}_template`;
            const svgTemplate = this.scene.cache.text.get(templateKey);
            
            if (svgTemplate) {
                for (let player = 0; player <= this.maxPlayers; player++) {
                    const textureKey = `unit_${id}_${player}`;
                    promises.push(this.createColorVariant(svgTemplate, textureKey, player));
                }
            }
        });
        
        // Wait for all processing to complete
        console.log(`[TemplateSVGAssetProvider] Waiting for ${promises.length} texture creation promises...`);
        await Promise.all(promises);
        
        // Create aliases for nature terrains
        this.createNatureTerrainAliases();
        
        console.log(`[TemplateSVGAssetProvider] Post-processing complete. Processed ${this.processedAssets.size} assets`);
        
        // List some created textures for debugging
        const sampleTextures = ['terrain_23_0', 'terrain_26_0', 'terrain_1_0', 'terrain_5_0'];
        sampleTextures.forEach(key => {
            if (this.scene.textures.exists(key)) {
                console.log(`[TemplateSVGAssetProvider] ✓ Texture exists: ${key}`);
            } else {
                console.log(`[TemplateSVGAssetProvider] ✗ Texture missing: ${key}`);
            }
        });
        
        // Now mark as ready
        this.ready = true;
        if (this.onComplete) {
            this.onComplete();
        }
    }
    
    private async createColorVariant(
        svgTemplate: string,
        textureKey: string,
        player: number
    ): Promise<void> {
        const colors = this.playerColors[player] || this.playerColors[0];
        
        // Check if the SVG has template variables
        const hasTemplateVars = svgTemplate.includes('{{');
        
        let processedSVG = svgTemplate;
        
        if (hasTemplateVars) {
            // Replace color placeholders in SVG
            processedSVG = svgTemplate
                .replace(/\{\{PRIMARY_COLOR\}\}/g, colors.primary)
                .replace(/\{\{SECONDARY_COLOR\}\}/g, colors.secondary)
                .replace(/\{\{ACCENT_COLOR\}\}/g, colors.accent || colors.primary)
                .replace(/\{\{PLAYER_ID\}\}/g, player.toString())
                .replace(/\{\{PLAYER_NUMBER\}\}/g, (player + 1).toString());
        }
        // For SVGs without templates, we'll apply tinting after rendering
        
        // Add gradient definitions if not present
        if (processedSVG.includes('{{PLAYER_GRADIENT}}')) {
            const gradientDef = `
                <defs>
                    <linearGradient id="playerGradient${player}" x1="0%" y1="0%" x2="100%" y2="100%">
                        <stop offset="0%" style="stop-color:${colors.primary};stop-opacity:1" />
                        <stop offset="100%" style="stop-color:${colors.secondary};stop-opacity:1" />
                    </linearGradient>
                    <radialGradient id="playerRadialGradient${player}">
                        <stop offset="0%" style="stop-color:${colors.primary};stop-opacity:1" />
                        <stop offset="100%" style="stop-color:${colors.secondary};stop-opacity:0.8" />
                    </radialGradient>
                </defs>
            `;
            
            processedSVG = processedSVG
                .replace('</svg>', `${gradientDef}</svg>`)
                .replace(/\{\{PLAYER_GRADIENT\}\}/g, `url(#playerGradient${player})`)
                .replace(/\{\{PLAYER_RADIAL_GRADIENT\}\}/g, `url(#playerRadialGradient${player})`);
        }
        
        // Convert processed SVG to texture (with optional tinting for player colors)
        await this.svgToTexture(processedSVG, textureKey, player, !hasTemplateVars);
        this.processedAssets.add(textureKey);
    }
    
    private async svgToTexture(svgString: string, textureKey: string, player: number = 0, applyTint: boolean = false): Promise<void> {
        return new Promise((resolve, reject) => {
            try {
                // Create blob from SVG string
                const blob = new Blob([svgString], { type: 'image/svg+xml;charset=utf-8' });
                const url = URL.createObjectURL(blob);
                
                // Create image element
                const img = new Image();
                
                img.onload = () => {
                    try {
                        // Create canvas and draw SVG at desired resolution
                        const canvas = document.createElement('canvas');
                        canvas.width = this.rasterSize;
                        canvas.height = this.rasterSize;
                        const ctx = canvas.getContext('2d');
                        
                        if (!ctx) {
                            throw new Error('Failed to get canvas context');
                        }
                        
                        // Enable image smoothing for better quality
                        ctx.imageSmoothingEnabled = true;
                        ctx.imageSmoothingQuality = 'high';
                        
                        // Draw the SVG to canvas at the target size
                        ctx.drawImage(img, 0, 0, this.rasterSize, this.rasterSize);
                        
                        // Apply player color tinting if needed (for city terrains without template vars)
                        if (applyTint && player > 0) {
                            // Get the terrain ID from the texture key (e.g., "terrain_1_2" -> "1")
                            const match = textureKey.match(/terrain_(\d+)_/);
                            if (match) {
                                const terrainId = parseInt(match[1]);
                                // Only apply tint to city terrains
                                if (this.cityTerrains.includes(terrainId)) {
                                    const colors = this.playerColors[player];
                                    // Apply color multiply effect
                                    ctx.globalCompositeOperation = 'multiply';
                                    ctx.fillStyle = colors.primary;
                                    ctx.fillRect(0, 0, this.rasterSize, this.rasterSize);
                                    ctx.globalCompositeOperation = 'destination-in';
                                    ctx.drawImage(canvas, 0, 0);
                                }
                            }
                        }
                        
                        // Add to Phaser's texture manager
                        this.scene.textures.addCanvas(textureKey, canvas);
                        
                        // Verify the texture was actually created
                        if (this.scene.textures.exists(textureKey)) {
                            console.log(`[TemplateSVGAssetProvider] Successfully created texture: ${textureKey}`);
                        } else {
                            console.error(`[TemplateSVGAssetProvider] Failed to create texture: ${textureKey}`);
                        }
                        
                        // Clean up
                        URL.revokeObjectURL(url);
                        resolve();
                    } catch (error) {
                        console.error(`[TemplateSVGAssetProvider] Error processing ${textureKey}:`, error);
                        URL.revokeObjectURL(url);
                        reject(error);
                    }
                };
                
                img.onerror = () => {
                    console.error(`[TemplateSVGAssetProvider] Failed to load SVG for ${textureKey}`);
                    URL.revokeObjectURL(url);
                    reject(new Error(`Failed to load SVG for ${textureKey}`));
                };
                
                img.src = url;
            } catch (error) {
                console.error(`[TemplateSVGAssetProvider] Error creating texture ${textureKey}:`, error);
                reject(error);
            }
        });
    }
    
    private createNatureTerrainAliases(): void {
        // For nature terrains, we'll handle the aliasing in the getTerrainTexture method
        // Since Phaser doesn't support addTexture, we can't create true aliases
        // Instead, getTerrainTexture will return the base texture for all players
        console.log('[TemplateSVGAssetProvider] Nature terrain aliases will be handled via getTerrainTexture');
    }
    
    getTerrainTexture(tileType: number, player: number): string {
        // Check if this terrain exists in the mapping
        if (this.themeMapping && !this.themeMapping.terrains[tileType.toString()]) {
            // Terrain not in theme, return a fallback or null
            console.warn(`[TemplateSVGAssetProvider] Terrain ${tileType} not found in theme ${this.themeName}`);
            return `terrain_${tileType}_${player}`; // Return expected key anyway for fallback
        }
        
        const textureKey = `terrain_${tileType}_${player}`;
        
        // For nature terrains, always use the neutral texture
        if (this.natureTerrains.includes(tileType)) {
            return `terrain_${tileType}_0`;
        }
        
        return textureKey;
    }
    
    getUnitTexture(unitType: number, player: number): string {
        // Check if this unit exists in the mapping
        if (this.themeMapping && !this.themeMapping.units[unitType.toString()]) {
            // Unit not in theme, return a fallback or null
            console.warn(`[TemplateSVGAssetProvider] Unit ${unitType} not found in theme ${this.themeName}`);
            return `unit_${unitType}_${player}`; // Return expected key anyway for fallback
        }
        
        return `unit_${unitType}_${player}`;
    }
    
    dispose(): void {
        // Clean up processed textures
        this.processedAssets.forEach(key => {
            if (this.scene.textures.exists(key)) {
                this.scene.textures.remove(key);
            }
        });
        
        this.processedAssets.clear();
        super.dispose();
    }
}
