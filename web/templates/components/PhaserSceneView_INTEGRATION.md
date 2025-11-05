# PhaserSceneView Integration Guide

This guide shows how to integrate `PhaserSceneView.html` into all pages using Phaser scenes.

## Scene Class Hierarchy

```
PhaserWorldScene (base class)
├── PhaserEditorScene (extends PhaserWorldScene)
└── PhaserGameScene (extends PhaserWorldScene)
```

All scene classes work with `PhaserSceneView.html` since they all extend `PhaserWorldScene` and use the same container initialization pattern.

## Current Pages Using Phaser Scenes

| Page | Scene Type | Current Layout | Migration Priority |
|------|-----------|----------------|-------------------|
| WorldEditorPage | PhaserEditorScene | Custom panel | ✅ High (uses PhaserEditorComponent) |
| WorldViewerPage | PhaserWorldScene | Simple container | ✅ High |
| GameViewerPageDockView | PhaserGameScene | Dockview panel | ⚠️ Medium |
| GameViewerPageGrid | PhaserGameScene | CSS Grid | ⚠️ Medium |
| GameViewerPageMobile | PhaserGameScene | Mobile layout | ⚠️ Low |
| StartGamePage | PhaserWorldScene | Preview | ⚠️ Low |

## Migration Pattern

### TypeScript Integration (Common for All Scene Types)

All scene classes share the same initialization pattern:

```typescript
class MyPage {
  private scene: PhaserWorldScene | PhaserEditorScene | PhaserGameScene;

  async initializeScene() {
    // Get container by SceneId specified in template
    const container = document.getElementById('my-scene-id');

    if (!container) {
      throw new Error('Scene container not found');
    }

    // Create appropriate scene type
    this.scene = new PhaserEditorScene(container, this.eventBus, this.debugMode);
    // OR
    // this.scene = new PhaserGameScene(container, this.eventBus, this.debugMode);
    // OR
    // this.scene = new PhaserWorldScene(container, this.eventBus, this.debugMode);

    // Initialize using LCMComponent lifecycle
    await this.scene.performLocalInit();
    this.scene.setupDependencies();
    await this.scene.activate();
  }

  // Access regions (same for all scene types)
  toggleToolbar() {
    const north = document.getElementById('phaser-scene-view-north');
    if (north) {
      north.classList.toggle('hidden');
    }
  }
}
```

## Page-Specific Migrations

### 1. WorldEditorPage (Priority: ✅ High)

**Current**: Uses `PhaserPanel.html` with toolbar + scene
**Scene Type**: `PhaserEditorScene`
**Target Layout**: North (toolbar) + Center (scene)

#### Template Changes

Replace `/templates/panels/PhaserPanel.html` with:

```html
{{/* Define Editor Toolbar */}}
{{ define "WorldEditorToolbar" }}
<div class="bg-white dark:bg-gray-800 border-b border-gray-300 dark:border-gray-600 p-2">
  <div class="flex items-center justify-between space-x-4">
    <!-- Clear Button, Brush Size, View Options -->
    <!-- (Keep existing toolbar content) -->
  </div>
</div>
{{ end }}

{{/* Use PhaserSceneView */}}
{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "WorldEditorToolbar" .)
  "SceneId" "editor-canvas-container"
  "CenterClass" "bg-gray-100 dark:bg-gray-900 p-2"
}}
```

#### TypeScript Changes

**File**: `/web/src/PhaserEditorComponent.ts`

```typescript
// BEFORE
private setupPhaserContainer(): void {
  let phaserContainer = this.findElement('#editor-canvas-container');
  // ... custom setup ...
}

// AFTER (No changes needed!)
private setupPhaserContainer(): void {
  // PhaserSceneView already provides #editor-canvas-container
  // with proper constraints - just find it
  let phaserContainer = this.findElement('#editor-canvas-container');

  if (!phaserContainer) {
    throw new Error('Scene container not found');
  }

  // Container already has proper sizing constraints from template
}
```

