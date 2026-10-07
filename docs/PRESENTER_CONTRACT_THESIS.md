# The Presenter-Contract Architecture

The canonical, app-agnostic version of this design reference now lives in goapplib, since
the pattern and its reference implementation (`@panyam/tsappkit`, `@panyam/tsappkit-solid`)
are shared across apps:

- `goapplib/docs/PRESENTER_CONTRACT_THESIS.md`

lilbattle is the case study in that document's appendix (game viewer = Go presenter in
WASM; world editor = in-process TS presenter). Read the goapplib copy for the principles;
this file is only a pointer so the doc is discoverable from the lilbattle tree.
