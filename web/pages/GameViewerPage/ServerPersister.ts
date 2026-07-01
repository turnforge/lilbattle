/**
 * ServerPersister receives protojson-encoded move groups from the WASM-side
 * SingletonGamesService (via the lilbattle.registerMovePersister JS bridge)
 * and POSTs them back to the server's ProcessMoves Connect endpoint.
 *
 * Same-origin fetch with `credentials: 'include'` attaches the session
 * cookie automatically. The server's auth middleware (oneauth) + the
 * per-game player check (services/authz.CanSubmitMoves) fire on every call.
 * An anonymous viewer's write gets rejected with 401; a logged-in
 * non-player gets rejected with a forbidden-status error; the player who
 * owns the current turn succeeds. Errors are thrown so the WASM Go side
 * can propagate them; the FE surfaces them via the game-log panel.
 *
 * NOTE: The WASM side already applied moves to the in-memory state before
 * SaveMoveGroup fires. A server rejection means the browser's optimistic
 * state now diverges from the canonical server state. For issue 174 we
 * accept the divergence and rely on page refresh to restore correctness
 * (per the ticket's log-only acceptance criterion). Reconciliation is a
 * separate follow-up.
 */
export class ServerPersister {
    /**
     * POST the (state, group) pair back to the server's ProcessMoves
     * endpoint. Throws on non-2xx so the WASM bridge sees a Promise
     * rejection and turns it into a Go error.
     */
    async save(gameId: string, stateJson: string, groupJson: string): Promise<void> {
        // ProcessMoves takes moves in the request; the group carries the
        // move list. Repack into the shape the server expects.
        const group = JSON.parse(groupJson);
        const body = JSON.stringify({
            gameId,
            moves: group.moves || [],
        });

        const url = '/api/lilbattle.v1.GamesService/ProcessMoves';
        const resp = await fetch(url, {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body,
        });

        if (resp.ok) {
            return;
        }

        // Non-2xx: read the body for whatever error text the Connect
        // handler produced (typically a JSON envelope with `code` and
        // `message`). Throw with a compact human-readable message.
        const text = await resp.text().catch(() => '');
        const suffix = text ? `: ${text.slice(0, 200)}` : '';
        throw new Error(`ProcessMoves ${resp.status} ${resp.statusText}${suffix}`);
    }
}
