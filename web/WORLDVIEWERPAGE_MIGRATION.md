# WorldViewerPage Migration to PhaserSceneView

## Migration Completed ✅

WorldViewerPage has been successfully migrated to use the new `PhaserSceneView` template component.

## Changes Made

### 1. Template Changes (`/web/templates/WorldViewerPage.html`)

**Before:**
```html
<div class="bg-white ... h-full" data-component="world-viewer">
    <div class="p-4 border-b ...">
        <h2>World Preview</h2>
    </div>
    <div id="phaser-viewer-container" class="w-full h-full" style="max-height:800px">
        <!-- Phaser canvas -->
    </div>
</div>
```

**Issues:**
- `style="max-height:800px"` caused rigid sizing
- No `min-height: 0` constraint
- Nested structure could cause sizing issues

**After:**
```html
<div class="bg-white ... h-full flex flex-col" data-component="world-viewer">
    <!-- Header -->
    <div class="p-4 border-b ... flex-shrink-0">
        <h2>World Preview</h2>
    </div>
    <!-- Phaser Scene using PhaserSceneView -->
    <div class="flex-1" style="min-height: 0;">
        {{ template "PhaserSceneView" ( dict
          "SceneId" "world-viewer-scene"
          "CenterClass" "bg-gray-900" )
        }}
    </div>
</div>
```

**Benefits:**
- ✅ Removed `max-height:800px` constraint
- ✅ Added `flex flex-col` for proper layout
- ✅ Header is `flex-shrink-0` (fixed size)
- ✅ Scene container has `flex-1` + `min-height: 0` (takes remaining space, doesn't push parent)
- ✅ Uses PhaserSceneView template for consistency

### 2. TypeScript Changes (`/web/src/WorldViewerPage.ts`)

**Before:**
```typescript
const phaserContainer = this.ensureElement('#phaser-viewer-container', 'phaser-viewer-container');
this.worldScene = new PhaserWorldScene(phaserContainer, this.eventBus, true);
```

**After:**
```typescript
// Uses PhaserSceneView template with SceneId: "world-viewer-scene"
const phaserContainer = this.ensureElement('#world-viewer-scene', 'world-viewer-scene');
this.worldScene = new PhaserWorldScene(phaserContainer, this.eventBus, true);
```

**Changes:**
- Container ID changed from `phaser-viewer-container` → `world-viewer-scene`
- Added comment explaining the template integration
- No other changes needed (same initialization pattern)

### 3. Template Include (`/web/templates/WorldViewerPage.html`)

Added at the top:
```html
{{# include "components/PhaserSceneView.html" #}}
```

This includes the PhaserSceneView template definition so it can be used.

## How It Works

### Sizing Flow (One-Way)

```
1. Card header (fixed height) ✓
2. Scene container (flex-1) takes remaining space ✓
3. PhaserSceneView wrapper (100% of parent) ✓
4. Scene div (100% of wrapper) ✓
5. Phaser canvas (constrained by scene div) ✓
6. Canvas NEVER influences parent size ❌ (blocked by min-height: 0)
```

### BorderLayout Structure

```
┌─────────────────────────────────────┐
│ Card Header: "World Preview"       │
├─────────────────────────────────────┤
│                                     │
│   PhaserSceneView                   │
│   ┌─────────────────────────────┐   │
│   │   CENTER (Scene Container)  │   │
│   │   - No North/South/East/West│   │
│   │   - Just the Phaser scene   │   │
│   │   - Takes 100% of space     │   │
│   │   - Never grows parent      │   │
│   └─────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

## Testing Checklist

To verify the migration:

- [ ] Parent container resizes → Canvas resizes smoothly
- [ ] No more `max-height:800px` constraint
- [ ] Canvas never causes parent to grow
- [ ] No unexpected scrollbars in scene area
- [ ] Scene remains visible at all viewport sizes
- [ ] TypeScript initializes correctly with new container ID
- [ ] World loads and displays properly
- [ ] Desktop and mobile layouts both work

## Key Differences from Before

### Removed
- ❌ `max-height:800px` - was causing rigid sizing
- ❌ Direct container ID `phaser-viewer-container`
- ❌ Simple nested div structure

### Added
- ✅ `flex flex-col` layout on card
- ✅ `flex-shrink-0` on header (fixed size)
- ✅ `flex-1` + `min-height: 0` on scene container (takes remaining space)
- ✅ PhaserSceneView template integration
- ✅ New container ID: `world-viewer-scene`

## Benefits

1. **No Circular Sizing** - Canvas can't push parent to grow
2. **Responsive** - Scene adapts to parent container naturally
3. **Consistent** - Uses same pattern as other Phaser pages
4. **Maintainable** - Future sizing fixes apply automatically
5. **Flexible** - Easy to add North/South/East/West regions later if needed

## Future Enhancements

If needed, we can easily add regions:

```html
{{/* Add optional controls bar */}}
{{ define "WorldViewerControls" }}
  <div class="p-2 bg-white dark:bg-gray-800 border-b">
    <button>Show Grid</button>
    <button>Show Coordinates</button>
  </div>
{{ end }}

{{ template "PhaserSceneView" ( dict
  "NorthContent" (template "WorldViewerControls" .)
  "SceneId" "world-viewer-scene"
  "CenterClass" "bg-gray-900" )
}}
```

## Related Files

- `/web/templates/WorldViewerPage.html` - Template with PhaserSceneView
- `/web/src/WorldViewerPage.ts` - TypeScript using new container ID
- `/web/templates/components/PhaserSceneView.html` - Reusable component
- `/web/templates/components/PhaserSceneView_README.md` - Component documentation

## Next Steps

1. ✅ **WorldViewerPage** - COMPLETE
2. ⏭️ **WorldEditorPage** - Next target (uses PhaserEditorComponent)
3. ⏭️ **GameViewerPages** - After editor page