**Key Point**: The container ID stays the same (`editor-canvas-container`), so TypeScript code requires minimal changes!

---

### 2. WorldViewerPage (Priority: ✅ High)

**Current**: Simple div container
**Scene Type**: `PhaserWorldScene`
**Target Layout**: Center only (no toolbars)

#### Template Changes

**File**: `/templates/WorldViewerPage.html`

```html
{{/* BEFORE */}}
<div id="world-viewer-container" style="width: 100%; height: 100vh;">
  <div id="phaser-container"></div>
</div>

{{/* AFTER */}}
{{ template "components/PhaserSceneView.html" dict
  "SceneId" "world-viewer-scene"
  "CenterClass" "bg-gray-900"
}}
```

#### TypeScript Changes

**File**: `/web/src/WorldViewerPage.ts`

```typescript
// BEFORE
const container = document.getElementById('phaser-container');

// AFTER
const container = document.getElementById('world-viewer-scene');

// Rest stays the same!
this.scene = new PhaserWorldScene(container, this.eventBus);
await this.scene.performLocalInit();
await this.scene.activate();
```

---

### 3. GameViewerPageDockView (Priority: ⚠️ Medium)

**Current**: Dockview panel with custom layout
**Scene Type**: `PhaserGameScene`
**Target Layout**: Optional North (controls) + Center (scene)

#### Template Changes

**File**: `/templates/GameViewerPageDockView.html` (or panel component)

```html
{{/* Define Optional Game Controls */}}
{{ define "GameViewerControls" }}
<div id="game-controls-bar" class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-2">
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-2">
      <button id="end-turn-btn" class="px-3 py-1 text-xs bg-green-600 text-white rounded">End Turn</button>
      <button id="cancel-btn" class="px-3 py-1 text-xs bg-gray-600 text-white rounded">Cancel</button>
    </div>
    <button id="toggle-controls" class="text-xs text-gray-600 hover:text-gray-800">
      Hide Controls
    </button>
  </div>
</div>
{{ end }}

{{/* Use PhaserSceneView with optional toolbar */}}
{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "GameViewerControls" .)
  "SceneId" "game-viewer-scene"
  "CenterClass" "bg-gray-900"
}}
```

#### TypeScript Changes

**File**: `/web/src/GameViewerPageDockView.ts`

```typescript
private initializeScene() {
  const container = document.getElementById('game-viewer-scene');

  if (!container) {
    throw new Error('Game viewer scene container not found');
  }

  this.gameScene = new PhaserGameScene(container, this.eventBus, this.debugMode);
  await this.gameScene.performLocalInit();
  this.gameScene.setupDependencies();
  await this.gameScene.activate();
}

// Toggle controls visibility
private setupControlsToggle() {
  const toggleBtn = document.getElementById('toggle-controls');
  const northRegion = document.getElementById('phaser-scene-view-north');

  toggleBtn?.addEventListener('click', () => {
    northRegion?.classList.toggle('hidden');
    toggleBtn.textContent = northRegion?.classList.contains('hidden')
      ? 'Show Controls'
      : 'Hide Controls';
  });
}
```

---

### 4. GameViewerPageGrid (Priority: ⚠️ Medium)

**Current**: CSS Grid with game log sidebar
**Scene Type**: `PhaserGameScene`
**Target Layout**: West (game log) + Center (scene) OR stay with grid

**Note**: This page uses CSS Grid for layout and may not benefit from PhaserSceneView. The scene container portion can still use PhaserSceneView for consistency:

```html
<div id="grid-game-scene-container">
  {{/* Use PhaserSceneView for scene area only */}}
  {{ template "components/PhaserSceneView.html" dict
    "SceneId" "game-viewer-scene"
    "CenterClass" "bg-gray-900"
  }}
</div>
```

---

### 5. GameViewerPageMobile (Priority: ⚠️ Low)

