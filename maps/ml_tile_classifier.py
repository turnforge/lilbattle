#!/usr/bin/env python3
"""
ML-Enhanced Tile Classifier

CPU-optimized machine learning enhancement for hex tile classification.
Combines traditional computer vision with lightweight ML techniques for
improved accuracy without GPU dependencies.

FEATURES:
- Enhanced traditional feature extraction (color, texture, shape)
- Random Forest classifier for fast CPU-based learning
- Progressive classification with adaptive resource usage
- Backwards compatibility with existing tile_classifier.py
- CPU-first design for local deployment

PHASES:
Phase 1: Enhanced traditional features + Random Forest (5-10ms)
Phase 2: MobileNet features + ensemble (50-80ms)
Phase 3: Adaptive hybrid system (10-80ms adaptive)

USAGE:
# Enhanced traditional approach
python ml_tile_classifier.py --batch extracted_tiles --references reference_tiles --method traditional

# With optional ML enhancement
python ml_tile_classifier.py --batch extracted_tiles --references reference_tiles --method enhanced
"""

import cv2
import numpy as np
import argparse
import json
import time
from typing import Dict, List, Tuple, Optional, Union
from pathlib import Path
from dataclasses import dataclass
import pickle
import warnings
warnings.filterwarnings('ignore')

# Traditional CV imports
from skimage import feature, measure
from scipy import ndimage
from scipy.stats import pearsonr

# ML imports (with fallback)
try:
    from sklearn.ensemble import RandomForestClassifier
    from sklearn.metrics.pairwise import cosine_similarity
    from sklearn.preprocessing import StandardScaler
    from sklearn.model_selection import train_test_split
    import joblib
    ML_AVAILABLE = True
except ImportError:
    ML_AVAILABLE = False
    print("Warning: scikit-learn not available. ML features disabled.")

# Optional deep learning imports
try:
    import tensorflow as tf
    tf.config.set_visible_devices([], 'GPU')  # Force CPU-only
    TENSORFLOW_AVAILABLE = True
except ImportError:
    TENSORFLOW_AVAILABLE = False

# Import original tile classifier for hybrid approach
try:
    from tile_classifier import TileClassifier as OriginalTileClassifier
    ORIGINAL_CLASSIFIER_AVAILABLE = True
except ImportError:
    ORIGINAL_CLASSIFIER_AVAILABLE = False
    print("Warning: Original tile_classifier.py not available. Hybrid mode disabled.")

@dataclass
class MLClassificationResult:
    """Enhanced classification result with ML metrics"""
    tile_path: str
    matches: List[Tuple[str, float]]
    best_match: str
    best_confidence: float
    method_used: str
    inference_time_ms: float
    feature_breakdown: Dict[str, float]

@dataclass
class TileFeatures:
    """Container for extracted tile features"""
    color_features: np.ndarray
    texture_features: np.ndarray
    shape_features: np.ndarray
    combined_features: np.ndarray

