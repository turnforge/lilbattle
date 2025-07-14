(function webpackUniversalModuleDefinition(root, factory) {
	if(typeof exports === 'object' && typeof module === 'object')
		module.exports = factory();
	else if(typeof define === 'function' && define.amd)
		define("weewar", [], factory);
	else if(typeof exports === 'object')
		exports["weewar"] = factory();
	else
		root["weewar"] = root["weewar"] || {}, root["weewar"]["LoginPage"] = factory();
})(this, () => {
return /******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./frontend/components/ThemeManager.ts":
/*!*********************************************!*\
  !*** ./frontend/components/ThemeManager.ts ***!
  \*********************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   ThemeManager: () => (/* binding */ ThemeManager)
/* harmony export */ });
class ThemeManager {
    static initialize() {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === ThemeManager.DARK ||
            (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            document.documentElement.classList.add('dark');
        }
        else {
            document.documentElement.classList.remove('dark');
        }
    }
    static setTheme(theme) {
        if (theme === ThemeManager.SYSTEM) {
            localStorage.removeItem('theme');
            if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.documentElement.classList.add('dark');
            }
            else {
                document.documentElement.classList.remove('dark');
            }
        }
        else if (theme === ThemeManager.DARK) {
            localStorage.setItem('theme', ThemeManager.DARK);
            document.documentElement.classList.add('dark');
        }
        else {
            localStorage.setItem('theme', ThemeManager.LIGHT);
            document.documentElement.classList.remove('dark');
        }
    }
    static getCurrentThemeSetting() {
        return localStorage.getItem('theme') || ThemeManager.SYSTEM;
    }
    static getNextTheme(currentSetting) {
        if (currentSetting === ThemeManager.LIGHT) {
            return ThemeManager.DARK;
        }
        else if (currentSetting === ThemeManager.DARK) {
            return ThemeManager.SYSTEM;
        }
        else {
            return ThemeManager.LIGHT;
        }
    }
    static getIconSVG(themeSetting) {
        switch (themeSetting) {
            case ThemeManager.LIGHT: return ThemeManager.LIGHT_ICON_SVG;
            case ThemeManager.DARK: return ThemeManager.DARK_ICON_SVG;
            case ThemeManager.SYSTEM:
            default: return ThemeManager.SYSTEM_ICON_SVG;
        }
    }
    static getThemeLabel(themeSetting) {
        switch (themeSetting) {
            case ThemeManager.LIGHT: return "Light Mode";
            case ThemeManager.DARK: return "Dark Mode";
            case ThemeManager.SYSTEM:
            default: return "System Default";
        }
    }
    static init() {
        ThemeManager.initialize();
    }
}
ThemeManager.LIGHT = 'light';
ThemeManager.DARK = 'dark';
ThemeManager.SYSTEM = 'system';
ThemeManager.LIGHT_ICON_SVG = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-full h-full"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" /></svg>`;
ThemeManager.DARK_ICON_SVG = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-full h-full"><path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" /></svg>`;
ThemeManager.SYSTEM_ICON_SVG = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-full h-full"><path stroke-linecap="round" stroke-linejoin="round" d="M9 17.25v1.007a3 3 0 01-.879 2.122L7.5 21h9l-.621-.621A3 3 0 0115 18.257V17.25m6-12V15a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 15V5.25m18 0A2.25 2.25 0 0018.75 3H5.25A2.25 2.25 0 003 5.25m18 0V12a2.25 2.25 0 01-2.25 2.25H5.25A2.25 2.25 0 013 12V5.25" /></svg>`;


/***/ }),

/***/ "./frontend/components/ToastManager.ts":
/*!*********************************************!*\
  !*** ./frontend/components/ToastManager.ts ***!
  \*********************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   ToastManager: () => (/* binding */ ToastManager)
