import { BasePage } from '../lib/BasePage';
import { EventBus } from '../lib/EventBus';
import { LCMComponent } from '../lib/LCMComponent';
import WeewarBundle from '../gen/wasmjs';
import { GamesServiceServiceClient } from '../gen/wasmjs/weewar/v1/gamesServiceClient';
import { SimulateAttackRequest, SimulateAttackResponse } from '../gen/wasmjs/weewar/v1/interfaces';

/**
 * Attack Simulator Page - Interactive combat simulator
 * Allows users to simulate combat between different units on different terrains
 */
class AttackSimulatorPage extends BasePage {
    private wasmBundle: WeewarBundle | null = null;
    private gamesServiceClient: GamesServiceServiceClient | null = null;

    // Canvas elements
    private attackerHexCanvas: HTMLCanvasElement;
    private defenderHexCanvas: HTMLCanvasElement;
    private attackerChartCanvas: HTMLCanvasElement;
    private defenderChartCanvas: HTMLCanvasElement;

    // Form elements
    private attackerUnitSelect: HTMLSelectElement;
    private attackerTerrainSelect: HTMLSelectElement;
    private attackerHealthInput: HTMLInputElement;
    private defenderUnitSelect: HTMLSelectElement;
    private defenderTerrainSelect: HTMLSelectElement;
    private defenderHealthInput: HTMLInputElement;
    private woundBonusInput: HTMLInputElement;
    private numSimulationsInput: HTMLInputElement;
    private simulateButton: HTMLButtonElement;

    // Stat elements
    private attackerMeanDamageEl: HTMLElement;
    private attackerKillProbEl: HTMLElement;
    private defenderMeanDamageEl: HTMLElement;
    private defenderKillProbEl: HTMLElement;

    constructor() {
        super('attack-simulator-page', new EventBus(), false);
    }

    // Override lifecycle methods from BasePage
    protected override initializeSpecificComponents(): LCMComponent[] {
        // Get canvas elements
        this.attackerHexCanvas = document.getElementById('attacker-hex') as HTMLCanvasElement;
        this.defenderHexCanvas = document.getElementById('defender-hex') as HTMLCanvasElement;
        this.attackerChartCanvas = document.getElementById('attacker-chart') as HTMLCanvasElement;
        this.defenderChartCanvas = document.getElementById('defender-chart') as HTMLCanvasElement;

        // Get form elements
        this.attackerUnitSelect = document.getElementById('attacker-unit') as HTMLSelectElement;
        this.attackerTerrainSelect = document.getElementById('attacker-terrain') as HTMLSelectElement;
        this.attackerHealthInput = document.getElementById('attacker-health') as HTMLInputElement;
        this.defenderUnitSelect = document.getElementById('defender-unit') as HTMLSelectElement;
        this.defenderTerrainSelect = document.getElementById('defender-terrain') as HTMLSelectElement;
        this.defenderHealthInput = document.getElementById('defender-health') as HTMLInputElement;
        this.woundBonusInput = document.getElementById('wound-bonus') as HTMLInputElement;
        this.numSimulationsInput = document.getElementById('num-simulations') as HTMLInputElement;
        this.simulateButton = document.getElementById('simulate-btn') as HTMLButtonElement;

        // Get stat elements
        this.attackerMeanDamageEl = document.getElementById('attacker-mean-damage')!;
        this.attackerKillProbEl = document.getElementById('attacker-kill-prob')!;
        this.defenderMeanDamageEl = document.getElementById('defender-mean-damage')!;
        this.defenderKillProbEl = document.getElementById('defender-kill-prob')!;

        // Initialize async components
        this.initAsync();

        // No child components
        return [];
    }

    protected override bindSpecificEvents(): void {
        // Simulate button
        this.simulateButton.addEventListener('click', () => this.runSimulation());

        // Auto-simulate on form changes
        const autoSimulate = () => this.runSimulation();
        this.attackerUnitSelect.addEventListener('change', autoSimulate);
        this.attackerTerrainSelect.addEventListener('change', autoSimulate);
        this.attackerHealthInput.addEventListener('input', autoSimulate);
        this.defenderUnitSelect.addEventListener('change', autoSimulate);
        this.defenderTerrainSelect.addEventListener('change', autoSimulate);
        this.defenderHealthInput.addEventListener('input', autoSimulate);
        this.woundBonusInput.addEventListener('input', autoSimulate);
        this.numSimulationsInput.addEventListener('input', autoSimulate);
    }

