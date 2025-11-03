import { ThemeManager } from '../lib/ThemeManager';
import { SplashScreen } from '../lib/SplashScreen';

/**
 * Manages the game listing page logic
 */
class GameListingPage {
    constructor() {
        ThemeManager.init();
        this.init();
    }

    /**
     * Initialize page
     */
    private init(): void {
        // Dismiss splash screen once page is ready
        SplashScreen.dismiss();
    }
}

// Initialize the GameListingPage when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
    (window as any).Page = new GameListingPage();
});
