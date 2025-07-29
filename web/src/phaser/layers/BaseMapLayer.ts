/**
 * Base Map Layer for Default Interactions
 * 
 * This layer works in hex coordinate space and handles default tile/unit clicks
 * when no other interactive layer consumes the event. Acts as the fallback
 * for basic map interactions.
 */

import * as Phaser from 'phaser';
import { BaseLayer, LayerConfig, ClickContext, LayerHitResult } from '../LayerSystem';

export interface MapLayerCallbacks {
    onTileClicked?: (q: number, r: number) => boolean;
    onUnitClicked?: (q: number, r: number) => boolean;
    onEmptySpaceClicked?: (q: number, r: number) => boolean;
}

/**
 * Base map layer that handles default tile and unit interactions
 */
export class BaseMapLayer extends BaseLayer {
    private callbacks: MapLayerCallbacks = {};
    
    constructor(scene: Phaser.Scene, callbacks: MapLayerCallbacks = {}) {
        super(scene, {
            name: 'base-map',
            coordinateSpace: 'hex',
            interactive: true,
            depth: 0, // Lowest priority - only handles events no other layer wants
        });
        
        this.callbacks = callbacks;
    }
    
    public hitTest(context: ClickContext): LayerHitResult | null {
        // Base map layer always consumes events that reach it
        // This ensures there's always a fallback handler
        return LayerHitResult.CONSUME;
    }
    
    public handleClick(context: ClickContext): boolean {
        console.log(`[BaseMapLayer] Handling default click at hex (${context.hexQ}, ${context.hexR})`);
        
        // Determine what was clicked and call appropriate callback
        if (context.unit) {
            console.log(`[BaseMapLayer] Unit click detected`);
            if (this.callbacks.onUnitClicked) {
                return this.callbacks.onUnitClicked(context.hexQ, context.hexR);
            }
        } else if (context.tile) {
            console.log(`[BaseMapLayer] Tile click detected`);
            if (this.callbacks.onTileClicked) {
                return this.callbacks.onTileClicked(context.hexQ, context.hexR);
            }
        } else {
            console.log(`[BaseMapLayer] Empty space click detected`);
            if (this.callbacks.onEmptySpaceClicked) {
                return this.callbacks.onEmptySpaceClicked(context.hexQ, context.hexR);
            }
        }
        
        // Default behavior if no specific callback handled the event
        console.log(`[BaseMapLayer] No specific handler for click, using default behavior`);
        return true; // Event handled
    }
    
    /**
     * Update the callbacks used by this layer
     */
    public setCallbacks(callbacks: MapLayerCallbacks): void {
        console.log(`[BaseMapLayer] Updating callbacks`);
        this.callbacks = { ...callbacks };
    }
    
    /**
     * Set tile click callback
     */
    public setTileClickCallback(callback?: (q: number, r: number) => boolean): void {
        this.callbacks.onTileClicked = callback;
    }
    
    /**
     * Set unit click callback
     */
    public setUnitClickCallback(callback?: (q: number, r: number) => boolean): void {
        this.callbacks.onUnitClicked = callback;
    }
    
    /**
     * Set empty space click callback
     */
    public setEmptySpaceClickCallback(callback?: (q: number, r: number) => boolean): void {
        this.callbacks.onEmptySpaceClicked = callback;
    }
}