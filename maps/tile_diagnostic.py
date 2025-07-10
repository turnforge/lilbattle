#!/usr/bin/env python3
"""
Tile Classification Diagnostic Tool

Helps diagnose why tile classification confidence is low by providing
detailed analysis of extracted vs reference tiles.

USAGE:
python tile_diagnostic.py --tile extracted_tiles/00_08.png --references reference_tiles --debug
"""

import cv2
import numpy as np
import argparse
from pathlib import Path
from ml_tile_classifier import MLTileClassifier, EnhancedFeatureExtractor

def analyze_single_tile(tile_path: str, reference_dir: str):
    """Analyze a single tile and provide diagnostic information."""
    
    print(f"=== DIAGNOSTIC ANALYSIS: {tile_path} ===\n")
    
    # Load tile
    tile_image = cv2.imread(tile_path, cv2.IMREAD_UNCHANGED)
    if tile_image is None:
        print(f"❌ Cannot load tile: {tile_path}")
        return
    
    print(f"✓ Loaded tile: {tile_image.shape}")
    
    # Initialize classifier for preprocessing
    classifier = MLTileClassifier(reference_dir, debug_mode=True)
    
    # Preprocess tile
    processed_tile = classifier._preprocess_tile(tile_image)
    if processed_tile is None:
        print("❌ Preprocessing failed")
        return
    
    print(f"✓ Preprocessed tile: {processed_tile.shape}")
    
    # Check if empty
    is_empty = classifier._is_empty_tile(processed_tile)
    if is_empty:
        print("⚠️  DETECTED AS EMPTY TILE")
        print(f"   Mean intensity: {np.mean(processed_tile):.3f}")
        print(f"   Std intensity: {np.std(processed_tile):.3f}")
        return
    
    print("✓ Tile appears to have content")
    
    # Extract features
    extractor = EnhancedFeatureExtractor()
    features = extractor.extract_features(processed_tile)
    
    print(f"\n--- FEATURE ANALYSIS ---")
    print(f"Color features: {len(features.color_features)} dims")
    print(f"  Range: [{np.min(features.color_features):.2f}, {np.max(features.color_features):.2f}]")
    print(f"  Mean: {np.mean(features.color_features):.2f}")
    
    print(f"Texture features: {len(features.texture_features)} dims") 
    print(f"  Range: [{np.min(features.texture_features):.2f}, {np.max(features.texture_features):.2f}]")
    print(f"  Mean: {np.mean(features.texture_features):.2f}")
    
    print(f"Shape features: {len(features.shape_features)} dims")
    print(f"  Range: [{np.min(features.shape_features):.2f}, {np.max(features.shape_features):.2f}]")
    print(f"  Mean: {np.mean(features.shape_features):.2f}")
    
    # Classify and get detailed results
    result = classifier.classify_tile(tile_path)
    
    print(f"\n--- CLASSIFICATION RESULTS ---")
    print(f"Best match: {result.best_match}")
    print(f"Confidence: {result.best_confidence:.4f}")
    print(f"Method: {result.method_used}")
    print(f"Inference time: {result.inference_time_ms:.1f}ms")
    
    print(f"\nTop 10 matches:")
    for i, (ref_name, conf) in enumerate(result.matches[:10]):
        print(f"  {i+1:2d}. {ref_name:15s}: {conf:.4f}")
    
    # Analyze top reference match
    if result.matches:
        top_match = result.matches[0][0]
        ref_path = Path(reference_dir) / f"{top_match}.png"
        
        if ref_path.exists():
            print(f"\n--- REFERENCE TILE ANALYSIS ---")
            print(f"Comparing with: {ref_path}")
            
            ref_image = cv2.imread(str(ref_path), cv2.IMREAD_UNCHANGED)
            ref_processed = classifier._preprocess_tile(ref_image)
            
            if ref_processed is not None:
                ref_features = extractor.extract_features(ref_processed)
                
                # Compare features
                color_sim = np.corrcoef(features.color_features, ref_features.color_features)[0,1]
                texture_sim = np.corrcoef(features.texture_features, ref_features.texture_features)[0,1]
                shape_sim = np.corrcoef(features.shape_features, ref_features.shape_features)[0,1]
                
                print(f"Feature correlations with best match:")
                print(f"  Color similarity: {color_sim:.3f}")
                print(f"  Texture similarity: {texture_sim:.3f}") 
                print(f"  Shape similarity: {shape_sim:.3f}")
                
                # Visual comparison stats
                print(f"\nVisual comparison:")
                print(f"  Extracted tile - mean: {np.mean(processed_tile):.3f}, std: {np.std(processed_tile):.3f}")
                print(f"  Reference tile - mean: {np.mean(ref_processed):.3f}, std: {np.std(ref_processed):.3f}")
    
    # Recommendations
    print(f"\n--- RECOMMENDATIONS ---")
    if result.best_confidence < 0.3:
        print("❌ Very low confidence - possible issues:")
        print("   • Tile extraction may be misaligned")
        print("   • No similar reference tile exists")
        print("   • Tile might be corrupted or unusual")
    elif result.best_confidence < 0.6:
        print("⚠️  Low confidence - consider:")
        print("   • Adding more reference examples") 
        print("   • Checking tile extraction quality")
        print("   • Using Phase 2 (deep features)")
    else:
        print("✓ Good confidence level")

