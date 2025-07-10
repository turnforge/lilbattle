#!/usr/bin/env python3
"""
Hex Splitter

Extracts individual hex tiles from WeeWar map images as separate PNG files
with transparent backgrounds to prevent neighbor tile bleeding.

FEATURES:
- Automatic mode: Uses grid_analyzer to detect hex grid structure
- Manual mode: Calculates grid parameters from user-provided dimensions  
- Row offset control: Handles both standard and inverted hex grid patterns
- Hexagonal masking: Creates clean tiles with transparent backgrounds

USAGE MODES:

1. AUTOMATIC MODE (uses grid analyzer):
   python hex_splitter.py --image map.png
   
2. MANUAL MODE (independent operation):
   python hex_splitter.py --image map.png --rows 7 --cols 7 --vert-spacing 53.5
   
3. ROW OFFSET CONTROL:
   --invert-offset    # Even rows offset instead of odd rows

ROW OFFSET PATTERNS:
- Standard (default): Row 0,2,4... at base X, Row 1,3,5... offset by hex_width/2
- Inverted: Row 0,2,4... offset by hex_width/2, Row 1,3,5... at base X

MANUAL MODE BENEFITS:
- Works when automatic detection fails
- No dependency on grid_analyzer accuracy
- Complete user control over grid parameters
- Faster execution (skips analysis step)

OUTPUT:
Creates individual hex tile files (0_0.png, 1_2.png, etc.) with:
- Proper hexagonal boundaries (flat-top orientation)
- Transparent backgrounds outside hex shape
- Consistent naming: {row}_{col}.png
"""

import cv2
import numpy as np
import argparse
from typing import List
from pathlib import Path
from dataclasses import dataclass
from grid_analyzer import HexGridAnalyzer, GridParams
from hex_generator import HexCellGenerator, HexCell


