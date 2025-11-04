import { BaseComponent } from './Component';
import { LCMComponent } from './LCMComponent';
import { EventBus } from './EventBus';

/**
 * MobileBottomDrawer - Reusable bottom drawer component for mobile layouts
 *
 * Features:
 * - Slides up from bottom covering 60-70% of viewport
 * - Backdrop overlay that dims the content behind
 * - Auto-closes when backdrop is tapped
 * - Smooth slide-up/down animations
 * - Holds any panel content
 */
export class MobileBottomDrawer extends BaseComponent implements LCMComponent {
    private backdropElement: HTMLElement;
    private drawerElement: HTMLElement;
    private contentElement: HTMLElement;
    private closeButton: HTMLElement | null;
    private isOpen: boolean = false;
    private onCloseCallback?: () => void;

    /**
     * Create a MobileBottomDrawer
     * @param rootElement - The root container element for the drawer
     * @param eventBus - Event bus for component communication
     * @param debugMode - Enable debug logging
     */
    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('mobile-bottom-drawer', rootElement, eventBus, debugMode);
    }

    /**
     * Phase 1: Initialize DOM and discover child components
     */
    async performLocalInit(): Promise<LCMComponent[]> {
        // Find drawer elements
        this.backdropElement = this.rootElement.querySelector('.drawer-backdrop') as HTMLElement;
        this.drawerElement = this.rootElement.querySelector('.drawer-container') as HTMLElement;
        this.contentElement = this.rootElement.querySelector('.drawer-content') as HTMLElement;
        this.closeButton = this.rootElement.querySelector('.drawer-close-btn') as HTMLElement;

        if (!this.backdropElement || !this.drawerElement || !this.contentElement) {
            throw new Error('MobileBottomDrawer: Required drawer elements not found');
        }

        // Bind event listeners
        this.bindEvents();

        return [];
    }

    /**
     * Bind event listeners for drawer interactions
     */
    private bindEvents(): void {
        // Close drawer on backdrop click
        this.backdropElement.addEventListener('click', (e) => {
            if (e.target === this.backdropElement) {
                this.close();
            }
        });

        // Close drawer on close button click
        if (this.closeButton) {
            this.closeButton.addEventListener('click', () => {
                this.close();
            });
        }

        // Close drawer on Escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isOpen) {
                this.close();
            }
        });
    }

    /**
     * Open the drawer with slide-up animation
     */
    public open(): void {
        if (this.isOpen) return;

        this.isOpen = true;

        // Show backdrop and drawer
        this.rootElement.classList.remove('hidden');

        // Force reflow to ensure transition works
        this.rootElement.offsetHeight;

        // Add visible class for backdrop fade-in
        this.backdropElement.classList.add('backdrop-visible');

        // Slide drawer up
        this.drawerElement.classList.remove('translate-y-full');

        // Emit event
        this.eventBus.emit('drawer-opened', { drawerId: this.componentId }, null, this);
    }

    /**
     * Close the drawer with slide-down animation
     */
    public close(): void {
        if (!this.isOpen) return;

        this.isOpen = false;

        // Fade out backdrop
        this.backdropElement.classList.remove('backdrop-visible');

        // Slide drawer down
        this.drawerElement.classList.add('translate-y-full');

        // Hide after animation completes (300ms)
        setTimeout(() => {
            if (!this.isOpen) {
                this.rootElement.classList.add('hidden');
            }
        }, 300);

        // Emit event
        this.eventBus.emit('drawer-closed', { drawerId: this.componentId }, null, this);

        // Call close callback if provided
        if (this.onCloseCallback) {
            this.onCloseCallback();
        }
    }

    /**
     * Toggle drawer open/closed
     */
    public toggle(): void {
        if (this.isOpen) {
            this.close();
        } else {
            this.open();
        }
    }

    /**
     * Check if drawer is currently open
     */
    public getIsOpen(): boolean {
        return this.isOpen;
    }

    /**
     * Set the content element for the drawer
     * @param element - The element to insert into the drawer content area
     */
    public setContent(element: HTMLElement): void {
        this.contentElement.innerHTML = '';
        this.contentElement.appendChild(element);
    }

    /**
     * Set a callback to be called when drawer closes
     * @param callback - Function to call on close
     */
    public setOnClose(callback: () => void): void {
        this.onCloseCallback = callback;
    }

    /**
     * Get the content container element
     */
    public getContentContainer(): HTMLElement {
        return this.contentElement;
    }
}
