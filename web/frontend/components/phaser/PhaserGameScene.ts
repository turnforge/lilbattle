import * as Phaser from 'phaser';
import { PhaserWorldScene } from './PhaserWorldScene';
import { World } from '../World';

/**
 * Callback interfaces for game-specific interactions
 */
export interface GameSceneCallbacks {
    onTileClicked?: (q: number, r: number) => void;
    onUnitClicked?: (q: number, r: number) => void;
}

/**
 * PhaserGameScene extends PhaserWorldScene with game-specific interactive features.
 * 
 * This scene adds:
 * - Click handling for terrain tiles (shows terrain info in TerrainStatsPanel)
 * - Click handling for units (shows movement/attack options)
 * - Game-specific UI feedback and highlighting
 * - Integration with GameViewerPage callbacks
 * 
 * Inherits from PhaserWorldScene:
 * - World as single source of truth for game data
 * - Tile and unit rendering using World data
 * - Camera controls and theme management
 * - Asset loading and coordinate conversion
 */
export class PhaserGameScene extends PhaserWorldScene {
    private callbacks: GameSceneCallbacks = {};
    
    constructor(config?: string | Phaser.Types.Scenes.SettingsConfig) {
        super(config || { key: 'PhaserGameScene' });
    }

    /**
     * Set callback functions for game interactions
     */
    public setCallbacks(callbacks: GameSceneCallbacks): void {
        this.callbacks = callbacks;
        console.log('[PhaserGameScene] Callbacks set:', Object.keys(callbacks));
    }

    /**
     * Override the base tile click handler to add game-specific logic
     */
    protected onTileClick(q: number, r: number): void {
        console.log(`[PhaserGameScene] Tile clicked: Q=${q}, R=${r}`);
        
        if (!this.world) {
            console.warn('[PhaserGameScene] No World available for tile click');
            return;
        }

        // Check if there's a unit at this position first
        const unit = this.world.getUnitAt(q, r);
        if (unit) {
            console.log(`[PhaserGameScene] Unit found at tile: Q=${q}, R=${r}, Type=${unit.unitType}, Player=${unit.player}`);
            
            // Call unit clicked callback
            if (this.callbacks.onUnitClicked) {
                this.callbacks.onUnitClicked(q, r);
            }
        } else {
            // No unit, this is a terrain click
            const tile = this.world.getTileAt(q, r);
            if (tile) {
                console.log(`[PhaserGameScene] Terrain clicked: Q=${q}, R=${r}, Type=${tile.tileType}, Player=${tile.player}`);
            } else {
                console.log(`[PhaserGameScene] Empty tile clicked: Q=${q}, R=${r}`);
            }
            
            // Call tile clicked callback
            if (this.callbacks.onTileClicked) {
                this.callbacks.onTileClicked(q, r);
            }
        }

        // Also emit the base event for any other listeners
        this.events.emit('tileClicked', { q, r });
    }

    /**
     * Get tile data at specific coordinates (for callback functions)
     */
    public getTileAt(q: number, r: number): any {
        if (!this.world) {
            return null;
        }
        return this.world.getTileAt(q, r);
    }

    /**
     * Get unit data at specific coordinates (for callback functions) 
     */
    public getUnitAt(q: number, r: number): any {
        if (!this.world) {
            return null;
        }
        return this.world.getUnitAt(q, r);
    }

    /**
     * Check if there's a unit at the specified coordinates
     */
    public hasUnitAt(q: number, r: number): boolean {
        if (!this.world) {
            return false;
        }
        return this.world.getUnitAt(q, r) !== null;
    }

    /**
     * Check if there's a tile at the specified coordinates
     */
    public hasTileAt(q: number, r: number): boolean {
        if (!this.world) {
            return false;
        }
        return this.world.getTileAt(q, r) !== null;
    }
}
