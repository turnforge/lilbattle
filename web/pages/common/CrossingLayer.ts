/**
 * Crossing Layer for Roads and Bridges
 *
 * This layer renders terrain improvements (roads on land, bridges on water)
 * as visual overlays on the hex grid. Roads and bridges are stored separately
 * from tiles to allow independent terrain modification while preserving crossings.
 *
 * Connection-based rendering:
 * - If a crossing has no neighbors with crossings, draws a horizontal line (left to right edge)
 * - If there are neighboring crossings, draws lines from center to center of connected tiles
 *
 * Depth: 5 (between tiles at 0 and units at 10)
 */

import * as Phaser from 'phaser';
import { BaseLayer, LayerConfig, ClickContext, LayerHitResult } from './LayerSystem';
import { hexToPixel, axialNeighbors } from './hexUtils';
import { CrossingType } from './World';

// =============================================================================
// Crossing Layer
// =============================================================================

/**
 * Layer for rendering roads and bridges with connection-based graphics
 */
export class CrossingLayer extends BaseLayer {
    private crossingGraphics = new Map<string, Phaser.GameObjects.Graphics>();
    private crossingData = new Map<string, CrossingType>();
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
     * Get neighboring hex coordinates
     */
    private getNeighbors(q: number, r: number): { q: number, r: number }[] {
        return axialNeighbors(q, r).map(([nq, nr]) => ({ q: nq, r: nr }));
    }

    /**
     * Check if two crossing types are compatible for connection
     * Roads only connect to roads, bridges only connect to bridges
     */
    private areTypesCompatible(type1: CrossingType, type2: CrossingType): boolean {
        const isRoad1 = type1 === CrossingType.CROSSING_TYPE_ROAD;
        const isRoad2 = type2 === CrossingType.CROSSING_TYPE_ROAD;
        // Both must be roads, or both must be bridges
        return isRoad1 === isRoad2;
    }

    /**
     * Get neighbors that have compatible crossings
     */
    private getConnectedNeighbors(q: number, r: number): { q: number, r: number }[] {
        const currentType = this.crossingData.get(this.getKey(q, r));
        if (!currentType) return [];

        return this.getNeighbors(q, r).filter(n => {
            const neighborType = this.crossingData.get(this.getKey(n.q, n.r));
            return neighborType && this.areTypesCompatible(currentType, neighborType);
        });
    }

    /**
     * Add or update a crossing at a hex coordinate
     */
    public setCrossing(q: number, r: number, crossingType: CrossingType): void {
        const key = this.getKey(q, r);

        if (crossingType === CrossingType.CROSSING_TYPE_UNSPECIFIED) {
            this.removeCrossing(q, r);
            return;
        }

        // Store the crossing data
        this.crossingData.set(key, crossingType);

        // Redraw this tile and all compatible neighbors that might be affected
        this.redrawTile(q, r);
        for (const neighbor of this.getNeighbors(q, r)) {
            const neighborType = this.crossingData.get(this.getKey(neighbor.q, neighbor.r));
            if (neighborType && this.areTypesCompatible(crossingType, neighborType)) {
                this.redrawTile(neighbor.q, neighbor.r);
            }
        }
    }

    /**
     * Remove crossing at a hex coordinate
     */
    public removeCrossing(q: number, r: number): void {
        const key = this.getKey(q, r);

        // Get neighbors before removing (to redraw them after)
        const connectedNeighbors = this.getConnectedNeighbors(q, r);

        // Remove graphics
        const graphics = this.crossingGraphics.get(key);
        if (graphics) {
            graphics.destroy();
            this.crossingGraphics.delete(key);
        }

        // Remove data
        this.crossingData.delete(key);

        // Redraw neighbors that were connected
        for (const neighbor of connectedNeighbors) {
            this.redrawTile(neighbor.q, neighbor.r);
        }
    }

    /**
     * Redraw the crossing graphic for a single tile based on its connections
     */
    private redrawTile(q: number, r: number): void {
        const key = this.getKey(q, r);
        const crossingType = this.crossingData.get(key);

        if (!crossingType) return;

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

        // Find connected neighbors
        const connectedNeighbors = this.getConnectedNeighbors(q, r);

        if (connectedNeighbors.length === 0) {
            // No connections - draw default horizontal crossing
            this.drawDefaultCrossing(graphics, crossingType);
        } else {
            // Draw connections to each neighbor
            // We only draw from current tile to neighbors, not back
            // This is handled by each tile drawing its own outgoing connections
            for (const neighbor of connectedNeighbors) {
                this.drawConnectionToNeighbor(graphics, q, r, neighbor.q, neighbor.r, crossingType);
            }
        }

        this.crossingGraphics.set(key, graphics);
    }

    /**
     * Draw the default crossing (horizontal line) when no neighbors are connected
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
     * Draw a connection line from current tile center toward a neighbor
     * We draw from center to the edge (halfway to neighbor center)
     */
    private drawConnectionToNeighbor(
        graphics: Phaser.GameObjects.Graphics,
        fromQ: number, fromR: number,
        toQ: number, toR: number,
        crossingType: CrossingType
    ): void {
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
     * Load crossings from a map of coordinate keys to CrossingType values
     */
    public loadCrossings(crossings: { [key: string]: CrossingType }): void {
        // Clear existing crossings
        this.clearAllCrossings();

        // First, store all crossing data
        for (const [key, crossingType] of Object.entries(crossings)) {
            const [q, r] = key.split(',').map(Number);
            if (!isNaN(q) && !isNaN(r) && crossingType !== CrossingType.CROSSING_TYPE_UNSPECIFIED) {
                this.crossingData.set(key, crossingType);
            }
        }

        // Then, draw all tiles (now that we know all neighbors)
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