    private async initAsync(): Promise<void> {
        // Load WASM
        await this.loadWASM();

        // Run initial simulation (dropdowns are already populated server-side)
        await this.runSimulation();
    }

    private async loadWASM(): Promise<void> {
        try {
            console.log('[AttackSimulator] Loading WASM bundle...');
            this.wasmBundle = new WeewarBundle();
            this.gamesServiceClient = new GamesServiceServiceClient(this.wasmBundle);
            await this.wasmBundle.loadWasm('/static/wasm/weewar-cli.wasm');
            await this.wasmBundle.waitUntilReady();
            console.log('[AttackSimulator] WASM loaded successfully');
        } catch (error) {
            console.error('[AttackSimulator] Failed to load WASM:', error);
            alert('Failed to load game engine');
        }
    }

    private async runSimulation(): Promise<void> {
        if (!this.gamesServiceClient) {
            console.error('[AttackSimulator] WASM not loaded yet');
            return;
        }

        // Build request
        const request: SimulateAttackRequest = {
            attackerUnitType: parseInt(this.attackerUnitSelect.value),
            attackerTerrain: parseInt(this.attackerTerrainSelect.value),
            attackerHealth: parseInt(this.attackerHealthInput.value),
            defenderUnitType: parseInt(this.defenderUnitSelect.value),
            defenderTerrain: parseInt(this.defenderTerrainSelect.value),
            defenderHealth: parseInt(this.defenderHealthInput.value),
            woundBonus: parseInt(this.woundBonusInput.value),
            numSimulations: parseInt(this.numSimulationsInput.value),
        };

        console.log('[AttackSimulator] Running simulation with:', request);

        try {
            // Call WASM RPC
            const response: SimulateAttackResponse = await this.gamesServiceClient.simulateAttack(request);
            console.log('[AttackSimulator] Simulation result:', response);

            // Update visualizations
            this.renderHexes(request);
            this.renderCharts(response);
            this.updateStats(response);
        } catch (error) {
            console.error('[AttackSimulator] Simulation failed:', error);
            alert('Simulation failed: ' + error);
        }
    }

    private renderHexes(request: SimulateAttackRequest): void {
        // Render attacker hex
        this.drawHex(
            this.attackerHexCanvas,
            this.getTerrainName(request.attackerTerrain),
            this.getUnitName(request.attackerUnitType),
            request.attackerHealth,
            '#3b82f6' // Blue for attacker
        );

        // Render defender hex
        this.drawHex(
            this.defenderHexCanvas,
            this.getTerrainName(request.defenderTerrain),
            this.getUnitName(request.defenderUnitType),
            request.defenderHealth,
            '#ef4444' // Red for defender
        );
    }

    private drawHex(
        canvas: HTMLCanvasElement,
        terrainName: string,
        unitName: string,
        health: number,
        color: string
    ): void {
        const ctx = canvas.getContext('2d')!;
        const width = canvas.width;
        const height = canvas.height;
        const centerX = width / 2;
        const centerY = height / 2;
        const size = 60;

        // Clear canvas
        ctx.clearRect(0, 0, width, height);

        // Draw hexagon
        ctx.beginPath();
        for (let i = 0; i < 6; i++) {
            const angle = (Math.PI / 3) * i - Math.PI / 2;
            const x = centerX + size * Math.cos(angle);
            const y = centerY + size * Math.sin(angle);
            if (i === 0) {
                ctx.moveTo(x, y);
            } else {
                ctx.lineTo(x, y);
            }
        }
        ctx.closePath();

        // Fill with terrain-based color
        ctx.fillStyle = this.getTerrainColor(terrainName);
        ctx.fill();
        ctx.strokeStyle = color;
        ctx.lineWidth = 3;
        ctx.stroke();

        // Draw unit name
        ctx.fillStyle = isDarkMode() ? '#fff' : '#000';
        ctx.font = 'bold 14px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(unitName, centerX, centerY - 10);

        // Draw health
        ctx.font = '12px sans-serif';
        ctx.fillText(`HP: ${health}`, centerX, centerY + 10);
    }

