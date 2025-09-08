import { BaseComponent } from '../lib/Component';
import { EventBus } from '../lib/EventBus';
import { LCMComponent } from '../lib/LCMComponent';
import { RulesTable } from './RulesTable';
import { ITheme } from '../assets/themes/BaseTheme';

interface UnitData {
    unitType: number;
    player: number;
}

/**
 * DamageDistributionPanel displays damage distribution between units
 * 
 * This component shows combat damage distributions when one unit attacks another,
 * displaying histograms of possible damage values and their probabilities.
 * The panel shows damage distributions for the selected unit against all other unit types.
 */
export class DamageDistributionPanel extends BaseComponent implements LCMComponent {
    private isUIBound = false;
    private isActivated = false;
    private currentUnit: UnitData | null = null;
    public rulesTable: RulesTable;
    private theme: ITheme | null = null;

    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('damage-distribution-panel', rootElement, eventBus, debugMode);
        this.rulesTable = new RulesTable();
    }

    // LCMComponent Phase 1: Initialize DOM structure
    public performLocalInit(): LCMComponent[] {
        if (this.isUIBound) {
            this.log('Already bound to DOM, skipping');
            return [];
        }

        this.log('Binding DamageDistributionPanel to DOM using template');
        this.isUIBound = true;
        this.log('DamageDistributionPanel bound to DOM successfully');
        
        // This is a leaf component - no children
        return [];
    }

    // Phase 2: No external dependencies needed
    public setupDependencies(): void {
        this.log('DamageDistributionPanel: No dependencies required');
    }

    // Phase 3: Activate component
    public activate(): void {
        if (this.isActivated) {
            this.log('Already activated, skipping');
            return;
        }

        this.log('Activating DamageDistributionPanel');
        this.isActivated = true;
        this.log('DamageDistributionPanel activated successfully');
    }

    // Phase 4: Deactivate component
    public deactivate(): void {
        this.log('Deactivating DamageDistributionPanel');
        this.currentUnit = null;
        this.isActivated = false;
        this.log('DamageDistributionPanel deactivated');
    }

    /**
     * Set the theme for getting unit names
     */
    public setTheme(theme: ITheme): void {
        this.theme = theme;
    }

    /**
     * Update the panel with damage distributions for a selected unit
     */
    public updateUnitInfo(unit: UnitData): void {
        if (!this.isActivated) {
            throw new Error('Component not activated, cannot update damage distribution');
        }

        this.currentUnit = unit;
        this.log('Updating damage distribution for unit:', unit);

        // Hide no-selection state and show distributions
        const noSelectionDiv = this.findElement('#no-unit-selected');
        const distributionsDiv = this.findElement('#damage-distributions');
        
        if (noSelectionDiv) noSelectionDiv.classList.add('hidden');
        if (distributionsDiv) distributionsDiv.classList.remove('hidden');
        
        // Update header with current unit info
        this.updateUnitHeader(unit);
        
        // Generate damage distribution table
        this.generateUnitCombatTable(unit.unitType);
    }

    /**
     * Clear unit selection and show empty state
     */
    public clearUnitInfo(): void {
        if (!this.isActivated) {
            return;
        }

        this.currentUnit = null;
        this.log('Clearing damage distribution');

        // Show no-selection state and hide distributions
        const noSelectionDiv = this.findElement('#no-unit-selected');
        const distributionsDiv = this.findElement('#damage-distributions');
        
        if (noSelectionDiv) noSelectionDiv.classList.remove('hidden');
        if (distributionsDiv) distributionsDiv.classList.add('hidden');
    }

    /**
     * Update unit header to show which unit's damage is being displayed
     */
    private updateUnitHeader(unit: UnitData): void {
        const unitNameElement = this.findElement('#selected-unit-name');
        
        if (unitNameElement) {
            const unitDef = this.rulesTable.getUnitDefinition(unit.unitType);
            const unitName = this.theme?.getUnitName(unit.unitType) || unitDef?.name || `Unit ${unit.unitType}`;
            unitNameElement.textContent = unitName;
        }
    }

    /**
     * Generate SVG histogram for damage distribution
     */
    private createDamageHistogram(damageDistribution: any): string {
        if (!damageDistribution || !damageDistribution.ranges || damageDistribution.ranges.length === 0) {
            return '<div class="text-gray-400 text-xs">No data</div>';
        }

        const ranges = damageDistribution.ranges;
        const width = 220;
        const height = 50;
        const topPadding = 5;
        const bottomPadding = 15; // Space for x-axis labels
        const chartHeight = height - topPadding - bottomPadding;
        const barSpacing = 1;
        const barWidth = (width / ranges.length) - barSpacing;
        
        // Find max probability for scaling
        const maxProbability = Math.max(...ranges.map((r: any) => r.probability || 0));
        
        // Create SVG with x-axis labels
        let svg = `<svg width="${width}" height="${height}" class="inline-block">`;
        
        ranges.forEach((range: any, index: number) => {
            const probability = range.probability || 0;
            const barHeight = (probability / maxProbability) * chartHeight;
            const x = index * (barWidth + barSpacing);
            const y = topPadding + chartHeight - barHeight;
            
            // Use the actual damage value (assuming minValue = maxValue for each bucket)
            const damageValue = Math.round(range.minValue);
            
            // Determine color based on damage value
            let fillColor = 'rgb(156, 163, 175)'; // gray-400
            if (damageValue >= 80) {
                fillColor = 'rgb(239, 68, 68)'; // red-500
            } else if (damageValue >= 60) {
                fillColor = 'rgb(251, 146, 60)'; // orange-400
            } else if (damageValue >= 40) {
                fillColor = 'rgb(250, 204, 21)'; // yellow-400
            } else if (damageValue >= 20) {
                fillColor = 'rgb(134, 239, 172)'; // green-300
            } else if (damageValue > 0) {
                fillColor = 'rgb(147, 197, 253)'; // blue-300
            }
            
            // Add bar with improved tooltip
            const percentageStr = (probability * 100).toFixed(1);
            const tooltipText = `${percentageStr}% chance of ${damageValue} damage`;
            
            svg += `<rect x="${x}" y="${y}" width="${barWidth}" height="${barHeight}" 
                         fill="${fillColor}" opacity="0.8" 
                         class="hover:opacity-100 transition-opacity cursor-help"
                         data-tooltip="${tooltipText}">
                      <title>${tooltipText}</title>
                    </rect>`;
            
            // Add x-axis label for this bar (show every bar or every other bar depending on space)
            if (index % Math.ceil(ranges.length / 11) === 0 || index === ranges.length - 1) {
                svg += `<text x="${x + barWidth/2}" y="${height - 2}" 
                             text-anchor="middle" 
                             fill="currentColor" 
                             opacity="0.6" 
                             font-size="9">${damageValue}</text>`;
            }
        });
        
        // Add baseline
        svg += `<line x1="0" y1="${topPadding + chartHeight}" x2="${width}" y2="${topPadding + chartHeight}" 
                     stroke="currentColor" stroke-opacity="0.3" stroke-width="1"/>`;
        
        // Add expected damage marker if it exists
        if (damageDistribution.expectedDamage !== undefined) {
            const expectedX = (damageDistribution.expectedDamage / 100) * width;
            svg += `<line x1="${expectedX}" y1="${topPadding}" x2="${expectedX}" y2="${topPadding + chartHeight}" 
                         stroke="rgb(239, 68, 68)" stroke-width="2" stroke-dasharray="2,2" opacity="0.6">
                      <title>Expected: ${damageDistribution.expectedDamage.toFixed(1)}</title>
                    </line>`;
        }
        
        svg += '</svg>';
        
        return svg;
    }

    /**
     * Generate unit-unit combat damage distribution table
     */
    private generateUnitCombatTable(unitId: number): void {
        const container = this.findElement('#unit-combat-properties');
        if (!container) return;

        // Get the table template
        const tableTemplate = document.getElementById('unit-combat-table-template') as HTMLTemplateElement;
        const rowTemplate = document.getElementById('combat-row-template') as HTMLTemplateElement;
        
        if (!tableTemplate || !rowTemplate) {
            console.warn('Unit combat table templates not found');
            return;
        }
        
        // Clear existing content
        container.innerHTML = '';
        
        // Clone the table template
        const tableElement = tableTemplate.content.cloneNode(true) as DocumentFragment;
        const tbody = tableElement.querySelector('tbody');
        
        if (!tbody) {
            console.warn('Table body not found in template');
            return;
        }
        
        // Get all available units for combat comparison
        const commonUnitIds = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15];
        let hasAnyCombat = false;
        
        commonUnitIds.forEach(targetUnitId => {
            const targetUnitDef = this.rulesTable.getUnitDefinition(targetUnitId);
            if (targetUnitDef && targetUnitDef.name) {
                // Get unit-unit combat properties
                const combatProps = this.rulesTable.getUnitUnitProperties(unitId, targetUnitId);
                
                if (combatProps && combatProps.damage) {
                    // Clone the row template
                    const rowElement = rowTemplate.content.cloneNode(true) as DocumentFragment;
                    const row = rowElement.querySelector('tr');
                    
                    if (row) {
                        // Fill in the row data
                        const targetNameCell = row.querySelector('[data-target-name]');
                        const damageHistogramCell = row.querySelector('[data-damage-histogram]');
                        
                        // Use theme-specific unit name if available
                        const targetUnitName = this.theme?.getUnitName(targetUnitId) || targetUnitDef.name;
                        if (targetNameCell) targetNameCell.textContent = targetUnitName;
                        
                        // Create histogram visualization
                        if (damageHistogramCell) {
                            const damageDistribution = combatProps.damage;
                            const histogram = this.createDamageHistogram(damageDistribution);
                            const summaryText = `<div class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                Range: ${Math.round(damageDistribution.minDamage)}-${Math.round(damageDistribution.maxDamage)}, 
                                Avg: ${damageDistribution.expectedDamage.toFixed(0)}
                            </div>`;
                            damageHistogramCell.innerHTML = `<div>${histogram}${summaryText}</div>`;
                        }
                        
                        // Add alternating row colors
                        if (tbody.children.length % 2 === 1) {
                            row.classList.add('bg-gray-50', 'dark:bg-gray-700');
                        }
                        
                        tbody.appendChild(rowElement);
                        hasAnyCombat = true;
                    }
                }
            }
        });
        
        // Only append the table if we have combat data to show
        if (hasAnyCombat) {
            container.appendChild(tableElement);
        } else {
            // Show a message if no combat data available
            container.innerHTML = '<div class="text-sm text-gray-500 dark:text-gray-400 italic p-4">No combat data available for this unit</div>';
        }
    }

    /**
     * Get current unit info (for external access)
     */
    public getCurrentUnit(): UnitData | null {
        return this.currentUnit;
    }

    /**
     * Check if unit is currently selected
     */
    public hasUnitSelected(): boolean {
        return this.currentUnit !== null;
    }

    protected destroyComponent(): void {
        this.deactivate();
    }
}