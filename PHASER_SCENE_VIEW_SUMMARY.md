# PhaserSceneView - Summary of Changes

## Problem Statement

PhaserWorldScene (and its derived classes) were experiencing recursive growth where resizing the parent container would cause the canvas to grow, which in turn would push the parent to grow larger, creating an infinite loop.

## Root Cause

1. **Flexbox defaults**: Flex items have `min-height: auto` by default, preventing them from shrinking below content size
2. **Minimum dimensions**: Canvas had minimum dimensions (800x600) forcing parent to grow
3. **No size constraints**: Scene container lacked proper constraints to prevent upward size influence

## Solution: PhaserSceneView Component

Created a reusable BorderLayout template that solves the circular sizing problem using the **`min-height: 0` pattern**.

### Key Innovation

```css
.center-region {
  flex: 1;           /* Take remaining space */
  min-height: 0;     /* CRITICAL: Allow shrinking below content size */
  min-width: 0;      /* Same for horizontal */
  overflow: hidden;  /* Prevent content overflow */
}
```

This breaks the circular dependency by allowing the container to be smaller than its content.

## Files Created

### 1. Core Template
- **`/web/templates/components/PhaserSceneView.html`**
  - Reusable BorderLayout component
  - Supports 5 regions: North, South, East, West, Center
  - All regions optional, only Center (scene) is required
  - Standard IDs for TypeScript access

### 2. Updated Example
- **`/web/templates/panels/PhaserPanelNew.html`**
  - Shows how to use PhaserSceneView for WorldEditorPage
  - Toolbar in North, Scene in Center

### 3. Documentation
- **`/web/templates/components/PhaserSceneView_README.md`**
  - Complete usage guide
  - Parameter reference
  - Common patterns
  - Troubleshooting

- **`/web/templates/examples/PhaserSceneView_Examples.html`**
  - 4 usage examples
  - JavaScript integration patterns

- **`/web/templates/components/PhaserSceneView_INTEGRATION.md`**
  - Migration guide for all pages
  - Page-specific migration plans
  - Scene type hierarchy
  - Testing checklist

### 4. Code Fixes (Applied)
- **`/web/src/PhaserEditorComponent.ts`**
  - Removed `minWidth: 800px` and `minHeight: 600px`
  - Set container to `100%` with `overflow: hidden`

- **`/web/src/phaser/PhaserWorldScene.ts`**
  - Added `autoCenter: Phaser.Scale.NO_CENTER`
  - Applied constraining CSS to canvas: `width: 100%`, `height: 100%`, `object-fit: contain`

- **`/web/templates/panels/PhaserPanel.html`**
  - Added `flex flex-col` to parent for proper flexbox
  - Added `min-height: 0` to scene container

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
1. North/South content → Height determined ✓
2. East/West content → Width determined ✓
3. Remaining space → Center calculated ✓
4. Center size → Canvas size ✓
5. Canvas → Parent size ✗ (BLOCKED by min-height: 0)
```

## Scene Class Hierarchy

All scene types work with PhaserSceneView:

```
PhaserWorldScene (base)
├── PhaserEditorScene (for WorldEditorPage)
└── PhaserGameScene (for GameViewerPages)
```

## Pages Using Phaser Scenes

| Page | Scene Type | Status |
|------|-----------|--------|
| WorldEditorPage | PhaserEditorScene | ✅ Ready to migrate |
| WorldViewerPage | PhaserWorldScene | ✅ Ready to migrate |
| GameViewerPageDockView | PhaserGameScene | ✅ Ready to migrate |
| GameViewerPageGrid | PhaserGameScene | ⚠️ Uses CSS Grid, may not need |
| GameViewerPageMobile | PhaserGameScene | ✅ Would benefit (removes calc()) |
| StartGamePage | PhaserWorldScene | ✅ Ready to migrate |

## Benefits

1. ✅ **No Circular Sizing** - Canvas never influences parent size
2. ✅ **Scene Type Agnostic** - Works with all scene classes (PhaserWorldScene, PhaserEditorScene, PhaserGameScene)
3. ✅ **Responsive** - Adapts to parent container changes naturally
4. ✅ **No JavaScript Required** - Pure CSS solution
5. ✅ **Reusable** - One template for all Phaser pages
6. ✅ **Minimal Migration** - Container IDs stay the same, TypeScript changes are minimal
7. ✅ **Maintainable** - Single source of truth for layout logic

## Quick Start

### For New Pages

```html
{{/* Simple viewer - scene only */}}
{{ template "components/PhaserSceneView.html" dict
  "SceneId" "my-scene"
}}

{{/* Editor with toolbar */}}
{{ define "MyToolbar" }}
  <div class="p-2">Toolbar content</div>
{{ end }}

{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "MyToolbar" .)
  "SceneId" "my-scene"
}}
```

### TypeScript (Same for All Scene Types)

```typescript
const container = document.getElementById('my-scene');

// Use appropriate scene type
this.scene = new PhaserEditorScene(container, this.eventBus);
// OR new PhaserGameScene(...) OR new PhaserWorldScene(...)

await this.scene.performLocalInit();
await this.scene.activate();
```

## Migration Order

**Recommended priority:**

1. ✅ **WorldEditorPage** (High impact, uses PhaserEditorComponent)
2. ✅ **WorldViewerPage** (Simple case, good test)
3. ⚠️ **GameViewerPageDockView** (Medium complexity)
4. ⚠️ **GameViewerPageMobile** (Would remove calc() complexity)
5. ⚠️ **StartGamePage** (Low priority)
6. ⚠️ **GameViewerPageGrid** (May not need - already uses CSS Grid)

## Testing Checklist

Before marking migration complete:

- [ ] Parent container resizes → Canvas resizes smoothly
- [ ] Toolbar wraps on narrow screens → Center adjusts correctly
- [ ] Canvas never causes parent to grow
- [ ] No unexpected scrollbars
- [ ] Scene remains visible at all viewport sizes
- [ ] TypeScript can access all regions by standard IDs
- [ ] Works with correct scene type (PhaserWorldScene/PhaserEditorScene/PhaserGameScene)

## Next Steps

1. **Test current fix** - Verify WorldEditorPage no longer has recursive growth
2. **Migrate WorldEditorPage** - Replace PhaserPanel.html with PhaserPanelNew.html
3. **Migrate WorldViewerPage** - Simplest case, validate pattern works
4. **Update remaining pages** - Follow integration guide for each page
5. **Remove old panel templates** - Once all pages migrated

## Additional Resources

- **Full Usage Guide**: `/web/templates/components/PhaserSceneView_README.md`
- **Integration Guide**: `/web/templates/components/PhaserSceneView_INTEGRATION.md`
- **Examples**: `/web/templates/examples/PhaserSceneView_Examples.html`