    private getTerrainColor(terrainName: string): string {
        const colors: Record<string, string> = {
            'Grass': '#90EE90',
            'Mountain': '#A0522D',
            'Water': '#4682B4',
            'Desert': '#F4A460',
            'Swamp': '#556B2F',
            'Forest': '#228B22',
        };
        return colors[terrainName] || '#CCC';
    }

    private renderCharts(response: SimulateAttackResponse): void {
        // Render attacker damage distribution
        this.drawBarChart(
            this.attackerChartCanvas,
            response.attackerDamageDistribution,
            'Damage to Defender',
            '#3b82f6'
        );

        // Render defender damage distribution
        this.drawBarChart(
            this.defenderChartCanvas,
            response.defenderDamageDistribution,
            'Damage to Attacker',
            '#ef4444'
        );
    }

    private drawBarChart(
        canvas: HTMLCanvasElement,
        distribution: { [key: number]: number },
        title: string,
        color: string
    ): void {
        const ctx = canvas.getContext('2d')!;
        const width = canvas.width;
        const height = canvas.height;

        // Clear canvas
        ctx.clearRect(0, 0, width, height);

        // Convert distribution to sorted array
        const data = Object.entries(distribution)
            .map(([damage, count]) => ({ damage: parseInt(damage), count }))
            .sort((a, b) => a.damage - b.damage);

        if (data.length === 0) {
            return;
        }

        // Chart dimensions
        const padding = 40;
        const chartWidth = width - 2 * padding;
        const chartHeight = height - 2 * padding;

        // Find max count for scaling
        const maxCount = Math.max(...data.map(d => d.count));

        // Bar width
        const barWidth = chartWidth / data.length;

        // Draw bars
        data.forEach((d, i) => {
            const barHeight = (d.count / maxCount) * chartHeight;
            const x = padding + i * barWidth;
            const y = height - padding - barHeight;

            ctx.fillStyle = color;
            ctx.fillRect(x, y, barWidth * 0.8, barHeight);

            // Draw damage value below bar
            ctx.fillStyle = isDarkMode() ? '#fff' : '#000';
            ctx.font = '10px sans-serif';
            ctx.textAlign = 'center';
            ctx.fillText(d.damage.toString(), x + barWidth * 0.4, height - padding + 15);
        });

        // Draw axes
        ctx.strokeStyle = isDarkMode() ? '#666' : '#ccc';
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(padding, height - padding);
        ctx.lineTo(width - padding, height - padding);
        ctx.moveTo(padding, padding);
        ctx.lineTo(padding, height - padding);
        ctx.stroke();
    }

    private updateStats(response: SimulateAttackResponse): void {
        this.attackerMeanDamageEl.textContent = response.attackerMeanDamage.toFixed(2);
        this.attackerKillProbEl.textContent = (response.attackerKillProbability * 100).toFixed(1) + '%';
        this.defenderMeanDamageEl.textContent = response.defenderMeanDamage.toFixed(2);
        this.defenderKillProbEl.textContent = (response.defenderKillProbability * 100).toFixed(1) + '%';
    }

    private getUnitName(unitId: number): string {
        const option = this.attackerUnitSelect.querySelector(`option[value="${unitId}"]`) ||
                       this.defenderUnitSelect.querySelector(`option[value="${unitId}"]`);
        return option?.textContent || 'Unknown';
    }

    private getTerrainName(terrainId: number): string {
        const option = this.attackerTerrainSelect.querySelector(`option[value="${terrainId}"]`) ||
                       this.defenderTerrainSelect.querySelector(`option[value="${terrainId}"]`);
        return option?.textContent || 'Unknown';
    }
}

function isDarkMode(): boolean {
    return document.documentElement.classList.contains('dark');
}

// Initialize the page when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    (window as any).Page = new AttackSimulatorPage();
});
