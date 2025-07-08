#!/usr/bin/env python3
"""
Grid Analyzer

Analyzes hex grid structure from edge-detected images.
"""

import cv2
import numpy as np
import math
from typing import Optional, Dict, List
from pathlib import Path
from dataclasses import dataclass


@dataclass
class GridParams:
    """Parameters defining the hex grid structure"""
    hex_width: int          # Width of hex tile in pixels
    hex_height: int         # Height of hex tile in pixels  
    rows: int              # Number of rows
    cols: int              # Number of columns
    row_offset: float      # X offset for odd rows (0 or hex_width/2)
    start_x: int           # X coordinate of first hex center
    start_y: int           # Y coordinate of first hex center
    spacing_x: float       # Horizontal spacing between centers
    spacing_y: float       # Vertical spacing between centers


class HexGridAnalyzer:
    """Analyzes hex grid structure from edge detection"""
    
    def __init__(self, debug_mode: bool = False):
        self.debug_mode = debug_mode
        self.debug_dir = Path("debug_images") if debug_mode else None
        
        if self.debug_mode:
            self.debug_dir.mkdir(exist_ok=True)
    
    def analyze_grid_structure(self, image: np.ndarray, expected_tiles: int = 34) -> Optional[GridParams]:
        """Analyze hex grid structure from map boundary"""
        # Get edge image
        edges = self._get_edge_image(image)
        
        if self.debug_mode:
            cv2.imwrite(str(self.debug_dir / "structure_edges.png"), edges)
        
        # Find map boundaries
        boundaries = self._find_map_boundaries(edges)
        if not boundaries:
            print("Failed to find map boundaries")
            return None
        
        if self.debug_mode:
            print(f"Map boundaries: {boundaries}")
        
        # Calculate hex grid parameters from boundaries and expected tile count
        params = self._calculate_grid_from_boundaries(image, boundaries, expected_tiles)
        
        if self.debug_mode:
            print(f"Calculated grid params: {params}")
        
        return params
    
    def _get_edge_image(self, image: np.ndarray) -> np.ndarray:
        """Get edge-detected image"""
        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8,8))
        enhanced = clahe.apply(gray)
        edges = cv2.Canny(enhanced, 30, 90)
        return edges
    
    def _find_map_boundaries(self, edges: np.ndarray) -> Optional[Dict]:
        """Find map boundaries using 4-directional edge images and OR combination"""
        height, width = edges.shape
        
        # Get 4-directional boundary edge images
        projections = self._get_4_directional_projections(edges)
        
        if self.debug_mode:
            self._save_projection_debug_4dir(projections, height, width)
        
        # Combine all 4 edge images using OR operation to get outer boundary
        combined_boundary = np.zeros((height, width), dtype=np.uint8)
        for direction, edge_img in projections.items():
            combined_boundary = cv2.bitwise_or(combined_boundary, edge_img)
        
        if self.debug_mode:
            cv2.imwrite(str(self.debug_dir / "combined_boundary.png"), combined_boundary)
        
        # Find boundaries from the combined edge image
        boundaries = {}
        
        # Find the actual extent of the combined boundary
        coords = np.where(combined_boundary > 0)
        if len(coords[0]) == 0:
            print("No boundary pixels found")
            return None
        
        boundaries['top'] = np.min(coords[0])
        boundaries['bottom'] = np.max(coords[0])
        boundaries['left'] = np.min(coords[1])
        boundaries['right'] = np.max(coords[1])
        
        boundaries['height'] = boundaries['bottom'] - boundaries['top']
        boundaries['width'] = boundaries['right'] - boundaries['left']
        
        # Analyze hex tile size from the boundary edge images
        hex_info = self._analyze_hex_dimensions(projections)
        boundaries.update(hex_info)
        
        if self.debug_mode:
            self._save_boundary_debug(edges, boundaries)
            print(f"Pattern spacings by direction: {hex_info.get('pattern_spacings', {})}")
            print(f"Detected hex tile side length: {hex_info.get('hex_side_length', 'unknown')}")
        
        return boundaries
    
    def _get_4_directional_projections(self, edges: np.ndarray) -> Dict[str, np.ndarray]:
        """Get boundary edge images from 4 directions"""
        height, width = edges.shape
        
        projections = {}
        edge_thickness = 2  # Thickness of edge lines for better visibility
        
        # Create 4 separate edge images (same size as original)
        view_from_top = np.zeros((height, width), dtype=np.uint8)
        view_from_bottom = np.zeros((height, width), dtype=np.uint8)
        view_from_left = np.zeros((height, width), dtype=np.uint8)
        view_from_right = np.zeros((height, width), dtype=np.uint8)
        
        # View from top: for each column, mark the first edge pixel from top
        for col in range(width):
            column_data = edges[:, col]
            if np.any(column_data > 0):
                first_edge = np.argmax(column_data > 0)
                # Mark the edge pixel with some thickness
                for t in range(edge_thickness):
                    if first_edge + t < height:
                        view_from_top[first_edge + t, col] = 255
        
        # View from bottom: for each column, mark the first edge pixel from bottom
        for col in range(width):
            column_data = edges[:, col]
            if np.any(column_data > 0):
                last_edge = height - 1 - np.argmax(column_data[::-1] > 0)
                # Mark the edge pixel with some thickness
                for t in range(edge_thickness):
                    if last_edge - t >= 0:
                        view_from_bottom[last_edge - t, col] = 255
        
        # View from left: for each row, mark the first edge pixel from left
        for row in range(height):
            row_data = edges[row, :]
            if np.any(row_data > 0):
                first_edge = np.argmax(row_data > 0)
                # Mark the edge pixel with some thickness
                for t in range(edge_thickness):
                    if first_edge + t < width:
                        view_from_left[row, first_edge + t] = 255
        
        # View from right: for each row, mark the first edge pixel from right
        for row in range(height):
            row_data = edges[row, :]
            if np.any(row_data > 0):
                last_edge = width - 1 - np.argmax(row_data[::-1] > 0)
                # Mark the edge pixel with some thickness
                for t in range(edge_thickness):
                    if last_edge - t >= 0:
                        view_from_right[row, last_edge - t] = 255
        
        projections['view_from_top'] = view_from_top
        projections['view_from_bottom'] = view_from_bottom
        projections['view_from_left'] = view_from_left
        projections['view_from_right'] = view_from_right
        
        return projections
    
    def _analyze_hex_dimensions(self, projections: Dict[str, np.ndarray]) -> Dict:
        """Analyze hex dimensions from pattern spacing in edge images"""
        hex_info = {}
        
        # Convert 2D edge images to 1D profiles for pattern analysis
        pattern_spacings = {}
        
        for direction, edge_img in projections.items():
            if direction in ['view_from_top', 'view_from_bottom']:
                # For top/bottom views, sum along vertical axis to get horizontal profile
                profile = np.sum(edge_img, axis=0)
            else:  # left/right views
                # For left/right views, sum along horizontal axis to get vertical profile
                profile = np.sum(edge_img, axis=1)
            
            spacing = self._find_pattern_spacing(profile)
            pattern_spacings[direction] = spacing
        
        # Use horizontal projections to determine hex spacing
        horizontal_spacings = [
            pattern_spacings.get('view_from_top', 0), 
            pattern_spacings.get('view_from_bottom', 0)
        ]
        vertical_spacings = [
            pattern_spacings.get('view_from_left', 0), 
            pattern_spacings.get('view_from_right', 0)
        ]
        
        # Take median of detected spacings (more robust than max)
        all_spacings = [s for s in horizontal_spacings + vertical_spacings if s > 10]
        
        if all_spacings:
            hex_side_length = int(np.median(all_spacings))
        else:
            hex_side_length = 60  # Fallback estimate
        
        hex_info['hex_side_length'] = hex_side_length
        hex_info['pattern_spacings'] = pattern_spacings
        
        return hex_info
    
    def _find_pattern_spacing(self, projection: np.ndarray) -> int:
        """Find the repeating pattern spacing in a projection"""
        if len(projection) < 10:
            return 0
        
        # Skip if projection is all zeros
        if np.max(projection) == 0:
            return 0
        
        # Look for peaks and valleys in the projection to find pattern spacing
        from scipy.signal import find_peaks
        
        # For sparse edge data, don't smooth too much - preserve the edge positions
        from scipy.ndimage import gaussian_filter1d
        smoothed = gaussian_filter1d(projection.astype(float), sigma=1)
        
        # Lower the threshold and distance for sparse edge data
        max_val = np.max(smoothed)
        if max_val == 0:
            return 0
            
        # Find peaks with lower threshold for sparse data
        peaks, _ = find_peaks(smoothed, height=max_val * 0.1, distance=10)
        
        if len(peaks) < 2:
            # Try finding any non-zero positions as potential peaks
            nonzero_positions = np.where(projection > 0)[0]
            if len(nonzero_positions) >= 2:
                # Use the spacing between non-zero regions
                spacings = np.diff(nonzero_positions)
                # Filter out very small spacings (likely same feature)
                valid_spacings = spacings[spacings > 5]
                if len(valid_spacings) > 0:
                    return int(np.median(valid_spacings))
            return 0
        
        # Calculate spacing between peaks
        peak_spacings = np.diff(peaks)
        
        if len(peak_spacings) > 0:
            # Return median spacing (most common hex spacing)
            return int(np.median(peak_spacings))
        
        return 0
    
    def _find_longest_continuous_line(self, projection: np.ndarray) -> int:
        """Find the length of the longest continuous non-zero segment"""
        if len(projection) == 0:
            return 0
        
        max_length = 0
        current_length = 0
        
        for value in projection:
            if value > 0:
                current_length += 1
                max_length = max(max_length, current_length)
            else:
                current_length = 0
        
        return max_length
    
    def _calculate_grid_from_boundaries(self, image: np.ndarray, boundaries: Dict, expected_tiles: int) -> GridParams:
        """Calculate hex grid parameters using detected hex side length"""
        
        map_width = boundaries['width']
        map_height = boundaries['height']
        hex_side_length = boundaries.get('hex_side_length', 0)
        
        if hex_side_length <= 10:  # If detection failed, use fallback
            print(f"Warning: hex side length detection failed ({hex_side_length}), using fallback")
            return self._fallback_grid_calculation(boundaries, expected_tiles)
        
        print(f"Using detected hex side length: {hex_side_length}")
        
        # Calculate hex dimensions from side length
        # For a regular hexagon: width ≈ 2 * side_length, height ≈ 1.73 * side_length
        hex_width = int(hex_side_length * 2)
        hex_height = int(hex_side_length * 1.73)  # sqrt(3) ≈ 1.73
        
        # Calculate how many hexes fit in the map
        # Use the detected pattern spacing directly instead of calculating from hex dimensions
        avg_horizontal_spacing = np.mean([boundaries['pattern_spacings'].get('view_from_top', hex_side_length), 
                                         boundaries['pattern_spacings'].get('view_from_bottom', hex_side_length)])
        avg_vertical_spacing = np.mean([boundaries['pattern_spacings'].get('view_from_left', hex_side_length), 
                                       boundaries['pattern_spacings'].get('view_from_right', hex_side_length)])
        
        # Use the detected spacings directly, with fallback if detection failed
        spacing_x = avg_horizontal_spacing if avg_horizontal_spacing > 0 else hex_side_length * 1.5
        spacing_y = avg_vertical_spacing if avg_vertical_spacing > 0 else hex_side_length * 1.3
        
        cols = int(map_width / spacing_x) + 1 if spacing_x > 0 else 7
        rows = int(map_height / spacing_y) + 1 if spacing_y > 0 else 7
        
        # Adjust if the grid is too large compared to expected tiles
        total_positions = rows * cols
        if total_positions > expected_tiles * 2:  # Too many positions
            # Try smaller grid
            cols = max(5, int(np.sqrt(expected_tiles * 1.5)))
            rows = max(5, int(np.sqrt(expected_tiles * 1.5)))
            
            # Recalculate spacing based on desired grid size
            spacing_x = map_width / cols
            spacing_y = map_height / rows
            hex_width = int(spacing_x * 1.33)  # Reverse calculation
            hex_height = int(spacing_y * 1.15)
        
        # Calculate starting positions (center of first hex)
        start_x = boundaries['left'] + hex_width // 2
        start_y = boundaries['top'] + hex_height // 2
        
        # Row offset for hex pattern (odd rows offset by half spacing)
        row_offset = spacing_x // 2
        
        print(f"Calculated grid: {rows}x{cols} = {rows * cols} positions")
        print(f"Hex dimensions: {hex_width}x{hex_height}")
        print(f"Spacing: {spacing_x:.1f}x{spacing_y:.1f}")
        
        return GridParams(
            hex_width=hex_width,
            hex_height=hex_height,
            rows=rows,
            cols=cols,
            row_offset=row_offset,
            start_x=start_x,
            start_y=start_y,
            spacing_x=spacing_x,
            spacing_y=spacing_y
        )
    
    def _fallback_grid_calculation(self, boundaries: Dict, expected_tiles: int) -> GridParams:
        """Fallback calculation when hex side length detection fails"""
        map_width = boundaries['width']
        map_height = boundaries['height']
        
        # Use square root approximation
        approx_side = int(np.sqrt(expected_tiles * 1.4))  # Slightly larger for hex shape
        
        rows = approx_side
        cols = approx_side
        
        spacing_x = map_width / cols
        spacing_y = map_height / rows
        
        hex_width = int(spacing_x)
        hex_height = int(spacing_y)
        
        start_x = boundaries['left'] + hex_width // 2
        start_y = boundaries['top'] + hex_height // 2
        row_offset = spacing_x // 2
        
        return GridParams(
            hex_width=hex_width,
            hex_height=hex_height,
            rows=rows,
            cols=cols,
            row_offset=row_offset,
            start_x=start_x,
            start_y=start_y,
            spacing_x=spacing_x,
            spacing_y=spacing_y
        )
    
    def _save_boundary_debug(self, edges: np.ndarray, boundaries: Dict):
        """Save debug image showing detected boundaries"""
        height, width = edges.shape
        
        # Create RGB image for better visualization
        debug_img = cv2.cvtColor(edges, cv2.COLOR_GRAY2BGR)
        
        # Draw boundary lines
        cv2.line(debug_img, (0, boundaries['top']), (width, boundaries['top']), (0, 255, 0), 2)  # Top - green
        cv2.line(debug_img, (0, boundaries['bottom']), (width, boundaries['bottom']), (0, 255, 0), 2)  # Bottom - green
        cv2.line(debug_img, (boundaries['left'], 0), (boundaries['left'], height), (255, 0, 0), 2)  # Left - blue
        cv2.line(debug_img, (boundaries['right'], 0), (boundaries['right'], height), (255, 0, 0), 2)  # Right - blue
        
        # Draw bounding box
        cv2.rectangle(debug_img, 
                     (boundaries['left'], boundaries['top']), 
                     (boundaries['right'], boundaries['bottom']), 
                     (0, 0, 255), 2)  # Red rectangle
        
        # Add text with dimensions
        cv2.putText(debug_img, f"W: {boundaries['width']}, H: {boundaries['height']}", 
                   (10, 30), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        
        cv2.imwrite(str(self.debug_dir / "map_boundaries.png"), debug_img)
    
    def _save_projection_debug_4dir(self, projections: Dict[str, np.ndarray], height: int, width: int):
        """Save debug visualization of 4-directional edge images"""
        # Save each individual edge image for clear visualization
        for direction, edge_img in projections.items():
            filename = f"edge_{direction}.png"
            cv2.imwrite(str(self.debug_dir / filename), edge_img)
        
        # Create a combined RGB visualization where each direction gets a color channel
        combined_img = np.zeros((height, width, 3), dtype=np.uint8)
        
        # Assign colors to each direction for the combined view
        if 'view_from_top' in projections:
            combined_img[:, :, 1] = projections['view_from_top']  # Green channel
        if 'view_from_bottom' in projections:
            combined_img[:, :, 0] = projections['view_from_bottom']  # Blue channel  
        if 'view_from_left' in projections:
            combined_img[:, :, 2] = projections['view_from_left']  # Red channel
        if 'view_from_right' in projections:
            # Combine with green channel (will appear cyan where overlapping)
            combined_img[:, :, 1] = cv2.bitwise_or(combined_img[:, :, 1], projections['view_from_right'])
        
        cv2.imwrite(str(self.debug_dir / "4dir_edges_combined.png"), combined_img)
        
        # Also create a simple grayscale combined view (OR of all edges)
        combined_gray = np.zeros((height, width), dtype=np.uint8)
        for direction, edge_img in projections.items():
            combined_gray = cv2.bitwise_or(combined_gray, edge_img)
        
        cv2.imwrite(str(self.debug_dir / "4dir_edges_gray.png"), combined_gray)
    
    def _save_projection_debug(self, projection: np.ndarray, direction: str, height: int, width: int, transpose: bool = False):
        """Save debug visualization of projection"""
        if direction == "horizontal":
            proj_img = np.zeros((height, width), dtype=np.uint8)
            for y, value in enumerate(projection):
                line_width = int((value / np.max(projection)) * width) if np.max(projection) > 0 else 0
                proj_img[y, :line_width] = 255
        else:  # vertical
            proj_img = np.zeros((height, width), dtype=np.uint8)
            for x, value in enumerate(projection):
                line_height = int((value / np.max(projection)) * height) if np.max(projection) > 0 else 0
                proj_img[-line_height:, x] = 255
        
        cv2.imwrite(str(self.debug_dir / f"{direction}_projection.png"), proj_img)


def main():
    """Test the grid analyzer"""
    # Load test image
    image_path = "../data/Maps/1_files/map-og.png"
    image = cv2.imread(image_path)
    
    if image is None:
        print(f"Could not load image: {image_path}")
        return
    
    # Analyze grid with expected tile count
    analyzer = HexGridAnalyzer(debug_mode=True)
    params = analyzer.analyze_grid_structure(image, expected_tiles=34)
    
    if params:
        print(f"Successfully analyzed grid structure:")
        print(f"  Dimensions: {params.hex_width}x{params.hex_height}")
        print(f"  Grid size: {params.rows} rows x {params.cols} cols = {params.rows * params.cols} total")
        print(f"  Spacing: {params.spacing_x:.1f}x{params.spacing_y:.1f}")
        print(f"  Row offset: {params.row_offset:.1f}")
        print(f"  Start position: ({params.start_x}, {params.start_y})")
    else:
        print("Failed to analyze grid structure")


if __name__ == "__main__":
    main()