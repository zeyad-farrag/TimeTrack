import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    files: ["packages/core/team/**/*.{js,jsx,ts,tsx}"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            "next/*",
            "react-router-dom",
            "react-dom",
            "react-dom/*",
            "@base-ui/react",
            "@base-ui/react/*",
            "tailwindcss/*",
          ],
        },
      ],
      "no-restricted-globals": ["error", "process", "window", "localStorage"],
    },
  },
  {
    files: ["packages/views/team/**/*.{js,jsx,ts,tsx}"],
    rules: {
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            "next/*",
            "react-router-dom",
            "@/packages/core/team/store",
            "@/packages/core/team/store/*",
            "../core/team/store",
            "../core/team/store/*",
            "../../core/team/store",
            "../../core/team/store/*",
          ],
        },
      ],
    },
  },
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
]);

export default eslintConfig;
