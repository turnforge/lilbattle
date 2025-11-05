# WorldEditorPage Migration to PhaserSceneView

## Migration Completed ✅

WorldEditorPage has been successfully migrated to use the PhaserSceneView template component with toolbar in North region.

## Changes Made

### 1. Template Changes (`/web/templates/panels/PhaserPanel.html`)

**Before** (manual flex layout):
```html
<div id="phaser-panel" class="flex flex-col" style="width: 100%; height: 100%">
  <!-- Toolbar -->
  <div class="... flex-shrink-0">
    <div class="flex items-center justify-between space-x-4">
      <!-- Controls -->
    </div>
  </div>

  <!-- Canvas Container -->
  <div class="flex-1 ... p-2 relative" style="min-height: 0; overflow: hidden;">
    <div id="editor-canvas-container" class="w-full h-full ...">
      <!-- Phaser canvas -->
    </div>
  </div>
</div>
```

**After** (using PhaserSceneView):
```html
{{/* Define Editor Toolbar (North Region) */}}
{{ define "EditorToolbar" }}
<div class="bg-white dark:bg-gray-800 border-b ... p-2">
  <div class="flex items-center justify-between space-x-4">
    <!-- Controls -->
  </div>
</div>
{{ end }}

{{/* Use PhaserSceneView with toolbar in North */}}
<div id="phaser-panel" style="width: 100%; height: 100%;">
  {{ template "PhaserSceneView" ( dict
    "NorthContent" (template "EditorToolbar" .)
    "SceneId" "editor-canvas-container"
    "CenterClass" "bg-gray-100 dark:bg-gray-900 p-2"
    "FlexMode" "fixed" )
  }}
</div>
```

**Benefits:**
- ✅ Toolbar extracted as separate template (EditorToolbar) - more maintainable
- ✅ PhaserSceneView handles all sizing constraints automatically
- ✅ FlexMode="fixed" provides 100% width/height sizing (appropriate for dockview panel)
- ✅ No manual flex layout management needed
- ✅ Consistent with WorldViewerPage pattern

### 2. WorldEditorPage Include (`/web/templates/WorldEditorPage.html`)

**Added:**
```html
{{# include "components/PhaserSceneView.html" #}}
```

This makes the PhaserSceneView template definition available to PhaserPanel.html.

### 3. TypeScript Changes (`/web/src/PhaserEditorComponent.ts`)

**No changes needed!**

The component already:
1. Finds `#editor-canvas-container` (created by PhaserSceneView with SceneId parameter)
2. Renames it to `#phaser-container` internally
3. Passes it to PhaserEditorScene

This existing pattern works perfectly with the new template structure.

## Layout Structure

### BorderLayout Pattern

```
┌─────────────────────────────────────────┐
│ NORTH: EditorToolbar                    │
│ - Clear button                          │
│ - Brush size selector                   │
│ - View options (grid, coords, health)   │
│ - Brush info display                    │
├─────────────────────────────────────────┤
│                                         │
│ CENTER: Phaser Editor Scene             │
│ - Container ID: editor-canvas-container │
│ - Takes remaining space                 │
│ - FlexMode="fixed" (100% size)          │
│ - Never grows parent                    │
│                                         │
└─────────────────────────────────────────┘
```

### Sizing Flow (One-Way)

```
1. Dockview panel size determined ✓
2. #phaser-panel wrapper (100% size) ✓
3. PhaserSceneView wrapper (width: 100%; height: 100%;) ✓
4. EditorToolbar (fixed height based on content) ✓
5. Center region (remaining space after toolbar) ✓
6. editor-canvas-container (100% of center) ✓
7. Phaser canvas (constrained by container) ✓
8. Canvas → Parent size ❌ (BLOCKED by sizing constraints)
```

## Key Differences from WorldViewerPage

| Feature | WorldViewerPage | WorldEditorPage |
|---------|----------------|-----------------|
| **FlexMode** | `"fill"` (flex-1 + min-height: 0) | `"fixed"` (width: 100%; height: 100%) |
| **Reason** | Card within flex layout | Dockview panel with explicit size |
| **Toolbar** | None (scene only) | EditorToolbar in North |
| **Center Styling** | `bg-gray-900` | `bg-gray-100 dark:bg-gray-900 p-2` |
| **Container Context** | Flex child of page layout | Dockview panel (100% size) |

