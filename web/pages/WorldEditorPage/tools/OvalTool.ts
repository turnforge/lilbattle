import { ShapeTool, HexCoord } from './ShapeTool';
import { World } from '../../common/World';
import { hexToRowCol } from '../../common/hexUtils';

/**
 * Oval/Ellipse drawing tool (axis-aligned).
 *
 * Workflow:
 * 1. First click: Set center point
 * 2. Second click: Set radiusX (horizontal radius in row/col space)
 * 3. Mouse move: Show oval preview with fixed radiusX, varying radiusY
 * 4. Third click: Complete oval with radiusY
 * 5. Escape: Cancel
 */
export class OvalTool implements ShapeTool {
  public readonly name = 'Oval';

  private center: HexCoord | null = null;
  private radiusXPoint: HexCoord | null = null;
  private radiusYPoint: HexCoord | null = null;
  private radiusX: number = 0;
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
    } else if (this.radiusXPoint === null) {
      // Second click: Store radiusX point and calculate radiusX
      this.radiusXPoint = { q, r };
      const centerRowCol = hexToRowCol(this.center.q, this.center.r);
      const pointRowCol = hexToRowCol(q, r);
      this.radiusX = Math.abs(pointRowCol.col - centerRowCol.col);
      return true; // More points needed
    } else {
      // Third click: Store radiusY point and complete
      this.radiusYPoint = { q, r };
      return false; // Shape complete
    }
  }

  getPreviewTiles(currentQ: number, currentR: number): HexCoord[] {
    if (this.center === null) {
      return []; // No preview until center is set
    }

    const centerRowCol = hexToRowCol(this.center.q, this.center.r);

    if (this.radiusXPoint === null) {
      // Preview radiusX as we move the mouse
      const currentRowCol = hexToRowCol(currentQ, currentR);
      const previewRadiusX = Math.abs(currentRowCol.col - centerRowCol.col);

      if (previewRadiusX === 0) {
        return [this.center]; // Show just center if radiusX is 0
      }

      // Show preview with equal radiusX and radiusY
      const tiles = this.world.ovalFrom(
        this.center.q,
        this.center.r,
        previewRadiusX,
        previewRadiusX, // Use same radius for circular preview
        false // Always outline for preview
      );

      return tiles.map(([q, r]) => ({ q, r }));
    } else {
      // RadiusX is fixed, preview radiusY
      const currentRowCol = hexToRowCol(currentQ, currentR);
      const previewRadiusY = Math.abs(currentRowCol.row - centerRowCol.row);

      if (this.radiusX === 0 && previewRadiusY === 0) {
        return [this.center];
      }

      // Show preview with fixed radiusX and varying radiusY
      const tiles = this.world.ovalFrom(
        this.center.q,
        this.center.r,
        Math.max(this.radiusX, 1), // Ensure at least 1
        Math.max(previewRadiusY, 1), // Ensure at least 1
        false // Always outline for preview
      );

      return tiles.map(([q, r]) => ({ q, r }));
    }
  }

  getResultTiles(): HexCoord[] {
    if (this.center === null || this.radiusXPoint === null || this.radiusYPoint === null) {
      return []; // No result if incomplete
    }

    // Calculate radiusY from the third point
    const centerRowCol = hexToRowCol(this.center.q, this.center.r);
    const radiusYPointRowCol = hexToRowCol(this.radiusYPoint.q, this.radiusYPoint.r);
    const radiusY = Math.abs(radiusYPointRowCol.row - centerRowCol.row);

    if (this.radiusX === 0 && radiusY === 0) {
      return [this.center]; // Single tile if both radii are 0
    }

    // Generate final oval with current fill setting
    const tiles = this.world.ovalFrom(
      this.center.q,
      this.center.r,
      Math.max(this.radiusX, 1),
      Math.max(radiusY, 1),
      this.filled
    );

    return tiles.map(([q, r]) => ({ q, r }));
  }

  getAnchorPoints(): HexCoord[] {
    const points: HexCoord[] = [];
    if (this.center !== null) {
      points.push(this.center);
    }
    if (this.radiusXPoint !== null) {
      points.push(this.radiusXPoint);
    }
    if (this.radiusYPoint !== null) {
      points.push(this.radiusYPoint);
    }
    return points;
  }

  reset(): void {
    this.center = null;
    this.radiusXPoint = null;
    this.radiusYPoint = null;
    this.radiusX = 0;
  }

  canComplete(): boolean {
    return this.center !== null && this.radiusXPoint !== null && this.radiusYPoint !== null;
  }

  requiresKeyboardConfirm(): boolean {
    return false; // Oval auto-completes after 3 clicks
  }

  getStatusText(): string {
    if (this.center === null) {
      return 'Click center of oval';
    } else if (this.radiusXPoint === null) {
      return 'Click to set horizontal radius (or press Escape to cancel)';
    } else if (this.radiusYPoint === null) {
      return 'Click to set vertical radius (or press Escape to cancel)';
    } else {
      return 'Oval complete';
    }
  }

  isFilled(): boolean {
    return this.filled;
  }

  setFilled(filled: boolean): void {
    this.filled = filled;
  }
}
