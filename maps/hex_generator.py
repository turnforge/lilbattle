#!/usr/bin/env python3
"""
Hex Generator

Generates systematic hex cell positions based on grid parameters with support
for both standard and inverted hex row offset patterns.

PURPOSE:
Takes GridParams (from grid_analyzer or manual calculation) and generates
the precise center coordinates for each hex cell in the grid. Handles the
characteristic "jagged" hex pattern where alternating rows are offset.

FEATURES:
- Standard hex pattern: Odd rows (1,3,5...) offset by hex_width/2
- Inverted hex pattern: Even rows (0,2,4...) offset by hex_width/2  
- Bounds checking: Only generates cells within image boundaries
- Debug visualization: Optional overlay showing hex positions
- Systematic naming: Each cell has (row, col) coordinates

HEX GRID GEOMETRY:
In a standard hex grid:
- Even rows (0,2,4...): Start at base X position
- Odd rows (1,3,5...): Offset by row_offset (typically hex_width/2)
- This creates the characteristic "jagged" hex tiling pattern

USAGE:
generator = HexCellGenerator(debug_mode=True)
hex_cells = generator.generate_hex_cells(image, params, invert_offset=False)
"""

import cv2
import numpy as np
from typing import List
from pathlib import Path
from dataclasses import dataclass
from grid_analyzer import GridParams


@dataclass
class HexCell:
    """Represents a hexagonal cell in the grid"""
    row: int
    col: int
    center_x: float
    center_y: float
    tile_id: int = 0
    confidence: float = 0.0


class HexCellGenerator:
    """Generates hex cell positions with configurable row offset patterns.
    
    This class converts GridParams into actual hex cell coordinates, handling
    the geometric calculations for proper hex grid layout. Supports both
    standard and inverted offset patterns for different hex grid orientations.
    
    KEY RESPONSIBILITIES:
    - Calculate precise hex center coordinates from grid parameters
    - Apply correct row offset pattern (standard or inverted)
    - Filter cells to ensure they fall within image boundaries  
    - Generate systematic (row, col) identifiers for each cell
    - Provide debug visualization of generated hex positions
    """
    
    def __init__(self, debug_mode: bool = False):
        self.debug_mode = debug_mode
        self.debug_dir = Path("debug_images") if debug_mode else None
        
        if self.debug_mode:
            self.debug_dir.mkdir(exist_ok=True)
    
    def generate_hex_cells(self, image: np.ndarray, params: GridParams, invert_offset: bool = False) -> List[HexCell]:
        """Generate hex cell positions with configurable row offset pattern.
        
        ALGORITHM:
        1. Iterate through each (row, col) position in the grid
        2. Calculate base X position: start_x + col * spacing_x  
        3. Apply row offset based on pattern:
           - Standard: Add row_offset to X for odd rows (1,3,5...)
           - Inverted: Add row_offset to X for even rows (0,2,4...)
        4. Calculate Y position: start_y + row * spacing_y
        5. Check if position is within image bounds
        6. Create HexCell object with calculated coordinates
        
        Args:
            image: Input image to check bounds against
            params: GridParams defining hex grid structure  
            invert_offset: If True, even rows are offset instead of odd rows
            
        Returns:
            List of HexCell objects with center coordinates
            
        ROW OFFSET PATTERNS:
        Standard (invert_offset=False):
        - Row 0: X = start_x + col * spacing_x
        - Row 1: X = start_x + col * spacing_x + row_offset  
        - Row 2: X = start_x + col * spacing_x
        - Row 3: X = start_x + col * spacing_x + row_offset
        
        Inverted (invert_offset=True):
        - Row 0: X = start_x + col * spacing_x + row_offset
        - Row 1: X = start_x + col * spacing_x
        - Row 2: X = start_x + col * spacing_x + row_offset  
        - Row 3: X = start_x + col * spacing_x
        """
        hex_cells = []
        
        print(f"Generating {params.rows} x {params.cols} = {params.rows * params.cols} hex cells")
        
        for row in range(params.rows):
            for col in range(params.cols):
                # Calculate hex center position using hex grid geometry
                x = params.start_x + col * params.spacing_x
                
                # CRITICAL: Apply row offset for hex pattern
                # This creates the characteristic "jagged" hex grid layout
                # 
                # WHY NEEDED: Hex grids have alternating row offsets to create efficient tiling
                # Standard WeeWar pattern: Even rows aligned, odd rows offset by hex_width/2
                # Some maps may use inverted pattern: Odd rows aligned, even rows offset
                #
                # IMPLEMENTATION:
                if invert_offset:
                    # INVERTED PATTERN: Even rows (0,2,4...) are offset
                    if row % 2 == 0:  
                        x += params.row_offset
                else:
                    # STANDARD PATTERN: Odd rows (1,3,5...) are offset  
                    if row % 2 == 1:  
                        x += params.row_offset
                
                y = params.start_y + row * params.spacing_y
                
                # Check if position is within image bounds
                if 0 <= x < image.shape[1] and 0 <= y < image.shape[0]:
                    hex_cell = HexCell(
                        row=row,
                        col=col,
                        center_x=x,
                        center_y=y,
                        tile_id=0,  # Will be classified later
                        confidence=0.0
                    )
                    hex_cells.append(hex_cell)
        
        if self.debug_mode:
            self._save_debug_hex_cells(image, hex_cells, self.debug_dir / "generated_cells.png")
        
        print(f"Generated {len(hex_cells)} valid hex cells")
        return hex_cells
    
    def _save_debug_hex_cells(self, image: np.ndarray, hex_cells: List[HexCell], path: Path):
        """Save debug image with hex cells marked"""
        debug_img = image.copy()
        
        for cell in hex_cells:
            # Draw center point
            cv2.circle(debug_img, (int(cell.center_x), int(cell.center_y)), 3, (0, 255, 0), -1)
            # Draw hex boundary circle
            cv2.circle(debug_img, (int(cell.center_x), int(cell.center_y)), 15, (255, 0, 0), 2)
            # Draw row/col text
            cv2.putText(debug_img, f"{cell.row},{cell.col}", 
                       (int(cell.center_x)-15, int(cell.center_y)-20), 
                       cv2.FONT_HERSHEY_SIMPLEX, 0.3, (0, 0, 255), 1)
        
        cv2.imwrite(str(path), debug_img)
        print(f"Debug image saved: {path}")


