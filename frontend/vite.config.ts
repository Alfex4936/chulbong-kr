import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import VitePluginHtmlEnv from "vite-plugin-html-env";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    VitePluginHtmlEnv(),
    VitePluginHtmlEnv({ compiler: true }),
  ],
  server: {
    proxy: {
      "/api/v1": {
        target: "https://api.k-pullup.com",
        changeOrigin: true,
        secure: false,
      },
      "/ws": {
        target: "wss://api.k-pullup.com",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
