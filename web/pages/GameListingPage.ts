import { ThemeManager } from '../lib/ThemeManager';
import { SplashScreen } from '../lib/SplashScreen';
import { BasePage } from '../lib/BasePage';

/**
 * Manages the game listing page logic
 */
class GameListingPage extends BasePage {}
GameListingPage.loadAfterPageLoaded("gameListingPage", GameListingPage, "GameListingPage")
