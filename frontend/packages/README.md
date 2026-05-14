# Package Boundaries

`packages/core/team/**` is framework-neutral domain code: it must not import `next/*`, `react-router-dom`, `react-dom`, UI libraries, Tailwind internals, browser globals, or `process.env`. `packages/views/team/**` may render React views but must not import `next/*` or `react-router-dom`; framework wiring belongs under `app/**`. These rules are enforced by `eslint.config.mjs`.