def main():
    """Generate hex grid from command line or test with default image"""
    import argparse
    from grid_analyzer import HexGridAnalyzer
    
    parser = argparse.ArgumentParser(description='Generate hex grid for WeeWar map images')
    parser.add_argument('--image', type=str, help='Path to the map image to analyze')
    parser.add_argument('--rows', type=int, help='Override number of rows (overrides detection)')
    parser.add_argument('--cols', type=int, help='Override number of columns (overrides detection)')
    parser.add_argument('--vert-spacing', type=float, help='Override vertical spacing in pixels (overrides detection)')
    parser.add_argument('--expected-tiles', type=int, default=34, help='Expected number of tiles in the map')
    parser.add_argument('--debug', action='store_true', help='Enable debug mode with visualization')
    
    args = parser.parse_args()
    
    # Use provided image path or default test image
    if args.image:
        image_path = args.image
    else:
        image_path = "../data/Maps/1_files/map-og.png"
        print(f"No image specified, using default: {image_path}")
    
    # Load image
    image = cv2.imread(image_path)
    
    if image is None:
        print(f"Could not load image: {image_path}")
        return
    
    print(f"Generating hex grid for: {image_path}")
    
    # First analyze grid structure
    analyzer = HexGridAnalyzer(debug_mode=args.debug)
    params = analyzer.analyze_grid_structure(image, expected_tiles=args.expected_tiles)
    
    if not params:
        print("Failed to analyze grid structure")
        return
    
    # Apply command-line overrides if provided
    if args.rows is not None:
        print(f"Overriding rows: {params.rows} -> {args.rows}")
        params.rows = args.rows
    
    if args.cols is not None:
        print(f"Overriding cols: {params.cols} -> {args.cols}")
        params.cols = args.cols
    
    if args.vert_spacing is not None:
        print(f"Overriding vertical spacing: {params.spacing_y:.1f} -> {args.vert_spacing}")
        params.spacing_y = args.vert_spacing
    
    print(f"Using grid parameters:")
    print(f"  Dimensions: {params.hex_width}x{params.hex_height}")
    print(f"  Grid size: {params.rows} rows x {params.cols} cols = {params.rows * params.cols} total")
    print(f"  Spacing: {params.spacing_x:.1f}x{params.spacing_y:.1f}")
    
    # Generate hex cells
    generator = HexCellGenerator(debug_mode=args.debug)
    hex_cells = generator.generate_hex_cells(image, params)
    
    print(f"Generated {len(hex_cells)} hex cells")
    if args.debug:
        for i, cell in enumerate(hex_cells[:10]):  # Show first 10
            print(f"  Cell {i}: row={cell.row}, col={cell.col}, pos=({cell.center_x:.1f}, {cell.center_y:.1f})")


if __name__ == "__main__":
    main()