**Current**: Mobile-optimized with calc() heights
**Scene Type**: `PhaserGameScene`
**Target Layout**: North (header) + Center (scene) + South (action bar)

#### Template Changes

```html
{{ define "MobileHeader" }}
<div id="mobile-header" class="bg-white dark:bg-gray-800 border-b border-gray-200">
  <!-- Header content -->
</div>
{{ end }}

{{ define "MobileActionBar" }}
<div id="mobile-bottom-bar" class="bg-white dark:bg-gray-800 border-t border-gray-200">
  <!-- Action buttons -->
</div>
{{ end }}

{{/* Use PhaserSceneView */}}
{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "MobileHeader" .)
  "SouthContent" (template "MobileActionBar" .)
  "SceneId" "mobile-game-scene"
  "CenterClass" "bg-gray-900"
  "WrapperClass" "mobile-game-wrapper"
}}
```

**Benefits**: No more `calc(100vh - 70px - 64px)` - flexbox handles it automatically!

---

### 6. StartGamePage (Priority: ⚠️ Low)

**Current**: World preview for game creation
**Scene Type**: `PhaserWorldScene`
**Target Layout**: Simple center-only

```html
{{ template "components/PhaserSceneView.html" dict
  "SceneId" "start-game-preview-scene"
  "CenterClass" "bg-gray-900 rounded-lg"
}}
```

---

## Migration Checklist

For each page migration:

### Template Updates
- [ ] Create region content templates (NorthContent, SouthContent, etc.)
- [ ] Replace existing layout with `PhaserSceneView.html` template call
- [ ] Specify `SceneId` matching TypeScript expectations
- [ ] Add appropriate `CenterClass` for styling

### TypeScript Updates
- [ ] Update container ID lookup to match `SceneId`
- [ ] Verify scene initialization uses correct scene type
- [ ] Update region access to use standard IDs (`phaser-scene-view-*`)
- [ ] Test toolbar toggle/hide functionality (if applicable)

### Testing
- [ ] Verify parent resizes → canvas resizes smoothly
- [ ] Check toolbar wraps correctly on narrow screens
- [ ] Confirm no scrollbars appear unexpectedly
- [ ] Test region visibility toggles (if used)
- [ ] Verify scene remains functional with all scene types

---

## Common Patterns by Scene Type

### PhaserEditorScene (WorldEditorPage)

**Typical Layout**: Toolbar + Scene
```html
{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "EditorToolbar" .)
  "SceneId" "editor-canvas-container"
  "CenterClass" "bg-gray-100 dark:bg-gray-900 p-2"
}}
```

### PhaserGameScene (GameViewerPages)

**Typical Layout**: Optional Controls + Scene + Optional Status
```html
{{ template "components/PhaserSceneView.html" dict
  "NorthContent" (template "GameControls" .)
  "SouthContent" (template "StatusBar" .)
  "SceneId" "game-viewer-scene"
  "CenterClass" "bg-gray-900"
}}
```

### PhaserWorldScene (Simple Viewers)

**Typical Layout**: Scene Only
```html
{{ template "components/PhaserSceneView.html" dict
  "SceneId" "world-viewer-scene"
  "CenterClass" "bg-gray-900"
}}
```

---

## Key Advantages

1. **Scene Type Agnostic**: Works with all scene classes since they share the same container pattern
2. **Minimal TypeScript Changes**: Container IDs stay the same
3. **Consistent Behavior**: All pages get circular sizing fix automatically
4. **Maintainable**: One template fixes issues in all pages
5. **Flexible**: Each page can customize regions as needed

---

## Next Steps

1. **Start with WorldEditorPage** (already using PhaserEditorComponent - highest impact)
2. **Migrate WorldViewerPage** (simplest case - good test)
3. **Update GameViewerPages** (more complex but high value)
4. **Test thoroughly** with all scene types
5. **Document any edge cases** discovered during migration
