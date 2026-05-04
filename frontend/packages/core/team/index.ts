// Boundary: core/team must not import Next.js, React DOM, router packages, UI libraries, or process.env.
// TanStack Query owns server state; Zustand owns client state. Views consume core through this public surface.
export const teamCorePackage = "team-core";