def compare_extraction_quality(extracted_dir: str, reference_dir: str):
    """Compare overall quality between extracted and reference tiles."""
    
    print(f"=== EXTRACTION QUALITY ANALYSIS ===\n")
    
    extracted_path = Path(extracted_dir)
    reference_path = Path(reference_dir)
    
    # Sample a few tiles from each
    extracted_tiles = list(extracted_path.glob("*.png"))[:5]
    reference_tiles = list(reference_path.glob("*.png"))[:5]
    
    print(f"Comparing {len(extracted_tiles)} extracted vs {len(reference_tiles)} reference tiles")
    
    classifier = MLTileClassifier(reference_dir, debug_mode=False)
    extractor = EnhancedFeatureExtractor()
    
    # Analyze extracted tiles
    extracted_stats = []
    for tile_path in extracted_tiles:
        tile_image = cv2.imread(str(tile_path), cv2.IMREAD_UNCHANGED)
        processed = classifier._preprocess_tile(tile_image)
        if processed is not None:
            extracted_stats.append({
                'mean': np.mean(processed),
                'std': np.std(processed),
                'size': processed.shape
            })
    
    # Analyze reference tiles  
    reference_stats = []
    for tile_path in reference_tiles:
        tile_image = cv2.imread(str(tile_path), cv2.IMREAD_UNCHANGED)
        processed = classifier._preprocess_tile(tile_image)
        if processed is not None:
            reference_stats.append({
                'mean': np.mean(processed),
                'std': np.std(processed), 
                'size': processed.shape
            })
    
    if extracted_stats and reference_stats:
        ext_means = [s['mean'] for s in extracted_stats]
        ref_means = [s['mean'] for s in reference_stats]
        
        print(f"\nIntensity comparison:")
        print(f"  Extracted tiles - mean: {np.mean(ext_means):.3f} ± {np.std(ext_means):.3f}")
        print(f"  Reference tiles - mean: {np.mean(ref_means):.3f} ± {np.std(ref_means):.3f}")
        
        mean_diff = abs(np.mean(ext_means) - np.mean(ref_means))
        if mean_diff > 0.2:
            print(f"  ⚠️  Large intensity difference: {mean_diff:.3f}")
            print("     This could indicate extraction/preprocessing issues")
        else:
            print("  ✓ Similar intensity ranges")

def main():
    parser = argparse.ArgumentParser(description='Diagnostic tool for tile classification')
    parser.add_argument('--tile', type=str, help='Single tile to analyze')
    parser.add_argument('--extracted', type=str, help='Directory with extracted tiles')
    parser.add_argument('--references', type=str, required=True, help='Reference tiles directory')
    parser.add_argument('--compare', action='store_true', help='Compare extraction quality')
    
    args = parser.parse_args()
    
    if args.tile:
        analyze_single_tile(args.tile, args.references)
    
    if args.compare and args.extracted:
        compare_extraction_quality(args.extracted, args.references)
    
    if not args.tile and not args.compare:
        parser.print_help()

if __name__ == "__main__":
    main()