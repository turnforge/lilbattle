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
        """Find map boundaries using 4-directional projections to get outer boundary only"""
        height, width = edges.shape
        
        # Get 4-directional boundary projections
        projections = self._get_4_directional_projections(edges)
        
        if self.debug_mode:
            self._save_projection_debug_4dir(projections, height, width)
        
        # Find boundaries from projections
        boundaries = {}
        
        # Top boundary: first non-zero in view_from_top
        top_profile = projections['view_from_top']
        boundaries['top'] = np.argmax(top_profile > 0) if np.any(top_profile > 0) else 0
        
        # Bottom boundary: last non-zero in view_from_bottom (reversed)
        bottom_profile = projections['view_from_bottom']
        boundaries['bottom'] = height - 1 - np.argmax(bottom_profile[::-1] > 0) if np.any(bottom_profile > 0) else height - 1
        
        # Left boundary: first non-zero in view_from_left
        left_profile = projections['view_from_left']
        boundaries['left'] = np.argmax(left_profile > 0) if np.any(left_profile > 0) else 0
        
        # Right boundary: last non-zero in view_from_right (reversed)
        right_profile = projections['view_from_right']
        boundaries['right'] = width - 1 - np.argmax(right_profile[::-1] > 0) if np.any(right_profile > 0) else width - 1
        
        boundaries['height'] = boundaries['bottom'] - boundaries['top']
        boundaries['width'] = boundaries['right'] - boundaries['left']
        
        # Analyze hex tile size from longest continuous lines
        hex_info = self._analyze_hex_dimensions(projections)
        boundaries.update(hex_info)
        
        if self.debug_mode:
            self._save_boundary_debug(edges, boundaries)
            print(f"Detected hex tile side length: {hex_info.get('hex_side_length', 'unknown')}")
        
        return boundaries
    
    def _get_4_directional_projections(self, edges: np.ndarray) -> Dict[str, np.ndarray]:
        """Get boundary projections from 4 directions"""
        height, width = edges.shape
        
        projections = {}
        
        # View from top: for each column, find the first edge pixel from top
        view_from_top = np.zeros(width)
        for col in range(width):
            column_data = edges[:, col]
            first_edge = np.argmax(column_data > 0) if np.any(column_data > 0) else 0
            view_from_top[col] = height - first_edge if first_edge > 0 else 0
        
        # View from bottom: for each column, find the first edge pixel from bottom
        view_from_bottom = np.zeros(width)
        for col in range(width):
            column_data = edges[:, col]
            last_edge = height - 1 - np.argmax(column_data[::-1] > 0) if np.any(column_data > 0) else height - 1
            view_from_bottom[col] = last_edge if last_edge < height - 1 else 0
        
        # View from left: for each row, find the first edge pixel from left
        view_from_left = np.zeros(height)
        for row in range(height):
            row_data = edges[row, :]
            first_edge = np.argmax(row_data > 0) if np.any(row_data > 0) else 0
            view_from_left[row] = width - first_edge if first_edge > 0 else 0
        
        # View from right: for each row, find the first edge pixel from right
        view_from_right = np.zeros(height)
        for row in range(height):
            row_data = edges[row, :]
            last_edge = width - 1 - np.argmax(row_data[::-1] > 0) if np.any(row_data > 0) else width - 1
            view_from_right[row] = last_edge if last_edge < width - 1 else 0
        
        projections['view_from_top'] = view_from_top
        projections['view_from_bottom'] = view_from_bottom
        projections['view_from_left'] = view_from_left
        projections['view_from_right'] = view_from_right
        
        return projections
    
    def _analyze_hex_dimensions(self, projections: Dict[str, np.ndarray]) -> Dict:
        """Analyze hex dimensions from pattern spacing in projections"""
        hex_info = {}
        
        # Analyze pattern spacing in each projection to find hex tile size
        pattern_spacings = {}
        
        for direction, projection in projections.items():
            spacing = self._find_pattern_spacing(projection)
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
        
        # Look for peaks and valleys in the projection to find pattern spacing
        from scipy.signal import find_peaks
        
        # Smooth the projection slightly
        from scipy.ndimage import gaussian_filter1d
        smoothed = gaussian_filter1d(projection.astype(float), sigma=2)
        
        # Find peaks (high points in projection)
        peaks, _ = find_peaks(smoothed, height=np.max(smoothed) * 0.3, distance=20)
        
        if len(peaks) < 2:
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
        
        # Use the detected spacings directly
        spacing_x = avg_horizontal_spacing
        spacing_y = avg_vertical_spacing
        
        cols = int(map_width / spacing_x) + 1
        rows = int(map_height / spacing_y) + 1
        
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
        """Save debug visualization of 4-directional projections"""
        fig_height = 800
        fig_width = 800
        
        # Create a combined visualization
        combined_img = np.zeros((fig_height, fig_width, 3), dtype=np.uint8)
        
        # Normalize and draw each projection
        for i, (direction, projection) in enumerate(projections.items()):
            if len(projection) == 0:
                continue
                
            max_val = np.max(projection) if np.max(projection) > 0 else 1
            normalized = (projection / max_val * 255).astype(np.uint8)
            
            # Position each projection in quarters
            if direction == 'view_from_top':
                # Top quarter
                y_start, y_end = 0, fig_height // 4
                x_start, x_end = 0, fig_width
                if len(normalized) > 0:
                    line_img = np.zeros((y_end - y_start, x_end - x_start), dtype=np.uint8)
                    x_coords = np.linspace(0, x_end - x_start - 1, len(normalized), dtype=int)
                    for j, val in enumerate(normalized):
                        if j < len(x_coords):
                            line_height = int((val / 255) * (y_end - y_start))
                            line_img[-line_height:, x_coords[j]] = 255
                    combined_img[y_start:y_end, x_start:x_end, 1] = line_img  # Green
            
            elif direction == 'view_from_bottom':
                # Second quarter
                y_start, y_end = fig_height // 4, fig_height // 2
                x_start, x_end = 0, fig_width
                if len(normalized) > 0:
                    line_img = np.zeros((y_end - y_start, x_end - x_start), dtype=np.uint8)
                    x_coords = np.linspace(0, x_end - x_start - 1, len(normalized), dtype=int)
                    for j, val in enumerate(normalized):
                        if j < len(x_coords):
                            line_height = int((val / 255) * (y_end - y_start))
                            line_img[:line_height, x_coords[j]] = 255
                    combined_img[y_start:y_end, x_start:x_end, 0] = line_img  # Blue
            
            elif direction == 'view_from_left':
                # Third quarter
                y_start, y_end = fig_height // 2, 3 * fig_height // 4
                x_start, x_end = 0, fig_width
                if len(normalized) > 0:
                    line_img = np.zeros((y_end - y_start, x_end - x_start), dtype=np.uint8)
                    y_coords = np.linspace(0, y_end - y_start - 1, len(normalized), dtype=int)
                    for j, val in enumerate(normalized):
                        if j < len(y_coords):
                            line_width = int((val / 255) * (x_end - x_start))
                            line_img[y_coords[j], :line_width] = 255
                    combined_img[y_start:y_end, x_start:x_end, 2] = line_img  # Red
            
            elif direction == 'view_from_right':
                # Fourth quarter
                y_start, y_end = 3 * fig_height // 4, fig_height
                x_start, x_end = 0, fig_width
                if len(normalized) > 0:
                    line_img = np.zeros((y_end - y_start, x_end - x_start), dtype=np.uint8)
                    y_coords = np.linspace(0, y_end - y_start - 1, len(normalized), dtype=int)
                    for j, val in enumerate(normalized):
                        if j < len(y_coords):
                            line_width = int((val / 255) * (x_end - x_start))
                            line_img[y_coords[j], -line_width:] = 255
                    combined_img[y_start:y_end, x_start:x_end, 1] = line_img  # Green
        
        # Add labels
        cv2.putText(combined_img, "View from Top (Green)", (10, 30), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        cv2.putText(combined_img, "View from Bottom (Blue)", (10, fig_height // 4 + 30), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        cv2.putText(combined_img, "View from Left (Red)", (10, fig_height // 2 + 30), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        cv2.putText(combined_img, "View from Right (Green)", (10, 3 * fig_height // 4 + 30), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 255), 2)
        
        cv2.imwrite(str(self.debug_dir / "4dir_projections.png"), combined_img)
    
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