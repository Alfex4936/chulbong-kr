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
        target: "http://128.134.184.2:9452",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api\/v1/, '')
      },
    },
  },
});
