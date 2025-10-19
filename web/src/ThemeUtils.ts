import { ITheme } from '../assets/themes/BaseTheme';

/**
 * ThemeUtils - Utilities for working with themes in the UI
 */
export class ThemeUtils {
    /**
     * Hydrate theme images after Go template renders HTML
     *
     * Finds all elements with class 'theme-unit-image' or 'theme-tile-image'
     * and populates them using the theme's setUnitImage/setTileImage methods.
     *
     * This bridges Go template rendering with TypeScript theme asset loading:
     * - Go renders HTML with data attributes (data-unit-id, data-player-id)
     * - This method reads those attributes and calls theme methods to load assets
     * - Theme handles the actual SVG/PNG loading and color application
     *
     * @param rootElement - The container element to search for theme images
     * @param theme - The theme instance to use for loading images
     * @param debugMode - Enable debug logging
     * @returns Promise that resolves when all images are hydrated
     */
    public static async hydrateThemeImages(
        rootElement: HTMLElement,
        theme: ITheme | null,
        debugMode: boolean = false
    ): Promise<void> {
        if (!theme) {
            if (debugMode) {
                console.log('[ThemeUtils] No theme set, skipping image hydration');
            }
            return;
        }

        const promises: Promise<void>[] = [];

        // Hydrate unit images
        const unitImages = rootElement.querySelectorAll<HTMLImageElement>('.theme-unit-image');
        unitImages.forEach(img => {
            const unitId = parseInt(img.dataset.unitId || '0', 10);
            const playerId = parseInt(img.dataset.playerId || '0', 10);

            if (unitId && img.parentElement) {
                // Use theme to set the image (async)
                promises.push(theme.setUnitImage(unitId, playerId, img.parentElement));
            }
        });

        // Hydrate tile images
        const tileImages = rootElement.querySelectorAll<HTMLImageElement>('.theme-tile-image');
        tileImages.forEach(img => {
            const tileId = parseInt(img.dataset.tileId || '0', 10);
            const playerId = parseInt(img.dataset.playerId || '0', 10);

            if (tileId && img.parentElement) {
                // Use theme to set the image (async)
                promises.push(theme.setTileImage(tileId, playerId, img.parentElement));
            }
        });

        // Wait for all images to load
        await Promise.all(promises);

        if (debugMode) {
            console.log(`[ThemeUtils] Hydrated ${unitImages.length} unit images and ${tileImages.length} tile images`);
        }
    }
}
