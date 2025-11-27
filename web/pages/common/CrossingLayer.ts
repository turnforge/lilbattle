/**
 * Crossing Layer for Roads and Bridges
 *
 * This layer renders terrain improvements (roads on land, bridges on water)
 * as visual overlays on the hex grid. Roads and bridges are stored separately
 * from tiles to allow independent terrain modification while preserving crossings.
 *
 * Depth: 5 (between tiles at 0 and units at 10)
 */

import * as Phaser from 'phaser';
import { BaseLayer, LayerConfig, ClickContext, LayerHitResult } from './LayerSystem';
import { hexToPixel } from './hexUtils';
import { CrossingType } from './World';

// =============================================================================
// Crossing Layer
// =============================================================================

/**
 * Layer for rendering roads and bridges
 */
export class CrossingLayer extends BaseLayer {
    private crossingGraphics = new Map<string, Phaser.GameObjects.Graphics>();
    private tileWidth: number;

    constructor(scene: Phaser.Scene, tileWidth: number) {
        super(scene, {
            name: 'crossings',
            coordinateSpace: 'hex',
            interactive: false, // Crossings are visual only, don't consume clicks
            depth: 5, // Between tiles (0) and units (10)
        });
        this.tileWidth = tileWidth;
    }

    public hitTest(context: ClickContext): LayerHitResult | null {
        // Crossings are visual only, never intercept clicks
        return LayerHitResult.TRANSPARENT;
    }

    /**
     * Add or update a crossing at a hex coordinate
     */
    public setCrossing(q: number, r: number, crossingType: CrossingType): void {
        const key = `${q},${r}`;

        // Remove existing crossing graphic if present
        this.removeCrossing(q, r);

        if (crossingType === CrossingType.CROSSING_TYPE_UNSPECIFIED) {
            return; // No crossing to render
        }

        // Create graphics for the crossing
        const graphics = this.scene.add.graphics();
        this.container.add(graphics);

        // Get world position
        const position = hexToPixel(q, r);
        graphics.setPosition(position.x, position.y);

        if (crossingType === CrossingType.CROSSING_TYPE_ROAD) {
            this.drawRoad(graphics);
        } else if (crossingType === CrossingType.CROSSING_TYPE_BRIDGE) {
            this.drawBridge(graphics);
        }

        this.crossingGraphics.set(key, graphics);
    }

    /**
     * Draw a road graphic
     * Roads are rendered as a brownish path pattern
     */
    private drawRoad(graphics: Phaser.GameObjects.Graphics): void {
        const halfWidth = this.tileWidth / 2;

        // Road base (brown/tan color)
        graphics.fillStyle(0x8B7355, 0.7); // Tan/dirt color

        // Draw a cross pattern for the road (connecting hex edges)
        const roadWidth = this.tileWidth * 0.25;
        const halfRoad = roadWidth / 2;

        // Horizontal road segment (left to right)
        graphics.fillRect(-halfWidth * 0.8, -halfRoad, halfWidth * 1.6, roadWidth);

        // Draw road edge lines
        graphics.lineStyle(2, 0x5D4E37, 0.8); // Darker brown edge
        graphics.strokeRect(-halfWidth * 0.8, -halfRoad, halfWidth * 1.6, roadWidth);
    }

    /**
     * Draw a bridge graphic
     * Bridges are rendered with wooden plank style
     */
    private drawBridge(graphics: Phaser.GameObjects.Graphics): void {
        const halfWidth = this.tileWidth / 2;

        // Bridge base (wooden color)
        graphics.fillStyle(0x8B4513, 0.8); // Saddle brown (wood)

        // Draw bridge as a horizontal path
        const bridgeWidth = this.tileWidth * 0.3;
        const halfBridge = bridgeWidth / 2;

        // Main bridge planks
        graphics.fillRect(-halfWidth * 0.75, -halfBridge, halfWidth * 1.5, bridgeWidth);

        // Bridge railings (darker wood)
        graphics.lineStyle(3, 0x654321, 1.0); // Dark brown
        graphics.strokeRect(-halfWidth * 0.75, -halfBridge, halfWidth * 1.5, bridgeWidth);

        // Draw plank lines across the bridge for texture
        graphics.lineStyle(1, 0x5D4E37, 0.6);
        const plankSpacing = this.tileWidth * 0.12;
        for (let x = -halfWidth * 0.6; x < halfWidth * 0.6; x += plankSpacing) {
            graphics.moveTo(x, -halfBridge);
            graphics.lineTo(x, halfBridge);
        }
        graphics.strokePath();

        // Support posts on ends (darker)
        graphics.fillStyle(0x4A3728, 1.0);
        const postWidth = this.tileWidth * 0.06;
        graphics.fillRect(-halfWidth * 0.75 - postWidth / 2, -halfBridge * 1.2, postWidth, bridgeWidth * 1.4);
        graphics.fillRect(halfWidth * 0.75 - postWidth / 2, -halfBridge * 1.2, postWidth, bridgeWidth * 1.4);
    }

    /**
     * Remove crossing at a hex coordinate
     */
    public removeCrossing(q: number, r: number): void {
        const key = `${q},${r}`;
        const graphics = this.crossingGraphics.get(key);

        if (graphics) {
            graphics.destroy();
            this.crossingGraphics.delete(key);
        }
    }

    /**
     * Clear all crossings
     */
    public clearAllCrossings(): void {
        for (const graphics of this.crossingGraphics.values()) {
            graphics.destroy();
        }
        this.crossingGraphics.clear();
    }

    /**
     * Load crossings from a map of coordinate keys to CrossingType values
     */
    public loadCrossings(crossings: { [key: string]: CrossingType }): void {
        // Clear existing crossings
        this.clearAllCrossings();

        // Add all crossings from the map
        for (const [key, crossingType] of Object.entries(crossings)) {
            const [q, r] = key.split(',').map(Number);
            if (!isNaN(q) && !isNaN(r)) {
                this.setCrossing(q, r, crossingType);
            }
        }
    }

    /**
     * Check if there's a crossing at the given hex coordinate
     */
    public hasCrossing(q: number, r: number): boolean {
        const key = `${q},${r}`;
        return this.crossingGraphics.has(key);
    }

    public destroy(): void {
        this.clearAllCrossings();
        super.destroy();
    }
}
