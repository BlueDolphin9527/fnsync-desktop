import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    // outDir: "dist",
    cssCodeSplit: false,
    brotliSize: false,
    assetsInlineLimit: 100 * 1024 * 8, // 100KB
    rollupOptions: {
      output: {
        manualChunks: () => "everything.js",
        assetFileNames: "bundle.css",
        entryFileNames: "bundle.js",
      },
    },
  },
});
