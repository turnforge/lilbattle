import { BaseComponent, EventBus, LCMComponent } from '@panyam/tsappkit';
import { ITheme } from '../../assets/themes/BaseTheme';

/**
 * GameEndedModal announces game-over to the user. The Go presenter signals
 * the transition via UpdateGameStatusRequest.finished=true; GameViewerPageBase
 * forwards that to show(winningPlayer). The modal stays decoupled from the
 * presenter wire format — it only knows "show winner N" / "hide."
 *
 * Closing the modal is non-destructive: input remains blocked because the
 * end-state is held on GameViewerPageBase, not on the modal. Users can dismiss
 * to look at the final board without re-enabling moves.
 */
export class GameEndedModal extends BaseComponent implements LCMComponent {
    private theme: ITheme | null = null;
    private modalOverlay: HTMLElement | null = null;
    private modalContent: HTMLElement | null = null;
    private modalBody: HTMLElement | null = null;

    constructor(rootElement: HTMLElement, eventBus: EventBus, debugMode: boolean = false) {
        super('game-ended-modal', rootElement, eventBus, debugMode);
    }

    public setTheme(theme: ITheme): void {
        this.theme = theme;
    }

    async performLocalInit(): Promise<LCMComponent[]> {
        this.modalOverlay = this.rootElement;
        this.modalContent = this.rootElement.querySelector('.modal-content');
        this.modalBody = this.rootElement.querySelector('.modal-body');

        if (!this.modalOverlay || !this.modalContent || !this.modalBody) {
            throw new Error('GameEndedModal: required modal elements not found');
        }

        this.setupCloseHandlers();
        return [];
    }

    private setupCloseHandlers(): void {
        // Overlay click dismisses but does not re-enable input.
        this.modalOverlay?.addEventListener('click', (e) => {
            if (e.target === this.modalOverlay) {
                this.hide();
            }
        });

        const closeBtn = this.rootElement.querySelector('.close-button');
        closeBtn?.addEventListener('click', () => this.hide());

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && this.isVisible()) {
                this.hide();
            }
        });
    }

    /**
     * Show the modal announcing the winning player. Re-entrant: calling show
     * twice with the same winner is a no-op aside from re-rendering the body.
     */
    public show(winningPlayer: number): void {
        if (!this.modalBody || !this.modalOverlay) return;
        this.renderBody(this.modalBody, winningPlayer);
        this.modalOverlay.classList.remove('hidden');
        this.modalOverlay.classList.add('flex');
    }

    public hide(): void {
        if (!this.modalOverlay) return;
        this.modalOverlay.classList.add('hidden');
        this.modalOverlay.classList.remove('flex');
    }

    public isVisible(): boolean {
        return this.modalOverlay?.classList.contains('flex') ?? false;
    }

    private renderBody(target: HTMLElement, winningPlayer: number): void {
        // Player 0 is the unspecified / draw / tie case — the backend only
        // sets winning_player when checkVictoryConditions returns a real
        // winner, so 0 here means "ended without a winner."
        const headline = winningPlayer > 0
            ? `Player ${winningPlayer} wins!`
            : 'Game over';
        const sub = winningPlayer > 0
            ? 'The match has been decided.'
            : 'The match has ended without a winner.';

        target.textContent = '';

        const wrap = document.createElement('div');
        wrap.className = 'flex flex-col items-center text-center gap-4 py-6';

        const trophy = document.createElement('div');
        trophy.className = 'text-5xl';
        trophy.textContent = '🏆';
        wrap.appendChild(trophy);

        const h3 = document.createElement('h3');
        h3.className = 'text-2xl font-bold text-gray-900 dark:text-white';
        h3.textContent = headline;
        wrap.appendChild(h3);

        const p = document.createElement('p');
        p.className = 'text-sm text-gray-600 dark:text-gray-300';
        p.textContent = sub;
        wrap.appendChild(p);

        const actions = document.createElement('div');
        actions.className = 'flex gap-3 mt-4';

        const back = document.createElement('a');
        back.href = '/games';
        back.className = 'px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors';
        back.textContent = 'Return to games';
        actions.appendChild(back);

        const close = document.createElement('button');
        close.className = 'close-button px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-200 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors';
        close.textContent = 'Close';
        actions.appendChild(close);

        wrap.appendChild(actions);
        target.appendChild(wrap);
    }
}
