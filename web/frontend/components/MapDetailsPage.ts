import { BasePage } from './BasePage';

/**
 * Main application initialization
 */
class MapDetailsPage extends BasePage {
    private currentMapId: string | null = null;
    private isLoadingMap: boolean = false; // Loading state

    constructor() {
        super();
        this.initializeSpecificComponents();
        this.bindSpecificEvents();
        this.loadInitialState();
    }

    protected initializeSpecificComponents(): void {
        const mapIdInput = document.getElementById("mapIdInput") as HTMLInputElement | null;
        const mapId = mapIdInput?.value.trim() || null; // Allow null if input not found/empty

        console.log('MapDetailsPage application initialized');
    }

    protected bindSpecificEvents(): void {
        const mobileMenuButton = document.getElementById('mobile-menu-button');
        if (mobileMenuButton) {
            mobileMenuButton.addEventListener('click', () => {
              // Do things like sidebar drawers etc
            });
        }

        const saveButton = document.querySelector('header button.bg-blue-600');
        if (saveButton) {
            saveButton.addEventListener('click', this.saveDocument.bind(this));
        }

        const exportButton = document.querySelector('header button.bg-gray-200');
        if (exportButton) {
            exportButton.addEventListener('click', this.exportDocument.bind(this));
        }
    }

    /** Load document data and set initial UI states */
    private loadInitialState(): void {
        // Theme button state is handled by BasePage

        const mapIdInput = document.getElementById("mapIdInput") as HTMLInputElement | null;
        const mapId = mapIdInput?.value.trim() || null;

        if (mapId) {
            this.currentMapId = mapId;
            console.log(`Found Map ID: ${this.currentMapId}. Loading data...`);
            this.loadMapData(this.currentMapId);
        } else {
            console.error("Map ID input element not found or has no value. Cannot load document.");
            this.showToast("Error", "Could not load document: Map ID missing.", "error");
        }
    }

    /**
     * Fetches map metadata, initializes section shells, and triggers content loading for each section.
     */
    private async loadMapData(mapId: string): Promise<void> {
        // TODO: Show global loading indicator
        console.log(`MapDetailsPage: Loading map ${mapId}...`);

        // here is where we would do "reload" via ajax - this coul dbe via ajax or via htmx
    }

    // Theme management is handled by BasePage

    /** Save document (Placeholder - needs full implementation later) */
    private saveDocument(): void {
        console.log("Save button clicked (Placeholder - Requires API integration for full save)");
        // This full save logic will be replaced by incremental saves triggered by component callbacks
        this.showToast('Save Action', 'Incremental saves handle updates. Full save TBD.', 'info');
    }

    /** Export document (Placeholder) */
    private exportDocument(): void {
        this.showToast('Export started', 'Your document is being prepared for export.', 'info');
        setTimeout(() => {
            this.showToast('Export complete', 'Document export simulation finished.', 'success');
        }, 1500);
    }

    public destroy(): void {
        // Clean up any specific resources for MapDetailsPage
        // Currently no specific cleanup needed
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const lc = new MapDetailsPage();
});
