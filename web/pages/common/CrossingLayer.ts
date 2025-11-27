/**
 * Crossing Layer for Roads and Bridges
 *
 * This layer renders terrain improvements (roads on land, bridges on water)
 * as visual overlays on the hex grid. Roads and bridges are stored separately
 * from tiles to allow independent terrain modification while preserving crossings.
 *
 * Explicit connectivity rendering:
 * - Each crossing stores which of its 6 hex neighbors it connects to via connectsTo array
 * - If a crossing has no connections (all false), draws a horizontal line (left to right edge)
 * - Otherwise, draws lines from center toward each connected direction
 *
 * Depth: 5 (between tiles at 0 and units at 10)
 */

import * as Phaser from 'phaser';
import { BaseLayer, LayerConfig, ClickContext, LayerHitResult } from './LayerSystem';
import { hexToPixel, getNeighborCoord } from './hexUtils';
import { CrossingType, Crossing } from './World';

// =============================================================================
// Crossing Layer
// =============================================================================

/**
 * Layer for rendering roads and bridges with explicit connection-based graphics
 */
export class CrossingLayer extends BaseLayer {
    private crossingGraphics = new Map<string, Phaser.GameObjects.Graphics>();
    private crossingData = new Map<string, Crossing>();
    private tileWidth: number;
    private tileHeight: number;

    constructor(scene: Phaser.Scene, tileWidth: number) {
        super(scene, {
            name: 'crossings',
            coordinateSpace: 'hex',
            interactive: false, // Crossings are visual only, don't consume clicks
            depth: 5, // Between tiles (0) and units (10)
        });
        this.tileWidth = tileWidth;
        this.tileHeight = tileWidth; // Assuming square-ish hexes
    }

    public hitTest(context: ClickContext): LayerHitResult | null {
        // Crossings are visual only, never intercept clicks
        return LayerHitResult.TRANSPARENT;
    }

    /**
     * Get hex key from coordinates
     */
    private getKey(q: number, r: number): string {
        return `${q},${r}`;
    }

    /**
     * Get the direction indices where this crossing has connections
     * Reads directly from connectsTo array
     */
    private getConnectionDirections(q: number, r: number): number[] {
        const crossing = this.crossingData.get(this.getKey(q, r));
        if (!crossing) return [];

        const directions: number[] = [];
        for (let i = 0; i < 6; i++) {
            if (crossing.connectsTo[i]) {
                directions.push(i);
            }
        }
        return directions;
    }

    /**
     * Add or update a crossing at a hex coordinate
     */
    public setCrossing(q: number, r: number, crossing: Crossing): void {
        const key = this.getKey(q, r);

        if (crossing.type === CrossingType.CROSSING_TYPE_UNSPECIFIED) {
            this.removeCrossing(q, r);
            return;
        }

        // Store the crossing data
        this.crossingData.set(key, crossing);

        // Redraw this tile
        this.redrawTile(q, r);
    }

    /**
     * Remove crossing at a hex coordinate
     */
    public removeCrossing(q: number, r: number): void {
        const key = this.getKey(q, r);

        // Remove graphics
        const graphics = this.crossingGraphics.get(key);
        if (graphics) {
            graphics.destroy();
            this.crossingGraphics.delete(key);
        }

        // Remove data
        this.crossingData.delete(key);
    }

    /**
     * Redraw the crossing graphic for a single tile based on its explicit connections
     */
    private redrawTile(q: number, r: number): void {
        const key = this.getKey(q, r);
        const crossing = this.crossingData.get(key);

        if (!crossing) return;

        // Remove existing graphic
        const existing = this.crossingGraphics.get(key);
        if (existing) {
            existing.destroy();
        }

        // Create new graphics
        const graphics = this.scene.add.graphics();
        this.container.add(graphics);

        // Get world position for this tile's center
        const position = hexToPixel(q, r);
        graphics.setPosition(position.x, position.y);

        // Get explicit connection directions
        const connectionDirections = this.getConnectionDirections(q, r);

        if (connectionDirections.length === 0) {
            // No connections - draw default horizontal crossing
            this.drawDefaultCrossing(graphics, crossing.type);
        } else {
            // Draw connections in each specified direction
            for (const direction of connectionDirections) {
                this.drawConnectionInDirection(graphics, q, r, direction, crossing.type);
            }
        }

        this.crossingGraphics.set(key, graphics);
    }

