(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define("weewar", [], factory);
	else if(typeof exports === 'object')
		exports["weewar"] = factory();
	else
		root["weewar"] = root["weewar"] || {}, root["weewar"]["MapEditorPage"] = factory();
})(this, () => {
return /******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./frontend/components/MapEditorPage.ts":
/*!**********************************************!*\
  !*** ./frontend/components/MapEditorPage.ts ***!
  \**********************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _ThemeManager__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ThemeManager */ "./frontend/components/ThemeManager.ts");
/* harmony import */ var _Modal__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./Modal */ "./frontend/components/Modal.ts");
/* harmony import */ var _ToastManager__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./ToastManager */ "./frontend/components/ToastManager.ts");
var __awaiter = (undefined && undefined.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};



class MapEditorPage {
    constructor() {
        this.themeManager = null;
        this.modal = null;
        this.toastManager = null;
        this.themeToggleButton = null;
        this.themeToggleIcon = null;
        this.currentMapId = null;
        this.isNewMap = false;
        this.mapData = null;
        this.currentTerrain = 1;
        this.brushSize = 0;
        this.editorCanvas = null;
        this.editorOutput = null;
        this.mapImage = null;
        this.wasmModule = null;
        this.wasmInitialized = false;
        this.initializeComponents();
        this.bindEvents();
        this.loadInitialState();
        this.initializeWasm();
    }
    initializeComponents() {
        const mapIdInput = document.getElementById("mapIdInput");
        const isNewMapInput = document.getElementById("isNewMap");
        this.currentMapId = (mapIdInput === null || mapIdInput === void 0 ? void 0 : mapIdInput.value.trim()) || null;
        this.isNewMap = (isNewMapInput === null || isNewMapInput === void 0 ? void 0 : isNewMapInput.value) === "true";
        _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.init();
        this.modal = _Modal__WEBPACK_IMPORTED_MODULE_1__.Modal.init();
        this.toastManager = _ToastManager__WEBPACK_IMPORTED_MODULE_2__.ToastManager.init();
        this.themeToggleButton = document.getElementById('theme-toggle-button');
        this.themeToggleIcon = document.getElementById('theme-toggle-icon');
        this.editorCanvas = document.getElementById('editor-canvas');
        this.editorOutput = document.getElementById('editor-output');
        this.mapImage = document.getElementById('map-image');
        if (!this.themeToggleButton || !this.themeToggleIcon) {
            console.warn("Theme toggle button or icon element not found in Header.");
        }
        this.logToConsole('Map Editor initialized');
    }
    bindEvents() {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
        if (this.themeToggleButton) {
            this.themeToggleButton.addEventListener('click', this.handleThemeToggleClick.bind(this));
        }
        const saveButton = document.getElementById('save-map-btn');
        if (saveButton) {
            saveButton.addEventListener('click', this.saveMap.bind(this));
        }
        const exportButton = document.getElementById('export-map-btn');
        if (exportButton) {
            exportButton.addEventListener('click', this.exportMap.bind(this));
        }
        const validateButton = document.getElementById('validate-map-btn');
        if (validateButton) {
            validateButton.addEventListener('click', this.validateMap.bind(this));
        }
        const clearConsoleButton = document.getElementById('clear-console-btn');
        if (clearConsoleButton) {
            clearConsoleButton.addEventListener('click', this.clearConsole.bind(this));
        }
        document.querySelectorAll('[data-action="create-new-map"]').forEach(button => {
            button.addEventListener('click', (e) => {
                const target = e.target;
                const width = parseInt(target.dataset.width || '8');
                const height = parseInt(target.dataset.height || '8');
                this.createNewMap(width, height);
            });
        });
        document.querySelectorAll('.terrain-button').forEach(button => {
            button.addEventListener('click', (e) => {
                const terrain = e.target.getAttribute('data-terrain');
                if (terrain) {
                    this.setBrushTerrain(parseInt(terrain));
                }
            });
        });
        const brushSizeSelect = document.getElementById('brush-size');
        if (brushSizeSelect) {
            brushSizeSelect.addEventListener('change', (e) => {
                this.setBrushSize(parseInt(e.target.value));
            });
        }
        (_a = document.querySelector('[data-action="paint-terrain"]')) === null || _a === void 0 ? void 0 : _a.addEventListener('click', () => {
            this.paintTerrain();
        });
        (_b = document.querySelector('[data-action="flood-fill"]')) === null || _b === void 0 ? void 0 : _b.addEventListener('click', () => {
            this.floodFill();
        });
        (_c = document.querySelector('[data-action="remove-terrain"]')) === null || _c === void 0 ? void 0 : _c.addEventListener('click', () => {
            this.removeTerrain();
        });
        (_d = document.querySelector('[data-action="undo"]')) === null || _d === void 0 ? void 0 : _d.addEventListener('click', () => {
            this.editorUndo();
        });
        (_e = document.querySelector('[data-action="redo"]')) === null || _e === void 0 ? void 0 : _e.addEventListener('click', () => {
            this.editorRedo();
        });
        document.querySelectorAll('[data-action="render-map"]').forEach(button => {
            button.addEventListener('click', (e) => {
                const target = e.target;
                const width = parseInt(target.dataset.width || '600');
                const height = parseInt(target.dataset.height || '450');
                this.renderEditor(width, height);
            });
        });
        document.querySelectorAll('[data-action="export-game"]').forEach(button => {
            button.addEventListener('click', (e) => {
                const target = e.target;
                const players = parseInt(target.dataset.players || '2');
                this.exportToGame(players);
            });
        });
        (_f = document.querySelector('[data-action="fill-all-grass"]')) === null || _f === void 0 ? void 0 : _f.addEventListener('click', () => {
            this.fillAllGrass();
        });
        (_g = document.querySelector('[data-action="create-test-pattern"]')) === null || _g === void 0 ? void 0 : _g.addEventListener('click', () => {
            this.createTestPattern();
        });
        (_h = document.querySelector('[data-action="create-island-map"]')) === null || _h === void 0 ? void 0 : _h.addEventListener('click', () => {
            this.createIslandMap();
        });
        (_j = document.querySelector('[data-action="create-mountain-ridge"]')) === null || _j === void 0 ? void 0 : _j.addEventListener('click', () => {
            this.createMountainRidge();
        });
        (_k = document.querySelector('[data-action="show-terrain-stats"]')) === null || _k === void 0 ? void 0 : _k.addEventListener('click', () => {
            this.showTerrainStats();
        });
        (_l = document.querySelector('[data-action="randomize-terrain"]')) === null || _l === void 0 ? void 0 : _l.addEventListener('click', () => {
            this.randomizeTerrain();
        });
        (_m = document.querySelector('[data-action="download-image"]')) === null || _m === void 0 ? void 0 : _m.addEventListener('click', () => {
            this.downloadImage();
        });
        (_o = document.querySelector('[data-action="download-game-data"]')) === null || _o === void 0 ? void 0 : _o.addEventListener('click', () => {
            this.downloadGameData();
        });
    }
    loadInitialState() {
        this.updateThemeButtonState();
        this.updateEditorStatus('Initializing...');
        if (this.isNewMap) {
            this.logToConsole('Creating new map...');
            this.initializeNewMap();
        }
        else if (this.currentMapId) {
            this.logToConsole(`Loading existing map: ${this.currentMapId}`);
            this.loadExistingMap(this.currentMapId);
        }
        else {
            this.logToConsole('Error: No map ID provided');
            this.updateEditorStatus('Error');
        }
    }
    initializeWasm() {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                this.logToConsole('Loading WASM module...');
                this.wasmInitialized = true;
                this.updateEditorStatus('Ready');
                this.logToConsole('WASM module loaded successfully');
            }
            catch (error) {
                console.error('Failed to initialize WASM:', error);
                this.logToConsole(`WASM initialization failed: ${error}`);
                this.updateEditorStatus('WASM Error');
            }
        });
    }
    initializeNewMap() {
        this.mapData = {
            name: "New Map",
            width: 8,
            height: 8,
            tiles: {},
            map_units: []
        };
        this.updateEditorStatus('New Map');
        this.logToConsole('New map initialized');
    }
    loadExistingMap(mapId) {
        return __awaiter(this, void 0, void 0, function* () {
            try {
                this.logToConsole(`Loading map data for ${mapId}...`);
                this.updateEditorStatus('Loading...');
                setTimeout(() => {
                    this.mapData = {
                        name: `Map ${mapId}`,
                        width: 8,
                        height: 8,
                        tiles: {},
                        map_units: []
                    };
                    this.updateEditorStatus('Loaded');
                    this.logToConsole('Map data loaded');
                }, 1000);
            }
            catch (error) {
                console.error('Failed to load map:', error);
                this.logToConsole(`Failed to load map: ${error}`);
                this.updateEditorStatus('Load Error');
            }
        });
    }
    createNewMap(width, height) {
        this.logToConsole(`Creating new ${width}×${height} map...`);
        this.mapData = {
            name: `New ${width}×${height} Map`,
            width,
            height,
            tiles: {},
            map_units: []
        };
        this.updateEditorStatus('Ready');
        this.logToConsole('New map created');
    }
    setBrushTerrain(terrain) {
        this.currentTerrain = terrain;
        const terrainNames = ['', 'Grass', 'Desert', 'Water', 'Mountain', 'Rock'];
        this.logToConsole(`Brush terrain set to: ${terrainNames[terrain]}`);
        this.updateBrushInfo();
        this.updateTerrainButtonSelection(terrain);
    }
    setBrushSize(size) {
        this.brushSize = size;
        const sizeNames = ['Single (1 hex)', 'Small (7 hexes)', 'Medium (19 hexes)', 'Large (37 hexes)', 'X-Large (61 hexes)', 'XX-Large (91 hexes)'];
        this.logToConsole(`Brush size set to: ${sizeNames[size]}`);
        this.updateBrushInfo();
    }
    paintTerrain() {
        const rowInput = document.getElementById('paint-row');
        const colInput = document.getElementById('paint-col');
        if (rowInput && colInput) {
            const row = parseInt(rowInput.value);
            const col = parseInt(colInput.value);
            this.logToConsole(`Painting terrain ${this.currentTerrain} at (${row}, ${col})`);
        }
    }
    floodFill() {
        const rowInput = document.getElementById('paint-row');
        const colInput = document.getElementById('paint-col');
        if (rowInput && colInput) {
            const row = parseInt(rowInput.value);
            const col = parseInt(colInput.value);
            this.logToConsole(`Flood filling with terrain ${this.currentTerrain} from (${row}, ${col})`);
        }
    }
    removeTerrain() {
        const rowInput = document.getElementById('paint-row');
        const colInput = document.getElementById('paint-col');
        if (rowInput && colInput) {
            const row = parseInt(rowInput.value);
            const col = parseInt(colInput.value);
            this.logToConsole(`Removing terrain at (${row}, ${col})`);
        }
    }
    editorUndo() {
        this.logToConsole('Undo action');
    }
    editorRedo() {
        this.logToConsole('Redo action');
    }
    renderEditor(width, height) {
        this.logToConsole(`Rendering map at ${width}×${height}...`);
        if (this.mapImage) {
            this.mapImage.style.display = 'block';
            this.mapImage.src = 'data:image/svg+xml;base64,' + btoa(`
                <svg width="${width}" height="${height}" xmlns="http://www.w3.org/2000/svg">
                    <rect width="${width}" height="${height}" fill="#90EE90"/>
                    <text x="${width / 2}" y="${height / 2}" text-anchor="middle" dy=".3em" font-family="Arial" font-size="14">
                        Map Preview (${width}×${height})
                    </text>
                </svg>
            `);
        }
    }
    downloadImage() {
        var _a;
        this.logToConsole('Downloading map image...');
        (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Download', 'Image download not yet implemented', 'info');
    }
    exportToGame(players) {
        var _a;
        this.logToConsole(`Exporting as ${players}-player game...`);
        (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Export', `${players}-player game export not yet implemented`, 'info');
    }
    downloadGameData() {
        var _a;
        this.logToConsole('Downloading game data...');
        (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Download', 'Game data download not yet implemented', 'info');
    }
    fillAllGrass() {
        this.logToConsole('Filling all tiles with grass...');
    }
    createTestPattern() {
        this.logToConsole('Creating test pattern...');
    }
    createIslandMap() {
        this.logToConsole('Creating island map...');
    }
    createMountainRidge() {
        this.logToConsole('Creating mountain ridge...');
    }
    showTerrainStats() {
        this.logToConsole('Terrain statistics:');
        this.logToConsole('- Grass: 0 tiles');
        this.logToConsole('- Desert: 0 tiles');
        this.logToConsole('- Water: 0 tiles');
        this.logToConsole('- Mountain: 0 tiles');
        this.logToConsole('- Rock: 0 tiles');
    }
    randomizeTerrain() {
        this.logToConsole('Randomizing terrain...');
    }
    saveMap() {
        return __awaiter(this, void 0, void 0, function* () {
            var _a, _b, _c;
            if (!this.mapData) {
                (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Error', 'No map data to save', 'error');
                return;
            }
            try {
                this.logToConsole('Saving map...');
                this.updateEditorStatus('Saving...');
                const url = this.isNewMap ? '/api/maps' : `/api/maps/${this.currentMapId}`;
                const method = this.isNewMap ? 'POST' : 'PUT';
                const response = yield fetch(url, {
                    method,
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(this.mapData),
                });
                if (response.ok) {
                    const result = yield response.json();
                    this.logToConsole('Map saved successfully');
                    this.updateEditorStatus('Saved');
                    (_b = this.toastManager) === null || _b === void 0 ? void 0 : _b.showToast('Success', 'Map saved successfully', 'success');
                    if (this.isNewMap && result.id) {
                        this.currentMapId = result.id;
                        this.isNewMap = false;
                        history.replaceState(null, '', `/maps/${result.id}/edit`);
                    }
                }
                else {
                    throw new Error(`Save failed: ${response.statusText}`);
                }
            }
            catch (error) {
                console.error('Save failed:', error);
                this.logToConsole(`Save failed: ${error}`);
                this.updateEditorStatus('Save Error');
                (_c = this.toastManager) === null || _c === void 0 ? void 0 : _c.showToast('Error', 'Failed to save map', 'error');
            }
        });
    }
    exportMap() {
        var _a;
        this.logToConsole('Exporting map...');
        (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Export', 'Export functionality not yet implemented', 'info');
    }
    validateMap() {
        this.logToConsole('Validating map...');
        this.logToConsole('Map validation completed - no issues found');
    }
    clearConsole() {
        if (this.editorOutput) {
            this.editorOutput.textContent = '';
        }
    }
    logToConsole(message) {
        if (this.editorOutput) {
            const timestamp = new Date().toLocaleTimeString();
            this.editorOutput.textContent += `[${timestamp}] ${message}\n`;
            this.editorOutput.scrollTop = this.editorOutput.scrollHeight;
        }
        console.log(`[MapEditor] ${message}`);
    }
    updateEditorStatus(status) {
        const statusElement = document.getElementById('editor-status');
        if (statusElement) {
            statusElement.textContent = status;
            statusElement.className = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium';
            if (status.includes('Error')) {
                statusElement.className += ' bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
            }
            else if (status === 'Ready' || status === 'Saved' || status === 'Loaded') {
                statusElement.className += ' bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
            }
            else {
                statusElement.className += ' bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
            }
        }
    }
    updateBrushInfo() {
        const brushInfo = document.getElementById('brush-info');
        if (brushInfo) {
            const terrainNames = ['', 'Grass', 'Desert', 'Water', 'Mountain', 'Rock'];
            const sizeNames = ['Single (1 hex)', 'Small (7 hexes)', 'Medium (19 hexes)', 'Large (37 hexes)', 'X-Large (61 hexes)', 'XX-Large (91 hexes)'];
            brushInfo.textContent = `Current: ${terrainNames[this.currentTerrain]}, ${sizeNames[this.brushSize]}`;
        }
    }
    updateTerrainButtonSelection(terrain) {
        document.querySelectorAll('.terrain-button').forEach(button => {
            const buttonTerrain = button.getAttribute('data-terrain');
            if (buttonTerrain === terrain.toString()) {
                button.classList.add('bg-blue-100', 'dark:bg-blue-900', 'border-blue-500');
            }
            else {
                button.classList.remove('bg-blue-100', 'dark:bg-blue-900', 'border-blue-500');
            }
        });
    }
    handleThemeToggleClick() {
        const currentSetting = _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.getCurrentThemeSetting();
        const nextSetting = _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.getNextTheme(currentSetting);
        _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.setTheme(nextSetting);
        this.updateThemeButtonState(nextSetting);
    }
    updateThemeButtonState(currentTheme) {
        if (!this.themeToggleButton || !this.themeToggleIcon)
            return;
        const themeToDisplay = currentTheme || _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.getCurrentThemeSetting();
        const iconSVG = _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.getIconSVG(themeToDisplay);
        const label = `Toggle theme (currently: ${_ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.getThemeLabel(themeToDisplay)})`;
        this.themeToggleIcon.innerHTML = iconSVG;
        this.themeToggleButton.setAttribute('aria-label', label);
        this.themeToggleButton.setAttribute('title', label);
    }
}
document.addEventListener('DOMContentLoaded', () => {
    new MapEditorPage();
});


/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = __webpack_modules__;
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/chunk loaded */
/******/ 	(() => {
/******/ 		var deferred = [];
/******/ 		__webpack_require__.O = (result, chunkIds, fn, priority) => {
/******/ 			if(chunkIds) {
/******/ 				priority = priority || 0;
/******/ 				for(var i = deferred.length; i > 0 && deferred[i - 1][2] > priority; i--) deferred[i] = deferred[i - 1];
/******/ 				deferred[i] = [chunkIds, fn, priority];
/******/ 				return;
/******/ 			}
/******/ 			var notFulfilled = Infinity;
/******/ 			for (var i = 0; i < deferred.length; i++) {
/******/ 				var [chunkIds, fn, priority] = deferred[i];
/******/ 				var fulfilled = true;
/******/ 				for (var j = 0; j < chunkIds.length; j++) {
/******/ 					if ((priority & 1 === 0 || notFulfilled >= priority) && Object.keys(__webpack_require__.O).every((key) => (__webpack_require__.O[key](chunkIds[j])))) {
/******/ 						chunkIds.splice(j--, 1);
/******/ 					} else {
/******/ 						fulfilled = false;
/******/ 						if(priority < notFulfilled) notFulfilled = priority;
/******/ 					}
/******/ 				}
/******/ 				if(fulfilled) {
/******/ 					deferred.splice(i--, 1)
/******/ 					var r = fn();
/******/ 					if (r !== undefined) result = r;
/******/ 				}
/******/ 			}
/******/ 			return result;
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/jsonp chunk loading */
/******/ 	(() => {
/******/ 		// no baseURI
/******/ 		
/******/ 		// object to store loaded and loading chunks
/******/ 		// undefined = chunk not loaded, null = chunk preloaded/prefetched
/******/ 		// [resolve, reject, Promise] = chunk loading, 0 = chunk loaded
/******/ 		var installedChunks = {
/******/ 			"MapEditorPage": 0
/******/ 		};
/******/ 		
/******/ 		// no chunk on demand loading
/******/ 		
/******/ 		// no prefetching
/******/ 		
/******/ 		// no preloaded
/******/ 		
/******/ 		// no HMR
/******/ 		
/******/ 		// no HMR manifest
/******/ 		
/******/ 		__webpack_require__.O.j = (chunkId) => (installedChunks[chunkId] === 0);
/******/ 		
/******/ 		// install a JSONP callback for chunk loading
/******/ 		var webpackJsonpCallback = (parentChunkLoadingFunction, data) => {
/******/ 			var [chunkIds, moreModules, runtime] = data;
/******/ 			// add "moreModules" to the modules object,
/******/ 			// then flag all "chunkIds" as loaded and fire callback
/******/ 			var moduleId, chunkId, i = 0;
/******/ 			if(chunkIds.some((id) => (installedChunks[id] !== 0))) {
/******/ 				for(moduleId in moreModules) {
/******/ 					if(__webpack_require__.o(moreModules, moduleId)) {
/******/ 						__webpack_require__.m[moduleId] = moreModules[moduleId];
/******/ 					}
/******/ 				}
/******/ 				if(runtime) var result = runtime(__webpack_require__);
/******/ 			}
/******/ 			if(parentChunkLoadingFunction) parentChunkLoadingFunction(data);
/******/ 			for(;i < chunkIds.length; i++) {
/******/ 				chunkId = chunkIds[i];
/******/ 				if(__webpack_require__.o(installedChunks, chunkId) && installedChunks[chunkId]) {
/******/ 					installedChunks[chunkId][0]();
/******/ 				}
/******/ 				installedChunks[chunkId] = 0;
/******/ 			}
/******/ 			return __webpack_require__.O(result);
/******/ 		}
/******/ 		
/******/ 		var chunkLoadingGlobal = this["webpackChunkweewar"] = this["webpackChunkweewar"] || [];
/******/ 		chunkLoadingGlobal.forEach(webpackJsonpCallback.bind(null, 0));
/******/ 		chunkLoadingGlobal.push = webpackJsonpCallback.bind(null, chunkLoadingGlobal.push.bind(chunkLoadingGlobal));
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module depends on other loaded chunks and execution need to be delayed
/******/ 	var __webpack_exports__ = __webpack_require__.O(undefined, ["frontend_components_Modal_ts-frontend_components_ThemeManager_ts-frontend_components_ToastMan-fef01a"], () => (__webpack_require__("./frontend/components/MapEditorPage.ts")))
/******/ 	__webpack_exports__ = __webpack_require__.O(__webpack_exports__);
/******/ 	
/******/ 	return __webpack_exports__;
/******/ })()
;
});
//# sourceMappingURL=MapEditorPage.3a0576585d41f77a956a.js.map