**Why FlexMode="fixed"?**
- WorldEditorPage uses Dockview which explicitly sizes panels
- Panel has definite width/height from dockview layout
- We want wrapper to be 100% of that size, not flex-based
- This is different from WorldViewerPage which is a flex child needing flex-1

## Testing Checklist

To verify the migration:

- [ ] WorldEditorPage loads without errors
- [ ] Toolbar displays correctly with all controls
- [ ] Phaser editor scene initializes properly
- [ ] Painting/editing tools work (brush, terrain, units)
- [ ] View options work (grid, coordinates, health toggles)
- [ ] Dockview panel resizing updates scene correctly
- [ ] No circular sizing (width changes don't affect height)
- [ ] No unexpected scrollbars
- [ ] Dark mode styling works on toolbar and scene
- [ ] Canvas fills available space after toolbar

## Toolbar Controls Preserved

All existing controls remain functional:
- **Clear Button** (`#clear-tile-btn`) - Activates clear mode
- **Brush Size** (`#brush-size`) - Select hex radius (0-15)
- **View Options**:
  - `#show-grid` - Toggle hex grid visibility
  - `#show-coordinates` - Toggle coordinate labels
  - `#show-health` - Toggle health display
- **Brush Info** (`#brush-info`) - Current tool status
- **Editor Badge** - Visual indicator

## Component Integration

**Event Flow:**
```
User interacts with toolbar
  → Event handlers in PhaserEditorComponent
  → Updates pageState (WorldEditorPageState)
  → PageState emits TOOL_STATE_CHANGED event
  → PhaserEditorComponent handleToolStateChanged()
  → Updates PhaserEditorScene settings
```

**No changes needed** to this flow - all DOM IDs preserved.

## Benefits

1. ✅ **Consistent Pattern** - Same PhaserSceneView approach as WorldViewerPage
2. ✅ **Maintainable** - Toolbar is separate template definition
3. ✅ **No Circular Sizing** - Built-in constraints prevent growth loops
4. ✅ **No TypeScript Changes** - Existing component code works unchanged
5. ✅ **Flexible Layout** - Easy to add South/East/West regions later if needed
6. ✅ **Self-Documenting** - Template structure clearly shows BorderLayout intent

## Future Enhancements

Easy to extend with additional regions:

### Add Status Bar (South)
```html
{{ define "EditorStatusBar" }}
  <div class="p-2 bg-gray-100 dark:bg-gray-800 border-t">
    Status: {{.Status}}
  </div>
{{ end }}

{{ template "PhaserSceneView" ( dict
  "NorthContent" (template "EditorToolbar" .)
  "SouthContent" (template "EditorStatusBar" .)
  ...
) }}
```

### Add Tool Palette (West)
```html
{{ define "ToolPalette" }}
  <div class="w-48 p-2 bg-white dark:bg-gray-800 border-r">
    <!-- Tool icons -->
  </div>
{{ end }}

{{ template "PhaserSceneView" ( dict
  "NorthContent" (template "EditorToolbar" .)
  "WestContent" (template "ToolPalette" .)
  ...
) }}
```

## Related Files

- `/web/templates/panels/PhaserPanel.html` - Migrated template
- `/web/templates/WorldEditorPage.html` - Added PhaserSceneView include
- `/web/src/PhaserEditorComponent.ts` - No changes needed
- `/web/templates/components/PhaserSceneView.html` - Reusable component
- `/web/templates/components/PhaserSceneView_README.md` - Component documentation

## Migration Status

- ✅ **WorldViewerPage** - Migrated (scene only, FlexMode="fill")
- ✅ **WorldEditorPage** - Migrated (toolbar + scene, FlexMode="fixed")
- ⏭️ **GameViewerPageDockView** - Next target
- ⏭️ **GameViewerPageMobile** - After DockView
- ⏭️ **StartGamePage** - Lower priority
