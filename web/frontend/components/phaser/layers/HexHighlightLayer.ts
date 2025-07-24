/**
 * Hex-based highlight layers for game interactions
 * 
 * These layers work in hex coordinate space and provide visual feedback
 * for movement, attack, and selection in the game.
 */

import * as Phaser from 'phaser';
import { BaseLayer, LayerConfig, ClickContext, LayerHitResult } from '../LayerSystem';
import { hexToPixel } from '../hexUtils';

// =============================================================================
// Hex Highlight Base Class
// =============================================================================

/**
 * Base class for all hex-based highlight layers
 */
export abstract class HexHighlightLayer extends BaseLayer {
    protected highlights = new Map<string, Phaser.GameObjects.Graphics>();
    protected tileWidth: number;
    
    constructor(scene: Phaser.Scene, config: LayerConfig & { tileWidth: number }) {
        super(scene, { ...config, coordinateSpace: 'hex' });
        this.tileWidth = config.tileWidth;
    }
    
    /**
     * Add highlight at hex coordinate
     */
    protected addHighlight(q: number, r: number, color: number, alpha: number = 0.3, strokeColor?: number, strokeWidth?: number): void {
        const key = `${q},${r}`;
        
        // Remove existing highlight if present
        this.removeHighlight(q, r);
        
        // Create new highlight
        const highlight = this.scene.add.graphics();
        
        // Set fill style
        highlight.fillStyle(color, alpha);
        if (strokeColor !== undefined && strokeWidth !== undefined) {
            highlight.lineStyle(strokeWidth, strokeColor, 1.0);
        }
        
        // Draw hexagon
        const position = hexToPixel(q, r);
        this.drawHexagon(highlight, position.x, position.y, this.tileWidth * 0.85);
        
        // Add to container and store reference
        this.container.add(highlight);
        this.highlights.set(key, highlight);
    }
    
    /**
     * Remove highlight at hex coordinate
     */
    protected removeHighlight(q: number, r: number): void {
        const key = `${q},${r}`;
        const highlight = this.highlights.get(key);
        
        if (highlight) {
            highlight.destroy();
            this.highlights.delete(key);
        }
    }
    
    /**
     * Check if there's a highlight at the given hex coordinate
     */
    protected hasHighlight(q: number, r: number): boolean {
        const key = `${q},${r}`;
        return this.highlights.has(key);
    }
    
    /**
     * Clear all highlights
     */
    protected clearHighlights(): void {
        for (const highlight of this.highlights.values()) {
            highlight.destroy();
        }
        this.highlights.clear();
    }
    
    /**
     * Draw hexagon shape on graphics object
     */
    private drawHexagon(graphics: Phaser.GameObjects.Graphics, x: number, y: number, size: number): void {
        const points: number[] = [];
        
        for (let i = 0; i < 6; i++) {
            const angle = (Math.PI / 3) * i;
            const px = x + size * Math.cos(angle);
            const py = y + size * Math.sin(angle);
            points.push(px, py);
        }
        
        graphics.fillPoints(points, true);
        if (graphics.lineStyle) {
            graphics.strokePoints(points, true);
        }
    }
    
    public destroy(): void {
        this.clearHighlights();
        super.destroy();
    }
}

// =============================================================================
// Selection Highlight Layer
// =============================================================================

/**
 * Shows yellow highlight for currently selected unit
 */
export class SelectionHighlightLayer extends HexHighlightLayer {
    private selectedCoord: { q: number; r: number } | null = null;
    
    constructor(scene: Phaser.Scene, tileWidth: number) {
        super(scene, {
            name: 'selection-highlight',
            coordinateSpace: 'hex',
            interactive: false, // Selection highlights don't consume clicks
            depth: 10, // High priority visual
            tileWidth
        });
    }
    
    public hitTest(context: ClickContext): LayerHitResult | null {
        // Selection highlights are visual only, never intercept clicks
        return LayerHitResult.TRANSPARENT;
    }
    
    /**
     * Show selection highlight at hex coordinate
     */
    public selectHex(q: number, r: number): void {
        console.log(`[SelectionHighlightLayer] Selecting hex (${q}, ${r})`);
        
        // Clear previous selection
        this.clearSelection();
        
        // Add new selection highlight (yellow with border)
        this.addHighlight(q, r, 0xFFFF00, 0.2, 0xFFFF00, 4);
        this.selectedCoord = { q, r };
    }
    
