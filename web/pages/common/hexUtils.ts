/**
 * Hex coordinate utility functions
 * These match the Go implementation from lib/map.go
 */

export interface HexCoord {
    q: number;
    r: number;
}

export interface PixelCoord {
    x: number;
    y: number;
}

// Tile constants matching the Go implementation
export const TILE_WIDTH = 64;
export const TILE_HEIGHT = 64;
export const Y_INCREMENT = 48;

/**
 * Convert hex coordinates to pixel coordinates
 * Matches lib/map.go CenterXYForTile
 */
export function hexToPixel(q: number, r: number, tileWidth=TILE_WIDTH, tileHeight=TILE_HEIGHT, yIncrement=Y_INCREMENT): PixelCoord {
  // Match the Go implementation from map.go CenterXYForTile
  const { row, col } = hexToRowCol(q, r);

  let y = yIncrement * row;
  let x = tileWidth * col;

  if ((row & 1) === 1) {
    x += tileWidth / 2;
  }
  return { x, y };
}

/**
 * Convert pixel coordinates to hex coordinates
 * Matches lib/map.go XYToQR
 */
export function pixelToHex(x: number, y: number, tileWidth=TILE_WIDTH, tileHeight=TILE_HEIGHT, yIncrement=Y_INCREMENT): HexCoord {
    // Match the Go implementation from map.go XYToQR
  const row = Math.floor((y + tileHeight / 2) / yIncrement);
  let halfDists = Math.floor(1 + Math.abs(x * 2 / tileWidth));
  if ((row & 1) !== 0) {
    halfDists = Math.floor(1 + Math.abs((x - tileWidth / 2) * 2 / tileWidth));
  }

  let col = Math.floor(halfDists / 2);
  if (x < 0) {
    col = -col;
  }

  return rowColToHex(row, col);
}

/**
 * Convert row/col coordinates to hex coordinates
 * RowColToHex: oddr_to_cube conversion
 */
export function rowColToHex(row: number, col: number): HexCoord {
    const x = col - Math.floor((row - (row & 1)) / 2);
    const z = row;
    const q = x;
    const r = z;
    return { q, r };
}

/**
 * Convert hex coordinates to row/col coordinates
 * HexToRowCol: cube_to_oddr conversion
 */
export function hexToRowCol(q: number, r: number): { row: number; col: number } {
    const row = r;
    const col = q + Math.floor((r - (r & 1)) / 2);
    return { row, col };
}


export const AxialNeighborDeltas = [
	{q: -1, r: 0}, // LEFT
	{q: 0, r: -1}, // TOP_LEFT
	{q: 1, r: -1}, // TOP_RIGHT
	{q: 1, r: 0},  // RIGHT
	{q: 0, r: 1},  // BOTTOM_RIGHT
	{q: -1, r: 1}, // BOTTOM_LEFT
]

export function axialNeighbors(q: number, r: number): [number, number][] {
  let out = [] as any;
	for (var i = 0;i < 6;i++) {
    out.push([q + AxialNeighborDeltas[i].q, r + AxialNeighborDeltas[i].r])
	}
  return out
}

/**
 * Calculate hex distance using cube coordinates
 * Distance is the maximum of the absolute differences in cube coordinates
 */
export function hexDistance(q1: number, r1: number, q2: number, r2: number): number {
    const dq = q2 - q1;
    const dr = r2 - r1;
    return Math.max(Math.abs(dq), Math.abs(dr), Math.abs(dq + dr));
}

/**
 * Get the direction index (0-5) from one hex to a neighbor hex
 * Returns null if the hexes are not neighbors
 *
 * Direction indices match AxialNeighborDeltas:
 * 0: LEFT (-1,0), 1: TOP_LEFT (0,-1), 2: TOP_RIGHT (1,-1),
 * 3: RIGHT (1,0), 4: BOTTOM_RIGHT (0,1), 5: BOTTOM_LEFT (-1,1)
 */
export function getDirectionIndex(fromQ: number, fromR: number, toQ: number, toR: number): number | null {
    const dq = toQ - fromQ;
    const dr = toR - fromR;
    for (let i = 0; i < 6; i++) {
        if (AxialNeighborDeltas[i].q === dq && AxialNeighborDeltas[i].r === dr) {
            return i;
        }
    }
    return null;
}

/**
 * Get the opposite direction index
 * 0 (LEFT) <-> 3 (RIGHT)
 * 1 (TOP_LEFT) <-> 4 (BOTTOM_RIGHT)
 * 2 (TOP_RIGHT) <-> 5 (BOTTOM_LEFT)
 */
export function getOppositeDirection(dirIndex: number): number {
    return (dirIndex + 3) % 6;
}

/**
 * Check if two hexes are immediate neighbors
 */
export function areNeighbors(q1: number, r1: number, q2: number, r2: number): boolean {
    return getDirectionIndex(q1, r1, q2, r2) !== null;
}

/**
 * Get the neighbor coordinate in a given direction
 */
export function getNeighborCoord(q: number, r: number, direction: number): [number, number] {
    const delta = AxialNeighborDeltas[direction];
    return [q + delta.q, r + delta.r];
}
