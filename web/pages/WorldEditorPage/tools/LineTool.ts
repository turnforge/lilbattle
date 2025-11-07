import { ShapeTool, HexCoord } from './ShapeTool';
import { World } from '../../common/World';

/**
 * Line/Path drawing tool.
 *
 * Workflow:
 * 1. Each click: Add vertex to path
 * 2. Mouse move: Show line segment from last vertex to cursor
 * 3. Enter: Confirm and draw the complete path
 * 4. Escape: Cancel entire path
 *
 * Note: Fill toggle is N/A for lines (always stroke only)
 */
export class LineTool implements ShapeTool {
  public readonly name = 'Line';

  private points: HexCoord[] = [];
  private world: World;

  constructor(world: World) {
    this.world = world;
  }

  addPoint(q: number, r: number): boolean {
    // Add point to the path
    this.points.push({ q, r });

    // Line tool always needs more points (completed via Enter key)
    return true;
  }

  getPreviewTiles(currentQ: number, currentR: number): HexCoord[] {
    if (this.points.length === 0) {
      return []; // No preview until first point is set
    }

    // Build preview path: all collected points + current mouse position
    const previewPoints = [...this.points, { q: currentQ, r: currentR }];

    // Generate line through all preview points
    const tiles = this.world.lineFrom(previewPoints);

    return tiles.map(([q, r]) => ({ q, r }));
  }

  getResultTiles(): HexCoord[] {
    if (this.points.length < 2) {
      // Need at least 2 points for a line
      return this.points;
    }

    // Generate final line through all collected points
    const tiles = this.world.lineFrom(this.points);

    return tiles.map(([q, r]) => ({ q, r }));
  }

  getAnchorPoints(): HexCoord[] {
    return [...this.points];
  }

  reset(): void {
    this.points = [];
  }

  canComplete(): boolean {
    // Need at least 2 points to complete a line
    return this.points.length >= 2;
  }

  requiresKeyboardConfirm(): boolean {
    return true; // Line requires Enter key to complete
  }

  getStatusText(): string {
    if (this.points.length === 0) {
      return 'Click to start line/path';
    } else if (this.points.length === 1) {
      return 'Click to add points, Enter to finish, Escape to cancel';
    } else {
      return `${this.points.length} points - Click to add more, Enter to finish, Escape to cancel`;
    }
  }

  isFilled(): boolean {
    return false; // Lines are never filled
  }

  setFilled(filled: boolean): void {
    // No-op: Lines don't support fill mode
  }
}