    /**
     * Clear current selection
     */
    public clearSelection(): void {
        if (this.selectedCoord) {
            console.log(`[SelectionHighlightLayer] Clearing selection`);
            this.removeHighlight(this.selectedCoord.q, this.selectedCoord.r);
            this.selectedCoord = null;
        }
    }
    
    /**
     * Get currently selected coordinate
     */
    public getSelection(): { q: number; r: number } | null {
        return this.selectedCoord;
    }
}

// =============================================================================
// Movement Highlight Layer
// =============================================================================

/**
 * Shows green highlights for valid movement positions
 */
export class MovementHighlightLayer extends HexHighlightLayer {
    private onMoveCallback?: (q: number, r: number) => void;
    
    constructor(scene: Phaser.Scene, tileWidth: number, onMoveCallback?: (q: number, r: number) => void) {
        super(scene, {
            name: 'movement-highlight',
            coordinateSpace: 'hex',
            interactive: true, // Movement highlights consume clicks
            depth: 5, // Below selection, above base map
            tileWidth
        });
        
        this.onMoveCallback = onMoveCallback;
    }
    
    public hitTest(context: ClickContext): LayerHitResult | null {
        if (!this.visible) return null;
        
        // Only consume clicks if there's a movement highlight at this position
        if (this.hasHighlight(context.hexQ, context.hexR)) {
            return LayerHitResult.CONSUME;
        }
        
        return LayerHitResult.TRANSPARENT;
    }
    
    public handleClick(context: ClickContext): boolean {
        console.log(`[MovementHighlightLayer] Movement click at (${context.hexQ}, ${context.hexR})`);
        
        if (this.onMoveCallback) {
            this.onMoveCallback(context.hexQ, context.hexR);
        }
        
        return true; // Event handled
    }
    
    /**
     * Show movement options
     */
    public showMovementOptions(coords: Array<{ q: number; r: number; cost?: number }>): void {
        console.log(`[MovementHighlightLayer] Showing ${coords.length} movement options`);
        
        // Clear existing highlights
        this.clearHighlights();
        
        // Add highlights for each valid movement position
        coords.forEach(coord => {
            // Green highlight with subtle border
            this.addHighlight(coord.q, coord.r, 0x00FF00, 0.15, 0x00FF00, 2);
        });
    }
    
    /**
     * Clear all movement highlights
     */
    public clearMovementOptions(): void {
        console.log(`[MovementHighlightLayer] Clearing movement options`);
        this.clearHighlights();
    }
}

// =============================================================================
// Attack Highlight Layer
// =============================================================================

/**
 * Shows red highlights for valid attack targets
 */
export class AttackHighlightLayer extends HexHighlightLayer {
    private onAttackCallback?: (q: number, r: number) => void;
    
    constructor(scene: Phaser.Scene, tileWidth: number, onAttackCallback?: (q: number, r: number) => void) {
        super(scene, {
            name: 'attack-highlight',
            coordinateSpace: 'hex',
            interactive: true, // Attack highlights consume clicks
            depth: 6, // Same level as movement, both are action highlights
            tileWidth
        });
        
        this.onAttackCallback = onAttackCallback;
    }
    
    public hitTest(context: ClickContext): LayerHitResult | null {
        if (!this.visible) return null;
        
        // Only consume clicks if there's an attack highlight at this position
        if (this.hasHighlight(context.hexQ, context.hexR)) {
            return LayerHitResult.CONSUME;
        }
        
        return LayerHitResult.TRANSPARENT;
    }
    
    public handleClick(context: ClickContext): boolean {
        console.log(`[AttackHighlightLayer] Attack click at (${context.hexQ}, ${context.hexR})`);
        
        if (this.onAttackCallback) {
            this.onAttackCallback(context.hexQ, context.hexR);
        }
        
        return true; // Event handled
    }
    
    /**
     * Show attack options
     */
    public showAttackOptions(coords: Array<{ q: number; r: number }>): void {
        console.log(`[AttackHighlightLayer] Showing ${coords.length} attack options`);
        
        // Clear existing highlights
        this.clearHighlights();
        
        // Add highlights for each valid attack target
        coords.forEach(coord => {
            // Red highlight with border
            this.addHighlight(coord.q, coord.r, 0xFF0000, 0.2, 0xFF0000, 2);
        });
    }
    
    /**
     * Clear all attack highlights
     */
    public clearAttackOptions(): void {
        console.log(`[AttackHighlightLayer] Clearing attack options`);
        this.clearHighlights();
    }
}