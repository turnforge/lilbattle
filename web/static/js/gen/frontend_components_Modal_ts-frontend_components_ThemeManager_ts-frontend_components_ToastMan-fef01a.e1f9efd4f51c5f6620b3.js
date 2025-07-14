"use strict";
(this["webpackChunkweewar"] = this["webpackChunkweewar"] || []).push([["frontend_components_Modal_ts-frontend_components_ThemeManager_ts-frontend_components_ToastMan-fef01a"],{

/***/ "./frontend/components/Modal.ts":
/*!**************************************!*\
  !*** ./frontend/components/Modal.ts ***!
  \**************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   Modal: () => (/* binding */ Modal)
/* harmony export */ });
/* harmony import */ var _TemplateLoader__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./TemplateLoader */ "./frontend/components/TemplateLoader.ts");

class Modal {
    constructor() {
        this.currentTemplateId = null;
        this.currentData = null;
        this.onSubmitCallback = null;
        this.onApplyCallback = null;
        this.modalContainer = document.getElementById('modal-container');
        this.modalBackdrop = document.getElementById('modal-backdrop');
        this.modalPanel = document.getElementById('modal-panel');
        this.modalContent = document.getElementById('modal-content');
        this.closeButton = document.getElementById('modal-close');
        this.templateLoader = new _TemplateLoader__WEBPACK_IMPORTED_MODULE_0__.TemplateLoader();
        this.bindEvents();
    }
    static getInstance() {
        if (!Modal.instance) {
            Modal.instance = new Modal();
        }
        return Modal.instance;
    }
    bindEvents() {
        if (this.closeButton) {
            this.closeButton.addEventListener('click', () => this.hide());
        }
        if (this.modalBackdrop) {
            this.modalBackdrop.addEventListener('click', (e) => {
                if (e.target === this.modalBackdrop) {
                    this.hide();
                }
            });
        }
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isVisible()) {
                this.hide();
            }
        });
        if (this.modalPanel) {
            this.modalPanel.addEventListener('click', (e) => {
                const target = e.target;
                const closeButton = target.closest('button[id$="-cancel"], button[id$="-close"]');
                if (closeButton) {
                    console.log(`Modal cancel/close button clicked: ${closeButton.id}`);
                    this.hide();
                    return;
                }
                const actionButton = target.closest('button[data-modal-action]');
                if (actionButton) {
                    const action = actionButton.getAttribute('data-modal-action');
                    console.log(`Modal action button clicked: ${action}`);
                    if (action === 'submit' && this.onSubmitCallback) {
                        this.onSubmitCallback(this.currentData);
                    }
                    else if (action === 'apply' && this.onApplyCallback) {
                        this.onApplyCallback(this.currentData);
                        this.hide();
                    }
                }
            });
        }
    }
    isVisible() {
        return this.modalContainer ? !this.modalContainer.classList.contains('hidden') : false;
    }
    show(templateId, data = null) {
        if (!this.modalContainer || !this.modalContent) {
            console.error("Modal container or content area not found.");
            return null;
        }
        const success = this.templateLoader.loadInto(templateId, this.modalContent);
        if (!success) {
            this.modalContainer.classList.remove('hidden');
            setTimeout(() => this.modalContainer.classList.add('modal-active'), 10);
            return null;
        }
        this.currentTemplateId = templateId;
        this.currentData = data || {};
        this.onSubmitCallback = (data === null || data === void 0 ? void 0 : data.onSubmit) || null;
        this.onApplyCallback = (data === null || data === void 0 ? void 0 : data.onApply) || null;
        if (data) {
            Object.entries(data).forEach(([key, value]) => {
                if (key !== 'onSubmit' && (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean')) {
                    if (this.modalContent)
                        this.modalContent.dataset[key] = String(value);
                }
            });
        }
        this.modalContainer.classList.remove('hidden');
        setTimeout(() => {
            this.modalContainer.classList.add('modal-active');
        }, 10);
        const firstElement = this.modalContent;
        return firstElement;
    }
    hide() {
        return new Promise((resolve) => {
            if (!this.modalContainer)
                return;
            this.modalContainer.classList.remove('modal-active');
            setTimeout(() => {
                this.modalContainer.classList.add('hidden');
                this.currentTemplateId = null;
                this.currentData = null;
                this.onSubmitCallback = null;
                this.onApplyCallback = null;
                if (this.modalContent)
                    this.modalContent.innerHTML = '';
                resolve();
            }, 200);
        });
    }
    getContentElement() {
        return this.modalContent;
    }
    getCurrentTemplate() {
        return this.currentTemplateId;
    }
    getCurrentData() {
        return this.currentData;
    }
    updateData(newData) {
        this.currentData = Object.assign(Object.assign({}, this.currentData), newData);
        if (this.modalContent && newData) {
            Object.entries(newData).forEach(([key, value]) => {
                if (key !== 'onSubmit' && (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean')) {
                    if (this.modalContent)
                        this.modalContent.dataset[key] = String(value);
                }
            });
        }
    }
    static init() {
        return Modal.getInstance();
    }
}
Modal.instance = null;


/***/ }),

/***/ "./frontend/components/TemplateLoader.ts":
/*!***********************************************!*\
  !*** ./frontend/components/TemplateLoader.ts ***!
  \***********************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   TemplateLoader: () => (/* binding */ TemplateLoader)
/* harmony export */ });
class TemplateLoader {
    constructor(registryName = "template-registry") {
        this.registryName = registryName;
    }
    _findTemplateWrapper(templateId) {
        const templateRegistry = document.getElementById(this.registryName);
        if (!templateRegistry) {
            console.error(`Template registry '#${this.registryName}' not found!`);
            return null;
        }
        const templateWrapper = templateRegistry.querySelector(`[data-template-id="${templateId}"]`);
        if (!templateWrapper) {
            console.error(`Template with ID "${templateId}" not found in registry '#${this.registryName}'.`);
            return null;
        }
        return templateWrapper;
    }
    loadHtml(templateId) {
        const templateWrapper = this._findTemplateWrapper(templateId);
        if (!templateWrapper) {
            return null;
        }
        return templateWrapper.innerHTML;
    }
    load(templateId) {
        const templateWrapper = this._findTemplateWrapper(templateId);
        if (!templateWrapper) {
            return [];
        }
        const templateRootElement = templateWrapper.cloneNode(true);
        if (!templateRootElement) {
            console.error(`Template content is empty for: ${templateId}`);
            return [];
        }
        return Array.from(templateRootElement.children);
    }
    loadInto(templateId, targetElement) {
        if (!targetElement) {
            console.error(`Cannot load template "${templateId}": Target element is null.`);
            return false;
        }
        const templateWrapper = this._findTemplateWrapper(templateId);
        if (!templateWrapper) {
            targetElement.innerHTML = `<div class="p-4 text-red-500">Error loading template '${templateId}' (Not Found)</div>`;
            return false;
        }
        targetElement.innerHTML = '';
        const childElements = Array.from(templateWrapper.children);
        if (childElements.length === 0) {
            console.warn(`Template "${templateId}" has no child elements to load.`);
        }
        else {
            childElements.forEach(child => {
                targetElement.appendChild(child.cloneNode(true));
            });
        }
        return true;
    }
}


/***/ }),

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

}]);
//# sourceMappingURL=frontend_components_Modal_ts-frontend_components_ThemeManager_ts-frontend_components_ToastMan-fef01a.e1f9efd4f51c5f6620b3.js.map