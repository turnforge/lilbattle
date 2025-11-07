import { ShapeTool, HexCoord } from './ShapeTool';
import { World } from '../../common/World';
import { hexDistance } from '../../common/hexUtils';

/**
 * Circle drawing tool.
 *
 * Workflow:
 * 1. First click: Set center point
 * 2. Mouse move: Show circle preview with radius = distance from center to cursor
 * 3. Second click: Complete circle
 * 4. Escape: Cancel
 */
export class CircleTool implements ShapeTool {
  public readonly name = 'Circle';

  private center: HexCoord | null = null;
  private radiusPoint: HexCoord | null = null;
  private filled: boolean = true;
  private world: World;

  constructor(world: World, filled: boolean = true) {
    this.world = world;
    this.filled = filled;
  }

  addPoint(q: number, r: number): boolean {
    if (this.center === null) {
      // First click: Store center
      this.center = { q, r };
      return true; // More points needed
    } else {
      // Second click: Store radius point and complete
      this.radiusPoint = { q, r };
      return false; // Shape complete
    }
  }

  getPreviewTiles(currentQ: number, currentR: number): HexCoord[] {
    if (this.center === null) {
      return []; // No preview until center is set
    }

    // Calculate radius as hex distance from center to current position
    const radius = hexDistance(this.center.q, this.center.r, currentQ, currentR);

    if (radius === 0) {
      // Show just the center point if radius is 0
      return [this.center];
    }

    // Show preview circle (always outline for preview)
    const tiles = this.world.circleFrom(
      this.center.q,
      this.center.r,
      radius,
      false // Always outline for preview
    );

    return tiles.map(([q, r]) => ({ q, r }));
  }

  getResultTiles(): HexCoord[] {
    if (this.center === null || this.radiusPoint === null) {
      return []; // No result if incomplete
    }

    // Calculate final radius
    const radius = hexDistance(
      this.center.q,
      this.center.r,
      this.radiusPoint.q,
      this.radiusPoint.r
    );

    if (radius === 0) {
      // Single tile if radius is 0
      return [this.center];
    }

    // Generate final circle with current fill setting
    const tiles = this.world.circleFrom(
      this.center.q,
      this.center.r,
      radius,
      this.filled
    );

    return tiles.map(([q, r]) => ({ q, r }));
  }

  getAnchorPoints(): HexCoord[] {
    const points: HexCoord[] = [];
    if (this.center !== null) {
      points.push(this.center);
    }
    if (this.radiusPoint !== null) {
      points.push(this.radiusPoint);
    }
    return points;
  }

  reset(): void {
    this.center = null;
    this.radiusPoint = null;
  }

  canComplete(): boolean {
    return this.center !== null && this.radiusPoint !== null;
  }

  requiresKeyboardConfirm(): boolean {
    return false; // Circle auto-completes after 2 clicks
  }

  getStatusText(): string {
    if (this.center === null) {
      return 'Click center of circle';
    } else if (this.radiusPoint === null) {
      return 'Click to set radius (or press Escape to cancel)';
    } else {
      return 'Circle complete';
    }
  }

  isFilled(): boolean {
    return this.filled;
  }

  setFilled(filled: boolean): void {
    this.filled = filled;
  }
}