class HexSplitter:
    """Splits hex grid images into individual tile images with hexagonal masking.
    
    This class handles the extraction of individual hex tiles from a complete hex grid
    image. It supports both automatic detection (via grid_analyzer) and manual mode
    (user-provided parameters) with configurable row offset patterns.
    
    KEY FEATURES:
    - Hexagonal masking with transparent backgrounds
    - Flat-top hex orientation (WeeWar standard)
    - Configurable row offset patterns (standard/inverted)
    - Bounds checking to ensure tiles are within image
    - Debug mode for detailed extraction logging
    """
    
    def __init__(self, output_dir: str = "hex_tiles", debug_mode: bool = False):
        self.output_dir = Path(output_dir)
        self.debug_mode = debug_mode
        
        # Create output directory
        self.output_dir.mkdir(exist_ok=True)
        
        if self.debug_mode:
            print(f"Output directory: {self.output_dir}")
    
    def split_hex_tiles(self, image: np.ndarray, params: GridParams, invert_offset: bool = False) -> List[str]:
        """Split the image into individual hex tiles with proper hex positioning.
        
        PROCESS:
        1. Generate hex cell positions using HexCellGenerator
        2. For each hex cell, extract a tile region with margin
        3. Create hexagonal mask for clean tile boundaries  
        4. Apply mask to create transparent background
        5. Save as individual PNG files with row_col naming
        
        Args:
            image: Input hex grid image to split
            params: GridParams defining hex grid structure
            invert_offset: If True, even rows are offset instead of odd rows
            
        Returns:
            List of file paths for extracted hex tiles
            
        ROW OFFSET BEHAVIOR:
        - invert_offset=False: Rows 0,2,4... at base X, rows 1,3,5... offset
        - invert_offset=True: Rows 0,2,4... offset, rows 1,3,5... at base X
        """
        
        # First generate hex cell positions
        generator = HexCellGenerator(debug_mode=False)
        hex_cells = generator.generate_hex_cells(image, params, invert_offset=invert_offset)
        
        if self.debug_mode:
            print(f"Splitting {len(hex_cells)} hex cells")
        
        # Extract each hex tile
        extracted_files = []
        for cell in hex_cells:
            tile_filename = self._extract_hex_tile(image, cell, params)
            if tile_filename:
                extracted_files.append(tile_filename)
        
        print(f"Extracted {len(extracted_files)} hex tiles to {self.output_dir}")
        return extracted_files
    
    def _extract_hex_tile(self, image: np.ndarray, cell: HexCell, params: GridParams) -> str:
        """Extract a single hex tile with transparent background"""
        
        # Calculate extraction region (slightly larger than hex to ensure full coverage)
        margin = 5
        tile_size = max(params.hex_width, params.hex_height) + 2 * margin
        
        # Calculate bounding box
        x_start = int(cell.center_x - tile_size // 2)
        y_start = int(cell.center_y - tile_size // 2)
        x_end = x_start + tile_size
        y_end = y_start + tile_size
        
        # Ensure bounds are within image
        height, width = image.shape[:2]
        x_start = max(0, x_start)
        y_start = max(0, y_start)
        x_end = min(width, x_end)
        y_end = min(height, y_end)
        
        if x_end <= x_start or y_end <= y_start:
            if self.debug_mode:
                print(f"Skipping cell {cell.row},{cell.col} - invalid bounds")
            return None
        
        # Extract the region
        tile_region = image[y_start:y_end, x_start:x_end]
        
        if tile_region.size == 0:
            return None
        
        # Create hexagonal mask
        # Use a larger radius to ensure we capture the full tile
        mask_radius = max(params.hex_width, params.hex_height) // 2 * 0.95  # Slightly smaller than extraction region
        mask = self._create_hex_mask(tile_region.shape, 
                                   cell.center_x - x_start, 
                                   cell.center_y - y_start,
                                   mask_radius)
        
        # Apply mask to create transparent background
        tile_with_alpha = self._apply_hex_mask(tile_region, mask)
        
        # Save the tile
        filename = f"{cell.row:02d}_{cell.col:02d}.png"
        output_path = self.output_dir / filename
        
        cv2.imwrite(str(output_path), tile_with_alpha)
        
        if self.debug_mode:
            print(f"Extracted tile {filename}: center=({cell.center_x:.1f}, {cell.center_y:.1f})")
        
        return str(output_path)
    
    def _create_hex_mask(self, shape: tuple, center_x: float, center_y: float, radius: float) -> np.ndarray:
        """Create a hexagonal mask for clean tile extraction"""
        height, width = shape[:2]
        mask = np.zeros((height, width), dtype=np.uint8)
        
        # Create actual hexagon vertices for flat-top hexagon (WeeWar style)
        # Start at top-right and go clockwise, offset by 30 degrees for flat-top
        hex_points = []
        for i in range(6):
            angle = (i * np.pi / 3) + (np.pi / 6)  # Add 30 degrees offset for flat-top
            x = center_x + radius * np.cos(angle)
            y = center_y + radius * np.sin(angle)
            hex_points.append([int(x), int(y)])
        
        # Convert to numpy array for OpenCV
        hex_points = np.array(hex_points, dtype=np.int32)
        
        # Fill the hexagonal region
        cv2.fillPoly(mask, [hex_points], 255)
        
        return mask
    
    def _apply_hex_mask(self, tile_region: np.ndarray, mask: np.ndarray) -> np.ndarray:
        """Apply hexagonal mask to create transparent background"""
        height, width = tile_region.shape[:2]
        
        # Create RGBA image (BGR + Alpha channel)
        if len(tile_region.shape) == 3:
            # Color image
            tile_rgba = cv2.cvtColor(tile_region, cv2.COLOR_BGR2BGRA)
        else:
            # Grayscale image
            tile_rgba = cv2.cvtColor(tile_region, cv2.COLOR_GRAY2BGRA)
        
        # Apply mask to alpha channel (255 = opaque, 0 = transparent)
        tile_rgba[:, :, 3] = mask
        
        return tile_rgba


def main():
    """Extract individual hex tiles from command line"""
    parser = argparse.ArgumentParser(description='Extract individual hex tiles from WeeWar map images')
    parser.add_argument('--image', type=str, required=True, help='Path to the map image to split')
    parser.add_argument('--output-dir', type=str, default='hex_tiles', help='Output directory for extracted tiles')
    parser.add_argument('--rows', type=int, help='Override number of rows (overrides detection)')
    parser.add_argument('--cols', type=int, help='Override number of columns (overrides detection)')
    parser.add_argument('--vert-spacing', type=float, help='Override vertical spacing in pixels (overrides detection)')
    parser.add_argument('--invert-offset', action='store_true', help='Invert hex row offset pattern (even rows offset instead of odd rows)')
    parser.add_argument('--expected-tiles', type=int, default=34, help='Expected number of tiles in the map')
    parser.add_argument('--debug', action='store_true', help='Enable debug mode with verbose output')
    
    args = parser.parse_args()
    
    # Load image
    image = cv2.imread(args.image)
    
    if image is None:
        print(f"Could not load image: {args.image}")
        return
    
    print(f"Splitting hex tiles from: {args.image}")
    
    # Check if we can run in manual mode (both rows and cols provided)
    if args.rows is not None and args.cols is not None:
        print("Running in manual mode - calculating GridParams from provided parameters")
        
        # MANUAL MODE: Calculate GridParams directly from user input without grid analyzer
        # This provides a fallback when automatic detection fails or for precise control
        image_height, image_width = image.shape[:2]
        
        # STEP 1: Calculate hex_width from image width and column count
        # Assumption: hex centers are evenly distributed across image width
        # For 7 columns in 448px image: hex_width = 448 / 7 = 64px
        hex_width = int(image_width / args.cols)
        
        # STEP 2: Determine vertical spacing and hex_height
        if args.vert_spacing is not None:
            # User provided vertical spacing - use directly
            spacing_y = args.vert_spacing
            # Reverse-calculate hex_height from spacing (typical: spacing = hex_height * 0.75)
            hex_height = int(spacing_y / 0.75)
        else:
            # Auto-calculate from image height and row count
            # Distribute image height evenly across rows
            spacing_y = image_height / args.rows if args.rows > 1 else image_height
            hex_height = int(spacing_y / 0.75)
        
        # STEP 3: Calculate grid positioning parameters
        spacing_x = hex_width          # Horizontal center-to-center distance
        row_offset = hex_width / 2     # Offset for alternating rows (standard hex pattern)
        start_x = hex_width / 2        # Start at center of first hex (not edge)
        start_y = hex_height / 2       # Start at center of first hex (not edge)
        
        params = GridParams(
            hex_width=hex_width,
            hex_height=hex_height,
            rows=args.rows,
            cols=args.cols,
            row_offset=row_offset,
            start_x=int(start_x),
            start_y=int(start_y),
            spacing_x=spacing_x,
            spacing_y=spacing_y
        )
        
        print(f"Calculated manual GridParams:")
        print(f"  Hex dimensions: {hex_width}x{hex_height}")
        print(f"  Spacing: {spacing_x:.1f}x{spacing_y:.1f}")
        print(f"  Row offset: {row_offset:.1f}")
        print(f"  Start position: ({start_x:.1f}, {start_y:.1f})")
        
    else:
        print("Running in automatic mode - using grid analyzer")
        
        # Analyze grid structure automatically
        analyzer = HexGridAnalyzer(debug_mode=args.debug)
        params = analyzer.analyze_grid_structure(image, expected_tiles=args.expected_tiles)
        
        if not params:
            print("Failed to analyze grid structure")
            print("Try running in manual mode with --rows and --cols parameters")
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
    print(f"  Grid size: {params.rows} rows x {params.cols} cols = {params.rows * params.cols} total")
    print(f"  Spacing: {params.spacing_x:.1f}x{params.spacing_y:.1f}")
    print(f"  Row offset pattern: {'Inverted (even rows offset)' if args.invert_offset else 'Standard (odd rows offset)'}")
    
    # Split into individual tiles
    splitter = HexSplitter(output_dir=args.output_dir, debug_mode=args.debug)
    extracted_files = splitter.split_hex_tiles(image, params, invert_offset=args.invert_offset)
    
    print(f"\nSuccessfully extracted {len(extracted_files)} hex tiles")
    print(f"Output directory: {args.output_dir}")


if __name__ == "__main__":
    main()
