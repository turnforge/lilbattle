# PhaserSceneView Component

A reusable BorderLayout template for Phaser scenes that solves the circular sizing problem and provides a flexible, maintainable pattern for all Phaser-based pages.

## Problem Solved

**Circular Sizing Issue**: When Phaser canvas is placed in a flex/grid container without proper constraints, it can cause recursive growth:
```
Parent resizes → Canvas resizes → Canvas pushes parent → Parent grows → Canvas grows → Loop!
```

**Solution**: This component uses `min-height: 0` and `min-width: 0` on flex children to break the circular dependency, ensuring:
- Parent size changes flow down to canvas ✅
- Canvas size NEVER influences parent ❌

## Architecture

### BorderLayout Pattern

```
┌─────────────────────────────────────────┐
│ NORTH (Optional Toolbar/Header)         │
├─────┬───────────────────────────┬───────┤
│     │                           │       │
│ W   │   CENTER (Phaser Scene)   │   E   │
│ E   │   - Takes remaining space │   A   │
│ S   │   - Never grows parent    │   S   │
│ T   │   - Constrained size      │   T   │
│     │                           │       │
├─────┴───────────────────────────┴───────┤
│ SOUTH (Optional Footer/Status Bar)      │
└─────────────────────────────────────────┘
```

### Sizing Flow (One-Way)

```
1. North/South content → Height determined by content
2. East/West content → Width determined by content
3. Remaining space → Center (scene container) calculated
4. Center size → Phaser canvas size (constrained)
5. Canvas NEVER pushes back ← Blocked by min-height: 0
```

## Usage

### Pattern: Block Inheritance (Recommended)

Use Go template block inheritance to extend PhaserSceneView regions. This is the cleanest approach:

**Step 1: Include base component and content templates**
```html
{{# include "components/PhaserSceneView.html" #}}
{{# include "panels/MyToolbar.html" #}}
```

**Step 2: Override blocks for regions you want to customize**
```html
{{ define "PhaserSceneView_North" }}
<div id="phaser-scene-view-north" class="flex-shrink-0" style="flex-shrink: 0">
  {{ template "MyToolbar" }}
</div>
{{ end }}
```

**Step 3: Call PhaserSceneView with parameters**
```html
{{ define "MyPanel" }}
  {{ template "PhaserSceneView" (dict
    "SceneId" "my-scene"
    "CenterClass" "bg-gray-900"
    "FlexMode" "fill" )
  }}
{{ end }}
```

### Example: WorldEditorPage Pattern

This is the complete pattern used in WorldEditorPage:

```html
<!-- panels/PhaserPanel.html -->
{{# include "components/PhaserSceneView.html" #}}
{{# include "panels/WorldEditorToolbar.html" #}}

{{/* Override North region with toolbar */}}
{{ define "PhaserSceneView_North" }}
<div id="phaser-scene-view-north" class="flex-shrink-0" style="flex-shrink: 0">
  {{ template "WorldEditorToolbar" }}
</div>
{{ end }}

{{/* Wrap in named template for reuse */}}
{{ define "PhaserPanel" }}
  {{ template "PhaserSceneView" (dict
    "SceneId" "phaser-container"
    "CenterClass" "bg-gray-100 dark:bg-gray-900 p-2"
    "FlexMode" "fixed" )
  }}
{{ end }}
```

Then from WorldEditorPage:
```html
{{# include "panels/PhaserPanel.html" #}}

<div id="canvas-panel-template" style="height: 100%;">
  {{ template "PhaserPanel" . }}
</div>
```

### Scene Only (No Regions)

For simple viewers, just call PhaserSceneView without defining any blocks:

```html
{{# include "components/PhaserSceneView.html" #}}

{{ template "PhaserSceneView" (dict
  "SceneId" "world-viewer-scene"
  "CenterClass" "bg-gray-900"
  "FlexMode" "fill" )
}}
```

No block overrides = scene only layout.

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `NorthContent` | HTML | - | Content for top region (toolbar/header) |
| `SouthContent` | HTML | - | Content for bottom region (footer/status) |
| `WestContent` | HTML | - | Content for left region (sidebar) |
| `EastContent` | HTML | - | Content for right region (sidebar) |
| `SceneId` | string | `"phaser-scene-container"` | ID for the Phaser scene container |
| `WrapperClass` | string | `""` | Additional CSS classes for wrapper |
| `CenterClass` | string | `""` | Additional CSS classes for center region |

## Standard IDs (for TypeScript)

All regions have predictable IDs for easy access in TypeScript:

- `phaser-scene-view-wrapper` - Main container
- `phaser-scene-view-north` - North region
- `phaser-scene-view-south` - South region
- `phaser-scene-view-east` - East region
- `phaser-scene-view-west` - West region
- `phaser-scene-view-center` - Center region (contains scene)
- `[SceneId]` - The actual scene container (customizable)

## TypeScript Integration