/* harmony export */ });
class ToastManager {
    constructor() {
        this.toasts = new Map();
        this.counter = 0;
        this.container = document.getElementById('toast-container');
        this.template = document.getElementById('toast-template');
    }
    static getInstance() {
        if (!ToastManager.instance) {
            ToastManager.instance = new ToastManager();
        }
        return ToastManager.instance;
    }
    showToast(title, message, type = 'info', duration = 4000) {
        if (!this.container || !this.template)
            return '';
        const id = `toast-${Date.now()}-${this.counter++}`;
        const toast = this.template.cloneNode(true);
        toast.id = id;
        toast.classList.remove('hidden');
        const titleElement = toast.querySelector('.toast-title');
        const messageElement = toast.querySelector('.toast-message');
        if (titleElement)
            titleElement.textContent = title;
        if (messageElement)
            messageElement.textContent = message;
        const iconContainer = toast.querySelector('.flex-shrink-0');
        if (iconContainer) {
            iconContainer.innerHTML = '';
            let icon;
            let borderColor;
            switch (type) {
                case 'success':
                    icon = '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>';
                    borderColor = 'border-green-500';
                    break;
                case 'error':
                    icon = '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" /></svg>';
                    borderColor = 'border-red-500';
                    break;
                case 'warning':
                    icon = '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" /></svg>';
                    borderColor = 'border-yellow-500';
                    break;
                case 'info':
                default:
                    icon = '<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" /></svg>';
                    borderColor = 'border-blue-500';
                    break;
            }
            iconContainer.innerHTML = icon;
            const borderElement = toast.querySelector('.border-l-4');
            if (borderElement) {
                borderElement.className = borderElement.className.replace(/border-[a-z]+-500/g, borderColor);
            }
        }
        const closeButton = toast.querySelector('.toast-close');
        if (closeButton) {
            closeButton.addEventListener('click', () => {
                this.hideToast(id);
            });
        }
        this.container.appendChild(toast);
        this.toasts.set(id, toast);
        setTimeout(() => {
            toast.classList.remove('scale-95', 'opacity-0');
            toast.classList.add('scale-100', 'opacity-100');
        }, 10);
        if (duration > 0) {
            setTimeout(() => {
                this.hideToast(id);
            }, duration);
        }
        return id;
    }
    hideToast(id) {
        const toast = this.toasts.get(id);
        if (!toast)
            return;
        toast.classList.remove('scale-100', 'opacity-100');
        toast.classList.add('scale-95', 'opacity-0');
        setTimeout(() => {
            toast.remove();
            this.toasts.delete(id);
        }, 300);
    }
    hideAllToasts() {
        this.toasts.forEach((_, id) => {
            this.hideToast(id);
        });
    }
    static init() {
        return ToastManager.getInstance();
    }
}
ToastManager.instance = null;


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
/************************************************************************/
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
/************************************************************************/
var __webpack_exports__ = {};
// This entry needs to be wrapped in an IIFE because it needs to be isolated against other modules in the chunk.
(() => {
/*!******************************************!*\
  !*** ./frontend/components/LoginPage.ts ***!
  \******************************************/
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _ThemeManager__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./ThemeManager */ "./frontend/components/ThemeManager.ts");
/* harmony import */ var _ToastManager__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ToastManager */ "./frontend/components/ToastManager.ts");


class LoginPage {
    constructor() {
        this.isSignUpMode = false;
        _ThemeManager__WEBPACK_IMPORTED_MODULE_0__.ThemeManager.init();
        _ToastManager__WEBPACK_IMPORTED_MODULE_1__.ToastManager.init();
        this.callbackURL = document.getElementById("callbackURL");
        this.form = document.getElementById('auth-form');
        this.titleElement = document.getElementById('auth-title');
        this.submitButton = document.getElementById('auth-submit-button');
        this.toggleLink = document.getElementById('auth-toggle-link');
        this.confirmPasswordGroup = document.getElementById('confirm-password-group');
        this.confirmPasswordInput = document.getElementById('confirm-password');
        this.emailInput = document.getElementById('email-address');
        if (!this.form || !this.titleElement || !this.submitButton || !this.toggleLink || !this.confirmPasswordGroup || !this.confirmPasswordInput || !this.emailInput) {
            console.error("LoginPage: Could not find all required authentication form elements.");
            return;
        }
        this.bindEvents();
        this.updateUI();
    }
    bindEvents() {
        var _a;
        (_a = this.toggleLink) === null || _a === void 0 ? void 0 : _a.addEventListener('click', (e) => {
            e.preventDefault();
            this.isSignUpMode = !this.isSignUpMode;
            this.updateUI();
        });
    }
    updateUI() {
        if (!this.form || !this.titleElement || !this.submitButton || !this.toggleLink || !this.confirmPasswordGroup || !this.confirmPasswordInput || !this.emailInput) {
            return;
        }
        if (this.isSignUpMode) {
            this.titleElement.textContent = 'Create your account';
            this.submitButton.textContent = 'Sign up';
            this.toggleLink.textContent = 'Already have an account? Sign In';
            this.confirmPasswordGroup.classList.remove('hidden');
            this.confirmPasswordInput.required = true;
            this.form.action = '/auth/signup?callbackURL=' + this.callbackURL.value;
            this.emailInput.autocomplete = 'email';
            this.form.querySelector('#password').autocomplete = 'new-password';
        }
        else {
            this.titleElement.textContent = 'Sign in to your account';
            this.submitButton.textContent = 'Sign in';
            this.toggleLink.textContent = 'Need an account? Sign Up';
            this.confirmPasswordGroup.classList.add('hidden');
            this.confirmPasswordInput.required = false;
            this.confirmPasswordInput.value = '';
            this.form.action = '/auth/login?callbackURL=' + this.callbackURL.value;
            this.emailInput.autocomplete = 'email';
            this.form.querySelector('#password').autocomplete = 'current-password';
        }
    }
    static init() {
        return new LoginPage();
    }
}
document.addEventListener('DOMContentLoaded', () => {
    LoginPage.init();
});

})();

/******/ 	return __webpack_exports__;
/******/ })()
;
});
//# sourceMappingURL=LoginPage.f3adba71637f38042770.js.map