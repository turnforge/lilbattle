/**
 * Asset Theme Preference Manager
 *
 * Manages user's asset theme preference (fantasy, classic, etc.)
 * Stores preference in both localStorage (frontend) and cookie (backend sync)
 */
export class AssetThemePreference {
    private static readonly STORAGE_KEY = 'assetTheme';
    private static readonly COOKIE_NAME = 'assetTheme';
    private static readonly COOKIE_MAX_AGE = 365 * 24 * 60 * 60; // 1 year in seconds
    private static readonly DEFAULT_THEME = 'fantasy';

    /**
     * Get the current asset theme preference
     * Priority: localStorage > default
     */
    public static get(): string {
        return localStorage.getItem(AssetThemePreference.STORAGE_KEY) || AssetThemePreference.DEFAULT_THEME;
    }

    /**
     * Save asset theme preference to both localStorage and cookie
     */
    public static set(theme: string): void {
        // Save to localStorage for frontend
        localStorage.setItem(AssetThemePreference.STORAGE_KEY, theme);

        // Save to cookie for backend
        const expires = new Date();
        expires.setTime(expires.getTime() + AssetThemePreference.COOKIE_MAX_AGE * 1000);
        document.cookie = `${AssetThemePreference.COOKIE_NAME}=${theme}; expires=${expires.toUTCString()}; path=/; SameSite=Strict`;
    }

    /**
     * Clear asset theme preference (will use default)
     */
    public static clear(): void {
        localStorage.removeItem(AssetThemePreference.STORAGE_KEY);
        document.cookie = `${AssetThemePreference.COOKIE_NAME}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; SameSite=Strict`;
    }

    /**
     * Get the default theme
     */
    public static getDefault(): string {
        return AssetThemePreference.DEFAULT_THEME;
    }
}