### Consistent Pattern

```typescript
class MyPage {
  private scene: PhaserWorldScene;

  async initializeScene() {
    // Use the SceneId you specified in the template
    const container = document.getElementById('my-scene-id');

    if (!container) {
      throw new Error('Scene container not found');
    }

    this.scene = new PhaserWorldScene(container, this.eventBus);
    await this.scene.performLocalInit();
    await this.scene.activate();
  }
}
```

### Accessing Regions

```typescript
// Toggle north toolbar visibility
toggleToolbar() {
  const north = document.getElementById('phaser-scene-view-north');
  if (north) {
    north.classList.toggle('hidden');
  }
}

// Hide east sidebar
hideEastSidebar() {
  const east = document.getElementById('phaser-scene-view-east');
  if (east) {
    east.style.display = 'none';
  }
}

// Get center dimensions for calculations
getCenterDimensions() {
  const center = document.getElementById('phaser-scene-view-center');
  return {
    width: center.clientWidth,
    height: center.clientHeight
  };
}
```

## Key CSS Constraints

These are the critical constraints that prevent circular sizing:

### 1. Flex Column Layout
```css
.wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
}
```

### 2. Fixed-Size Regions
```css
.north, .south, .east, .west {
  flex-shrink: 0; /* Don't shrink, natural size */
}
```

### 3. Center Takes Remaining Space
```css
.center {
  flex: 1;           /* Take remaining space */
  min-width: 0;      /* Allow shrinking below content size */
  min-height: 0;     /* CRITICAL: Breaks circular dependency */
  overflow: hidden;  /* Prevent content overflow */
}
```

### 4. Scene Container
```css
.scene-container {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}
```

## Common Patterns

### WorldViewerPage (Scene Only)
```html
{{ template "components/PhaserSceneView.html" dict
  "SceneId" "world-viewer-scene"
  "CenterClass" "bg-gray-900"
}}
```

### WorldEditorPage (Toolbar + Scene)
```html
{{ define "EditorToolbar" }}
  <!-- Toolbar content here -->
{{ end }}

{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "EditorToolbar" .)
  "SceneId" "world-editor-scene"
  "CenterClass" "bg-gray-100 dark:bg-gray-900 p-2"
}}
```

### GameViewerPage (Optional Toolbar + Scene)
```html
{{ define "GameToolbar" }}
  <div id="game-toolbar" class="...">
    <!-- Toolbar with toggle button -->
  </div>
{{ end }}

{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "GameToolbar" .)
  "SceneId" "game-viewer-scene"
  "CenterClass" "bg-gray-900"
}}

<script>
// Hide toolbar via JavaScript
document.getElementById('phaser-scene-view-north').classList.toggle('hidden');
</script>
```

## Benefits

1. **No Circular Sizing** - Canvas never influences parent size
2. **Responsive** - Adapts to parent container changes naturally
3. **No JavaScript Required** - Pure CSS solution
4. **Reusable** - One template for all Phaser pages
5. **Predictable IDs** - Consistent naming for TypeScript access
6. **Flexible** - Any combination of regions
7. **Maintainable** - Single source of truth for layout logic

## Migration Guide

### Old Pattern (PhaserPanel.html)
```html
<div id="phaser-panel" style="width: 100%; height: 100%">
  <div class="toolbar">...</div>
  <div class="flex-1 min-h-0">
    <div id="editor-canvas-container" class="w-full h-full">
      <!-- Phaser here -->
    </div>
  </div>
</div>
```

### New Pattern (Using PhaserSceneView)
```html
{{ define "Toolbar" }}
  <div class="toolbar">...</div>
{{ end }}

{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "Toolbar" .)
  "SceneId" "editor-canvas-container"
}}
```

**Changes in TypeScript**: None! The scene container ID remains the same.

## Testing Checklist

When using this component, verify:

- ✅ Parent container resizes → Canvas resizes smoothly
- ✅ Toolbar wraps on narrow screens → Center adjusts correctly
- ✅ Canvas never causes parent to grow
- ✅ No scrollbars appear unexpectedly
- ✅ Scene remains visible at all viewport sizes
- ✅ TypeScript can access all regions by standard IDs

## Troubleshooting

### Canvas Not Resizing
- Ensure parent has explicit height (not `height: auto`)
- Check that Phaser is configured with `Phaser.Scale.RESIZE` mode

### Parent Still Growing
- Verify `min-height: 0` is set on center region
- Check for `min-width`/`min-height` on canvas element
- Ensure no child has `flex-shrink: 1` when it should be `0`

### Regions Not Showing
- Check that content is actually passed to the template
- Verify conditional rendering: `{{ if .NorthContent }}`
- Ensure region content has visible elements

## Future Enhancements

Potential additions:
- Collapsible regions (with animation)
- Resizable splitters between regions
- Region persistence (save size/visibility state)
- Mobile-optimized variants
