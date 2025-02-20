import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
export default defineConfig({
  plugins: [react()],
  css: {
    postcss: {
      plugins: [
        async () => {
          const tailwindcss = await import("tailwindcss");
          return tailwindcss.default;
        },
        async () => {
          const autoprefixer = await import("autoprefixer");
          return autoprefixer.default;
        },
      ],
    },
  },
});