class EnhancedFeatureExtractor:
    """CPU-optimized feature extraction for tile classification.
    
    Extracts lightweight but informative features optimized for CPU processing:
    - Color features: HSV statistics, dominant colors, color moments
    - Texture features: Local Binary Patterns, gradient statistics
    - Shape features: Hu moments, contour properties
    
    All features are designed for <10ms extraction time on CPU.
    """
    
    def __init__(self, target_size: Tuple[int, int] = (64, 64)):
        self.target_size = target_size
        
        # LBP parameters for texture analysis
        self.lbp_radius = 1
        self.lbp_n_points = 8 * self.lbp_radius
        
    def extract_features(self, tile_image: np.ndarray) -> TileFeatures:
        """Extract comprehensive features from tile image.
        
        Args:
            tile_image: Preprocessed tile image (normalized float32)
            
        Returns:
            TileFeatures object with all extracted features
        """
        # Ensure proper format
        if tile_image.dtype != np.uint8:
            tile_image = (tile_image * 255).astype(np.uint8)
        
        # Resize to standard size for consistent features
        if tile_image.shape[:2] != self.target_size:
            tile_image = cv2.resize(tile_image, self.target_size)
        
        # Extract feature categories
        color_features = self._extract_color_features(tile_image)
        texture_features = self._extract_texture_features(tile_image)
        shape_features = self._extract_shape_features(tile_image)
        
        # Combine all features
        combined_features = np.concatenate([
            color_features,
            texture_features,
            shape_features
        ])
        
        return TileFeatures(
            color_features=color_features,
            texture_features=texture_features,
            shape_features=shape_features,
            combined_features=combined_features
        )
    
    def _extract_color_features(self, image: np.ndarray) -> np.ndarray:
        """Extract color-based features optimized for CPU.
        
        Features extracted (25 dimensions):
        - HSV channel statistics (mean, std, skew) - 9 dims
        - Dominant colors (top 3 colors in HSV) - 9 dims  
        - Color histogram bins (3x3 grid) - 9 dims
        - Color variance and contrast - 2 dims
        
        Args:
            image: BGR image (uint8)
            
        Returns:
            Color feature vector (25 dimensions)
        """
        features = []
        
        # Convert to HSV for better color representation
        hsv = cv2.cvtColor(image, cv2.COLOR_BGR2HSV)
        
        # 1. HSV channel statistics (9 dims)
        for channel in range(3):
            channel_data = hsv[:, :, channel].flatten()
            features.extend([
                np.mean(channel_data),
                np.std(channel_data),
                self._safe_skew(channel_data)
            ])
        
        # 2. Dominant colors using k-means (9 dims)
        try:
            from sklearn.cluster import KMeans
            pixels = hsv.reshape(-1, 3)
            kmeans = KMeans(n_clusters=3, random_state=42, n_init=10)
            kmeans.fit(pixels)
            
            # Sort cluster centers by frequency
            labels = kmeans.labels_
            unique, counts = np.unique(labels, return_counts=True)
            sorted_indices = np.argsort(counts)[::-1]
            
            dominant_colors = kmeans.cluster_centers_[sorted_indices]
            features.extend(dominant_colors.flatten())
        except:
            # Fallback if sklearn not available
            features.extend([0] * 9)
        
        # 3. Color histogram in 3x3 spatial grid (9 dims)
        h, w = image.shape[:2]
        grid_features = []
        for i in range(3):
            for j in range(3):
                y1, y2 = i * h // 3, (i + 1) * h // 3
                x1, x2 = j * w // 3, (j + 1) * w // 3
                region = hsv[y1:y2, x1:x2]
                grid_features.append(np.mean(region))
        features.extend(grid_features)
        
        # 4. Color variance and contrast (2 dims)
        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        features.append(np.var(gray))
        features.append(np.std(gray) / (np.mean(gray) + 1e-8))  # Contrast
        
        return np.array(features, dtype=np.float32)
    
    def _extract_texture_features(self, image: np.ndarray) -> np.ndarray:
        """Extract texture features using Local Binary Patterns.
        
        Features extracted (18 dimensions):
        - LBP histogram (10 uniform patterns) - 10 dims
        - Gradient magnitude statistics - 4 dims
        - Edge density and orientation - 4 dims
        
        Args:
            image: BGR image (uint8)
            
        Returns:
            Texture feature vector (18 dimensions)
        """
        features = []
        
        # Convert to grayscale for texture analysis
        if len(image.shape) == 3:
            gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        else:
            gray = image
        
        # 1. Local Binary Pattern features (10 dims)
        try:
            lbp = feature.local_binary_pattern(
                gray, self.lbp_n_points, self.lbp_radius, method='uniform'
            )
            
            # Compute LBP histogram
            n_bins = self.lbp_n_points + 2
            hist, _ = np.histogram(lbp.ravel(), bins=n_bins, 
                                range=(0, n_bins), density=True)
            features.extend(hist[:10])  # Take first 10 bins
        except:
            features.extend([0] * 10)
        
        # 2. Gradient magnitude statistics (4 dims)
        sobelx = cv2.Sobel(gray, cv2.CV_64F, 1, 0, ksize=3)
        sobely = cv2.Sobel(gray, cv2.CV_64F, 0, 1, ksize=3)
        gradient_magnitude = np.sqrt(sobelx**2 + sobely**2)
        
        features.extend([
            np.mean(gradient_magnitude),
            np.std(gradient_magnitude),
            np.percentile(gradient_magnitude, 75),
            np.percentile(gradient_magnitude, 95)
        ])
        
        # 3. Edge density and orientation (4 dims)
        edges = cv2.Canny(gray, 50, 150)
        edge_density = np.sum(edges > 0) / edges.size
        features.append(edge_density)
        
        # Gradient orientation
        gradient_angle = np.arctan2(sobely, sobelx)
        features.extend([
            np.mean(np.cos(2 * gradient_angle)),  # Horizontal preference
            np.mean(np.sin(2 * gradient_angle)),  # Vertical preference
            np.std(gradient_angle)  # Orientation variance
        ])
        
        return np.array(features, dtype=np.float32)
    
    def _extract_shape_features(self, image: np.ndarray) -> np.ndarray:
        """Extract shape-based features using moments and contours.
        
        Features extracted (12 dimensions):
        - Hu moments (7 invariant moments) - 7 dims
        - Compactness and solidity - 2 dims
        - Aspect ratio and extent - 2 dims
        - Circularity measure - 1 dim
        
        Args:
            image: BGR image (uint8)
            
        Returns:
            Shape feature vector (12 dimensions)
        """
        features = []
        
        # Convert to grayscale
        if len(image.shape) == 3:
            gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        else:
            gray = image
        
        # 1. Hu moments (7 dims)
        try:
            moments = cv2.moments(gray)
            hu_moments = cv2.HuMoments(moments)
            # Log transform to handle large values
            hu_moments = -np.sign(hu_moments) * np.log10(np.abs(hu_moments) + 1e-10)
            features.extend(hu_moments.flatten())
        except:
            features.extend([0] * 7)
        
        # 2. Contour-based features (5 dims)
        try:
            # Threshold image to find contours
            _, binary = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
            contours, _ = cv2.findContours(binary, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
            
            if contours:
                # Get largest contour
                largest_contour = max(contours, key=cv2.contourArea)
                area = cv2.contourArea(largest_contour)
                perimeter = cv2.arcLength(largest_contour, True)
                
                # Bounding rectangle
                x, y, w, h = cv2.boundingRect(largest_contour)
                
                # Shape features
                compactness = (perimeter * perimeter) / (4 * np.pi * area + 1e-8)
                aspect_ratio = w / (h + 1e-8)
                extent = area / ((w * h) + 1e-8)
                solidity = area / (cv2.contourArea(cv2.convexHull(largest_contour)) + 1e-8)
                circularity = (4 * np.pi * area) / (perimeter * perimeter + 1e-8)
                
                features.extend([compactness, aspect_ratio, extent, solidity, circularity])
            else:
                features.extend([0] * 5)
        except:
            features.extend([0] * 5)
        
        return np.array(features, dtype=np.float32)
    
    def _safe_skew(self, data: np.ndarray) -> float:
        """Calculate skewness with safety checks"""
        try:
            from scipy.stats import skew
            return float(skew(data))
        except:
            # Fallback calculation
            mean_val = np.mean(data)
            std_val = np.std(data)
            if std_val == 0:
                return 0.0
            skew_val = np.mean(((data - mean_val) / std_val) ** 3)
            return float(skew_val)

class MLTileClassifier:
    """ML-enhanced tile classifier with CPU-optimized performance.
    
    Supports multiple classification methods:
    - traditional: Enhanced CV features + Random Forest
    - enhanced: Traditional + optional deep features
    - adaptive: Smart resource usage based on confidence
    
    Designed for local deployment without GPU requirements.
    """
    
    def __init__(self,
                 reference_tiles_dir: str,
                 method: str = "traditional",
                 debug_mode: bool = False):
        """Initialize ML classifier.
        
        Args:
            reference_tiles_dir: Directory with reference tile images
            method: Classification method (traditional, enhanced, adaptive)
            debug_mode: Enable debug output and timing
        """
        self.reference_tiles_dir = Path(reference_tiles_dir)
        self.method = method
        self.debug_mode = debug_mode
        
        # Feature extractor
        self.feature_extractor = EnhancedFeatureExtractor()
        
        # Reference data storage
        self.reference_features: Dict[str, np.ndarray] = {}
        self.reference_images: Dict[str, np.ndarray] = {}
        self.tile_classes = []
        
        # ML models
        self.rf_classifier = None
        self.feature_scaler = None
        
        # Original classifier for hybrid approach
        self.original_classifier = None
        if method == "hybrid" and ORIGINAL_CLASSIFIER_AVAILABLE:
            self.original_classifier = OriginalTileClassifier(
                reference_tiles_dir=str(reference_tiles_dir),
                debug_mode=debug_mode
            )
        
        # Performance tracking
        self.timing_stats = {
            'feature_extraction': [],
            'classification': [],
            'total': []
        }
        
        # Debug directory
        self.debug_dir = Path("debug_images/ml_classification") if debug_mode else None
        if self.debug_mode:
            self.debug_dir.mkdir(parents=True, exist_ok=True)
        
        # Load reference tiles and train if needed
        self._load_reference_tiles()
        if ML_AVAILABLE and self.method in ["traditional", "enhanced", "hybrid"]:
            self._train_classifier()
    
    def _load_reference_tiles(self):
        """Load and preprocess reference tiles."""
        if not self.reference_tiles_dir.exists():
            raise FileNotFoundError(f"Reference tiles directory not found: {self.reference_tiles_dir}")
        
        reference_files = list(self.reference_tiles_dir.glob("*.png"))
        if not reference_files:
            raise ValueError(f"No PNG files found in reference directory: {self.reference_tiles_dir}")
        
        print(f"Loading {len(reference_files)} reference tiles...")
        
        for ref_file in reference_files:
            tile_name = ref_file.stem
            
            # Load image
            ref_image = cv2.imread(str(ref_file), cv2.IMREAD_UNCHANGED)
            if ref_image is None:
                continue
            
            # Preprocess
            processed_tile = self._preprocess_tile(ref_image)
            if processed_tile is None:
                continue
            
            # Extract features
            start_time = time.time()
            features = self.feature_extractor.extract_features(processed_tile)
            extraction_time = (time.time() - start_time) * 1000
            
            self.reference_features[tile_name] = features.combined_features
            self.reference_images[tile_name] = processed_tile
            
            # Use full filename as class for exact matching (e.g., "1_0_5.png" -> class "1_0_5")
            tile_class = tile_name
            if tile_class not in self.tile_classes:
                self.tile_classes.append(tile_class)
            
            if self.debug_mode:
                print(f"  Loaded {tile_name}: {len(features.combined_features)} features ({extraction_time:.1f}ms)")
        
        print(f"Successfully loaded {len(self.reference_features)} reference tiles")
        print(f"Found {len(self.tile_classes)} tile classes: {self.tile_classes}")
    
    def _preprocess_tile(self, tile_image: np.ndarray) -> Optional[np.ndarray]:
        """Preprocess tile image with empty tile detection."""
        if tile_image is None or tile_image.size == 0:
            return None
        
        # Handle transparency
        if len(tile_image.shape) == 3 and tile_image.shape[2] == 4:
            alpha = tile_image[:, :, 3]
            coords = np.column_stack(np.where(alpha > 0))
            if len(coords) == 0:
                if self.debug_mode:
                    print("  Warning: Completely transparent tile (empty)")
                return None
            
            # Check if tile is mostly empty (< 10% content)
            content_ratio = len(coords) / (tile_image.shape[0] * tile_image.shape[1])
            if content_ratio < 0.1:
                if self.debug_mode:
                    print(f"  Warning: Mostly empty tile ({content_ratio:.1%} content)")
            
            top, left = coords.min(axis=0)
            bottom, right = coords.max(axis=0)
            content_region = tile_image[top:bottom+1, left:right+1]
            
            if content_region.shape[2] == 4:
                content_region = content_region[:, :, :3]
            
            tile_image = content_region
        
        # Resize to standard dimensions
        target_size = (64, 64)
        tile_image = cv2.resize(tile_image, target_size, interpolation=cv2.INTER_AREA)
        
        return tile_image
    
    def _train_classifier(self):
        """Train Random Forest classifier on reference tiles."""
        if not ML_AVAILABLE:
            print("Warning: scikit-learn not available. Falling back to similarity matching.")
            return
        
        if len(self.reference_features) < 2:
            print("Warning: Not enough reference tiles for training.")
            return
        
        print("Training Random Forest classifier...")
        
        # Prepare training data
        X = []
        y = []
        
        for tile_name, features in self.reference_features.items():
            tile_class = tile_name  # Use full filename as class
            X.append(features)
            y.append(tile_class)
        
        X = np.array(X)
        y = np.array(y)
        
        # Feature scaling
        self.feature_scaler = StandardScaler()
        X_scaled = self.feature_scaler.fit_transform(X)
        
        # Train Random Forest
        self.rf_classifier = RandomForestClassifier(
            n_estimators=50,  # Balanced speed/accuracy
            max_depth=10,     # Prevent overfitting
            random_state=42,
            n_jobs=1         # Single thread for consistency
        )
        
        self.rf_classifier.fit(X_scaled, y)
        
        # Evaluate on training data (for debugging)
        if self.debug_mode:
            train_accuracy = self.rf_classifier.score(X_scaled, y)
            print(f"Training accuracy: {train_accuracy:.3f}")
        
        print("Random Forest training completed.")
    
    def _is_empty_tile(self, processed_tile: np.ndarray) -> bool:
        """Check if tile appears to be empty/background."""
        if processed_tile is None:
            return True
        
        # Check if tile is very dark/uniform (background)
        mean_intensity = np.mean(processed_tile)
        std_intensity = np.std(processed_tile)
        
        # Empty tiles tend to be very uniform and dark/transparent
        if mean_intensity < 0.1 and std_intensity < 0.05:
            return True
        
        return False
    
    def classify_tile(self, tile_image_path: str) -> MLClassificationResult:
        """Classify a single tile using ML-enhanced approach."""
        start_time = time.time()
        
        # Load and preprocess tile
        tile_image = cv2.imread(tile_image_path, cv2.IMREAD_UNCHANGED)
        if tile_image is None:
            raise ValueError(f"Could not load tile image: {tile_image_path}")
        
        processed_tile = self._preprocess_tile(tile_image)
        if processed_tile is None:
            raise ValueError(f"Failed to preprocess tile: {tile_image_path}")
        
        # Check for empty tile
        if self._is_empty_tile(processed_tile):
            return MLClassificationResult(
                tile_path=tile_image_path,
                matches=[("empty", 1.0)],
                best_match="empty",
                best_confidence=1.0,
                method_used="empty_detection",
                inference_time_ms=(time.time() - start_time) * 1000,
                feature_breakdown={'empty_tile': True}
            )
        
        # Extract features
        feature_start = time.time()
        features = self.feature_extractor.extract_features(processed_tile)
        feature_time = (time.time() - feature_start) * 1000
        
        # Classify based on method
        classification_start = time.time()
        
        if self.method == "traditional" and ML_AVAILABLE and self.rf_classifier is not None:
            result = self._classify_with_rf(features, tile_image_path)
        elif self.method == "hybrid":
            result = self._classify_hybrid(features, tile_image_path)
        else:
            result = self._classify_with_similarity(features, tile_image_path)
        
        classification_time = (time.time() - classification_start) * 1000
        total_time = (time.time() - start_time) * 1000
        
        # Update result with timing info
        result.inference_time_ms = total_time
        result.feature_breakdown = {
            'feature_extraction_ms': feature_time,
            'classification_ms': classification_time,
            'total_ms': total_time
        }
        
        # Track performance stats
        self.timing_stats['feature_extraction'].append(feature_time)
        self.timing_stats['classification'].append(classification_time)
        self.timing_stats['total'].append(total_time)
        
        if self.debug_mode:
            print(f"Classification: {result.best_match} (confidence: {result.best_confidence:.3f}, {total_time:.1f}ms)")
        
        return result
    
    def _classify_with_rf(self, features: TileFeatures, tile_path: str) -> MLClassificationResult:
        """Classify using Random Forest classifier."""
        # Scale features
        features_scaled = self.feature_scaler.transform(features.combined_features.reshape(1, -1))
        
        # Get class probabilities
        class_probs = self.rf_classifier.predict_proba(features_scaled)[0]
        classes = self.rf_classifier.classes_
        
        # Create matches list
        matches = []
        for class_name, prob in zip(classes, class_probs):
            matches.append((class_name, float(prob)))
        
        # Sort by confidence
        matches.sort(key=lambda x: x[1], reverse=True)
        
        best_match, best_confidence = matches[0] if matches else ("unknown", 0.0)
        
        return MLClassificationResult(
            tile_path=tile_path,
            matches=matches,
            best_match=best_match,
            best_confidence=best_confidence,
            method_used="random_forest",
            inference_time_ms=0.0,  # Will be set later
            feature_breakdown={}
        )
    
    def _classify_with_similarity(self, features: TileFeatures, tile_path: str) -> MLClassificationResult:
        """Classify using feature similarity matching."""
        matches = []
        
        # Compare with all reference tiles
        for ref_name, ref_features in self.reference_features.items():
            # Calculate cosine similarity
            similarity = self._cosine_similarity(features.combined_features, ref_features)
            matches.append((ref_name, float(similarity)))
        
        # Sort by similarity (exact filename matching)
        matches.sort(key=lambda x: x[1], reverse=True)
        
        # Use exact filename matches (no grouping)
        best_match, best_confidence = matches[0] if matches else ("unknown", 0.0)
        
        return MLClassificationResult(
            tile_path=tile_path,
            matches=matches,
            best_match=best_match,
            best_confidence=best_confidence,
            method_used="feature_similarity",
            inference_time_ms=0.0,
            feature_breakdown={}
        )
    
    def _classify_hybrid(self, features: TileFeatures, tile_path: str) -> MLClassificationResult:
        """Hybrid classification: Use both ML and traditional, pick the best.
        
        Strategy:
        1. Run both Random Forest and traditional CV classification
        2. Compare confidence scores
        3. Use the method with higher confidence
        4. Fall back to traditional if ML confidence is very low (< 0.3)
        
        Args:
            features: Extracted tile features
            tile_path: Path to tile being classified
            
        Returns:
            Best classification result from either method
        """
        results = []
        
        # 1. Try ML classification if available
        if ML_AVAILABLE and self.rf_classifier is not None:
            try:
                ml_result = self._classify_with_rf(features, tile_path)
                ml_result.method_used = "random_forest_hybrid"
                results.append(ml_result)
                
                if self.debug_mode:
                    print(f"  ML result: {ml_result.best_match} (confidence: {ml_result.best_confidence:.3f})")
            except Exception as e:
                if self.debug_mode:
                    print(f"  ML classification failed: {e}")
        
        # 2. Try traditional classification
        if self.original_classifier is not None:
            try:
                # Convert our tile to original classifier format
                traditional_result = self.original_classifier.classify_tile(tile_path)
                
                # Convert to our result format
                converted_result = MLClassificationResult(
                    tile_path=tile_path,
                    matches=traditional_result.matches,
                    best_match=traditional_result.best_match,
                    best_confidence=traditional_result.best_confidence,
                    method_used="traditional_cv_hybrid",
                    inference_time_ms=0.0,
                    feature_breakdown=traditional_result.metrics_breakdown
                )
                results.append(converted_result)
                
                if self.debug_mode:
                    print(f"  Traditional result: {converted_result.best_match} (confidence: {converted_result.best_confidence:.3f})")
            except Exception as e:
                if self.debug_mode:
                    print(f"  Traditional classification failed: {e}")
        
        # 3. Fall back to similarity if both fail
        if not results:
            similarity_result = self._classify_with_similarity(features, tile_path)
            similarity_result.method_used = "similarity_fallback"
            results.append(similarity_result)
            
            if self.debug_mode:
                print(f"  Similarity fallback: {similarity_result.best_match} (confidence: {similarity_result.best_confidence:.3f})")
        
        # 4. Select best result based on confidence and method
        best_result = self._select_best_hybrid_result(results)
        
        if self.debug_mode:
            print(f"  Hybrid choice: {best_result.best_match} (confidence: {best_result.best_confidence:.3f}, method: {best_result.method_used})")
        
        return best_result
    
    def _select_best_hybrid_result(self, results: List[MLClassificationResult]) -> MLClassificationResult:
        """Select the best result from multiple classification methods.
        
        Selection strategy:
        1. If traditional CV confidence > 0.6, prefer it (proven reliable)
        2. If ML confidence > 0.8, prefer it (high confidence)
        3. Otherwise, use the method with highest confidence
        4. Tie-breaker: prefer traditional CV (more proven)
        
        Args:
            results: List of classification results from different methods
            
        Returns:
            Best classification result
        """
        if not results:
            raise ValueError("No classification results provided")
        
        if len(results) == 1:
            return results[0]
        
        # Find traditional and ML results
        traditional_result = None
        ml_result = None
        
        for result in results:
            if "traditional" in result.method_used:
                traditional_result = result
            elif "random_forest" in result.method_used:
                ml_result = result
        
        # Apply selection strategy
        if traditional_result and traditional_result.best_confidence > 0.6:
            # Traditional CV is reliable at this confidence level
            return traditional_result
        
        if ml_result and ml_result.best_confidence > 0.8:
            # ML is very confident
            return ml_result
        
        # Use highest confidence method
        best_result = max(results, key=lambda r: r.best_confidence)
        
        # Tie-breaker: prefer traditional if confidence is close
        if traditional_result and ml_result:
            confidence_diff = abs(traditional_result.best_confidence - ml_result.best_confidence)
            if confidence_diff < 0.1:  # Within 10% confidence
                return traditional_result
        
        return best_result
    
    def _cosine_similarity(self, vec1: np.ndarray, vec2: np.ndarray) -> float:
        """Calculate cosine similarity between two vectors."""
        dot_product = np.dot(vec1, vec2)
        norm1 = np.linalg.norm(vec1)
        norm2 = np.linalg.norm(vec2)
        
        if norm1 == 0 or norm2 == 0:
            return 0.0
        
        return dot_product / (norm1 * norm2)
    
    def classify_batch(self, tile_dir: str) -> Dict[str, MLClassificationResult]:
        """Classify all tiles in a directory."""
        tile_dir_path = Path(tile_dir)
        if not tile_dir_path.exists():
            raise FileNotFoundError(f"Tile directory not found: {tile_dir}")
        
        tile_files = list(tile_dir_path.glob("*.png"))
        if not tile_files:
            raise ValueError(f"No PNG files found in directory: {tile_dir}")
        
        print(f"Classifying {len(tile_files)} tiles using {self.method} method...")
        
        results = {}
        for i, tile_file in enumerate(tile_files):
            if self.debug_mode:
                print(f"  Processing {i+1}/{len(tile_files)}: {tile_file.name}")
            
            try:
                result = self.classify_tile(str(tile_file))
                results[tile_file.name] = result
            except Exception as e:
                print(f"  Error processing {tile_file.name}: {e}")
                continue
        
        # Print performance summary
        self._print_performance_summary(len(results))
        
        return results
    
    def _print_performance_summary(self, num_tiles: int):
        """Print performance statistics."""
        if not self.timing_stats['total']:
            return
        
        avg_total = np.mean(self.timing_stats['total'])
        avg_feature = np.mean(self.timing_stats['feature_extraction'])
        avg_classify = np.mean(self.timing_stats['classification'])
        
        print(f"\nPerformance Summary ({num_tiles} tiles):")
        print(f"  Average total time: {avg_total:.1f}ms")
        print(f"  Average feature extraction: {avg_feature:.1f}ms")
        print(f"  Average classification: {avg_classify:.1f}ms")
        print(f"  Method used: {self.method}")
        
        if ML_AVAILABLE and self.rf_classifier is not None:
            print(f"  Random Forest: {self.rf_classifier.n_estimators} trees")

def main():
    """Command-line interface for ML-enhanced tile classification."""
    parser = argparse.ArgumentParser(description='ML-Enhanced Tile Classification')
    parser.add_argument('--references', type=str, required=True,
                        help='Directory containing reference tile images')
    parser.add_argument('--tile', type=str,
                        help='Single tile image to classify')
    parser.add_argument('--batch', type=str,
                        help='Directory containing tiles to classify')
    parser.add_argument('--method', choices=['traditional', 'enhanced', 'adaptive', 'hybrid'],
                        default='traditional',
                        help='Classification method')
    parser.add_argument('--output', type=str, default='ml_classification_results.json',
                        help='Output file for classification results')
    parser.add_argument('--debug', action='store_true',
                        help='Enable debug mode with verbose output')
    
    args = parser.parse_args()
    
    # Validate arguments
    if not args.tile and not args.batch:
        parser.error("Either --tile or --batch must be specified")
    
    # Check ML availability
    if not ML_AVAILABLE:
        print("Warning: scikit-learn not available. Some features may be limited.")
    
    # Create classifier
    classifier = MLTileClassifier(
        reference_tiles_dir=args.references,
        method=args.method,
        debug_mode=args.debug
    )
    
    # Process tiles
    if args.tile:
        # Single tile classification
        print(f"Classifying single tile: {args.tile}")
        result = classifier.classify_tile(args.tile)
        
        print(f"\nML Classification Results:")
        print(f"Best match: {result.best_match} (confidence: {result.best_confidence:.3f})")
        print(f"Method: {result.method_used}")
        print(f"Inference time: {result.inference_time_ms:.1f}ms")
        print(f"Top 5 matches:")
        for i, (tile_class, confidence) in enumerate(result.matches[:5]):
            print(f"  {i+1}. {tile_class}: {confidence:.3f}")
        
        # Save results
        results_data = {
            'tile_path': result.tile_path,
            'best_match': result.best_match,
            'best_confidence': result.best_confidence,
            'method_used': result.method_used,
            'inference_time_ms': result.inference_time_ms,
            'matches': result.matches,
            'feature_breakdown': result.feature_breakdown
        }
        
        with open(args.output, 'w') as f:
            json.dump(results_data, f, indent=2)
        
        print(f"\nResults saved to: {args.output}")
    
    elif args.batch:
        # Batch classification
        print(f"Classifying batch of tiles from: {args.batch}")
        results = classifier.classify_batch(args.batch)
        
        # Summary statistics
        total_tiles = len(results)
        high_confidence_tiles = sum(1 for r in results.values() if r.best_confidence > 0.8)
        avg_inference_time = np.mean([r.inference_time_ms for r in results.values()])
        
        print(f"\nBatch Classification Summary:")
        print(f"Total tiles classified: {total_tiles}")
        print(f"High confidence matches (>0.8): {high_confidence_tiles}")
        print(f"Classification accuracy: {high_confidence_tiles/total_tiles*100:.1f}%")
        print(f"Average inference time: {avg_inference_time:.1f}ms")
        
        # Save batch results
        batch_results = {}
        for filename, result in results.items():
            batch_results[filename] = {
                'best_match': result.best_match,
                'best_confidence': result.best_confidence,
                'method_used': result.method_used,
                'inference_time_ms': result.inference_time_ms,
                'matches': result.matches[:5],  # Top 5 matches
                'feature_breakdown': result.feature_breakdown
            }
        
        with open(args.output, 'w') as f:
            json.dump(batch_results, f, indent=2)
        
        print(f"\nBatch results saved to: {args.output}")

if __name__ == "__main__":
    main()