    /**
     * Draw the default crossing (horizontal line) when no connections are specified
     */
    private drawDefaultCrossing(graphics: Phaser.GameObjects.Graphics, crossingType: CrossingType): void {
        const halfWidth = this.tileWidth / 2;
        const lineWidth = crossingType === CrossingType.CROSSING_TYPE_BRIDGE ? 12 : 10;

        if (crossingType === CrossingType.CROSSING_TYPE_ROAD) {
            // Road: tan/brown color
            graphics.lineStyle(lineWidth, 0x8B7355, 0.9);
            graphics.lineBetween(-halfWidth * 0.7, 0, halfWidth * 0.7, 0);
            // Edge lines
            graphics.lineStyle(2, 0x5D4E37, 0.8);
            graphics.lineBetween(-halfWidth * 0.7, -lineWidth / 2, halfWidth * 0.7, -lineWidth / 2);
            graphics.lineBetween(-halfWidth * 0.7, lineWidth / 2, halfWidth * 0.7, lineWidth / 2);
        } else {
            // Bridge: wooden brown
            graphics.lineStyle(lineWidth, 0x8B4513, 0.9);
            graphics.lineBetween(-halfWidth * 0.7, 0, halfWidth * 0.7, 0);
            // Railings
            graphics.lineStyle(3, 0x654321, 1.0);
            graphics.lineBetween(-halfWidth * 0.7, -lineWidth / 2, halfWidth * 0.7, -lineWidth / 2);
            graphics.lineBetween(-halfWidth * 0.7, lineWidth / 2, halfWidth * 0.7, lineWidth / 2);
        }
    }

    /**
     * Draw a connection line from current tile center toward a neighbor in the given direction
     * We draw from center to the edge (halfway to neighbor center)
     */
    private drawConnectionInDirection(
        graphics: Phaser.GameObjects.Graphics,
        fromQ: number, fromR: number,
        direction: number,
        crossingType: CrossingType
    ): void {
        // Get neighbor coordinate in this direction
        const [toQ, toR] = getNeighborCoord(fromQ, fromR, direction);

        // Calculate relative position of neighbor center from our center
        const fromPos = hexToPixel(fromQ, fromR);
        const toPos = hexToPixel(toQ, toR);

        // Direction vector from current tile to neighbor (relative to our position at 0,0)
        const dx = toPos.x - fromPos.x;
        const dy = toPos.y - fromPos.y;

        // Draw from center (0,0) to halfway point (edge of our hex)
        const endX = dx / 2;
        const endY = dy / 2;

        const lineWidth = crossingType === CrossingType.CROSSING_TYPE_BRIDGE ? 12 : 10;

        if (crossingType === CrossingType.CROSSING_TYPE_ROAD) {
            // Road: tan/brown with edge lines
            graphics.lineStyle(lineWidth, 0x8B7355, 0.9);
            graphics.lineBetween(0, 0, endX, endY);

            // Draw edge lines parallel to the path
            const length = Math.sqrt(endX * endX + endY * endY);
            if (length > 0) {
                const perpX = (-endY / length) * (lineWidth / 2);
                const perpY = (endX / length) * (lineWidth / 2);
                graphics.lineStyle(2, 0x5D4E37, 0.8);
                graphics.lineBetween(perpX, perpY, endX + perpX, endY + perpY);
                graphics.lineBetween(-perpX, -perpY, endX - perpX, endY - perpY);
            }
        } else {
            // Bridge: wooden with railings
            graphics.lineStyle(lineWidth, 0x8B4513, 0.9);
            graphics.lineBetween(0, 0, endX, endY);

            // Draw railings parallel to the path
            const length = Math.sqrt(endX * endX + endY * endY);
            if (length > 0) {
                const perpX = (-endY / length) * (lineWidth / 2);
                const perpY = (endX / length) * (lineWidth / 2);
                graphics.lineStyle(3, 0x654321, 1.0);
                graphics.lineBetween(perpX, perpY, endX + perpX, endY + perpY);
                graphics.lineBetween(-perpX, -perpY, endX - perpX, endY - perpY);
            }
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
        this.crossingData.clear();
    }

    /**
     * Load crossings from a map of coordinate keys to Crossing objects
     */
    public loadCrossings(crossings: { [key: string]: Crossing }): void {
        // Clear existing crossings
        this.clearAllCrossings();

        // First, store all crossing data
        for (const [key, crossing] of Object.entries(crossings)) {
            const [q, r] = key.split(',').map(Number);
            if (!isNaN(q) && !isNaN(r) && crossing.type !== CrossingType.CROSSING_TYPE_UNSPECIFIED) {
                this.crossingData.set(key, crossing);
            }
        }

        // Then, draw all tiles
        for (const [key] of this.crossingData) {
            const [q, r] = key.split(',').map(Number);
            this.redrawTile(q, r);
        }
    }

    /**
     * Check if there's a crossing at the given hex coordinate
     */
    public hasCrossing(q: number, r: number): boolean {
        return this.crossingData.has(this.getKey(q, r));
    }

    public destroy(): void {
        this.clearAllCrossings();
        super.destroy();
    }
}
