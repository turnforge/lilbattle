(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define("weewar", [], factory);
	else if(typeof exports === 'object')
		exports["weewar"] = factory();
	else
		root["weewar"] = root["weewar"] || {}, root["weewar"]["GameInstanceDetailsPage"] = factory();
})(this, () => {
return /******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./frontend/components/GameInstanceDetailsPage.ts":
/*!********************************************************!*\
  !*** ./frontend/components/GameInstanceDetailsPage.ts ***!
  \********************************************************/
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



class GameInstanceDetailsPage {
    constructor() {
        this.themeManager = null;
        this.modal = null;
        this.toastManager = null;
        this.themeToggleButton = null;
        this.themeToggleIcon = null;
        this.currentGameInstanceId = null;
        this.isLoadingGameInstance = false;
        this.initializeComponents();
        this.bindEvents();
        this.loadInitialState();
    }
    initializeComponents() {
        const designIdInput = document.getElementById("designIdInput");
        const designId = (designIdInput === null || designIdInput === void 0 ? void 0 : designIdInput.value.trim()) || null;
        _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.init();
        this.modal = _Modal__WEBPACK_IMPORTED_MODULE_1__.Modal.init();
        this.toastManager = _ToastManager__WEBPACK_IMPORTED_MODULE_2__.ToastManager.init();
        this.themeToggleButton = document.getElementById('theme-toggle-button');
        this.themeToggleIcon = document.getElementById('theme-toggle-icon');
        if (!this.themeToggleButton || !this.themeToggleIcon) {
            console.warn("Theme toggle button or icon element not found in Header.");
        }
        console.log('LeetCoach application initialized');
    }
    bindEvents() {
        if (this.themeToggleButton) {
            this.themeToggleButton.addEventListener('click', this.handleThemeToggleClick.bind(this));
        }
        const mobileMenuButton = document.getElementById('mobile-menu-button');
        if (mobileMenuButton) {
            mobileMenuButton.addEventListener('click', () => {
            });
        }
        const saveButton = document.querySelector('header button.bg-blue-600');
        if (saveButton) {
            saveButton.addEventListener('click', this.saveDocument.bind(this));
        }
        const exportButton = document.querySelector('header button.bg-gray-200');
        if (exportButton) {
            exportButton.addEventListener('click', this.exportDocument.bind(this));
        }
    }
    loadInitialState() {
        var _a;
        this.updateThemeButtonState();
        const designIdInput = document.getElementById("designIdInput");
        const designId = (designIdInput === null || designIdInput === void 0 ? void 0 : designIdInput.value.trim()) || null;
        if (designId) {
            this.currentGameInstanceId = designId;
            console.log(`Found GameInstance ID: ${this.currentGameInstanceId}. Loading data...`);
            this.loadGameInstanceData(this.currentGameInstanceId);
        }
        else {
            console.error("GameInstance ID input element not found or has no value. Cannot load document.");
            (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast("Error", "Could not load document: GameInstance ID missing.", "error");
        }
    }
    loadGameInstanceData(designId) {
        return __awaiter(this, void 0, void 0, function* () {
            console.log(`GameInstanceDetailsPage: Loading design ${designId}...`);
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
    saveDocument() {
        var _a;
        console.log("Save button clicked (Placeholder - Requires API integration for full save)");
        (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Save Action', 'Incremental saves handle updates. Full save TBD.', 'info');
    }
    exportDocument() {
        if (this.toastManager) {
            this.toastManager.showToast('Export started', 'Your document is being prepared for export.', 'info');
            setTimeout(() => {
                var _a;
                (_a = this.toastManager) === null || _a === void 0 ? void 0 : _a.showToast('Export complete', 'Document export simulation finished.', 'success');
            }, 1500);
        }
    }
}
document.addEventListener('DOMContentLoaded', () => {
    const lc = new GameInstanceDetailsPage();
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
/******/ 			"GameInstanceDetailsPage": 0
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
/******/ 	var __webpack_exports__ = __webpack_require__.O(undefined, ["frontend_components_Modal_ts-frontend_components_ThemeManager_ts-frontend_components_ToastMan-fef01a"], () => (__webpack_require__("./frontend/components/GameInstanceDetailsPage.ts")))
/******/ 	__webpack_exports__ = __webpack_require__.O(__webpack_exports__);
/******/ 	
/******/ 	return __webpack_exports__;
/******/ })()
;
});
//# sourceMappingURL=GameInstanceDetailsPage.6689fbd3b0f706d0dd42.js.map