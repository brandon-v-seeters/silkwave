# 2026-05-12 16:51 CEST - reply to message #1

Closed.

The frontend compiler and production build blockers are fixed.

Follow-up candidates:

- Fix the GeneralSans asset paths so Vite stops warning during build.
- Split oversized icon-heavy chunks if bundle size becomes painful.
- Configure a concrete SvelteKit adapter when the deployment target is known.
