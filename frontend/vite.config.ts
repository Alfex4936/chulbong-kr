import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import VitePluginHtmlEnv from "vite-plugin-html-env";
import fs from "fs";
// import mkcert from "vite-plugin-mkcert";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    // mkcert(),
    VitePluginHtmlEnv(),
    VitePluginHtmlEnv({ compiler: true }),
  ],
  server: {
    // https: true,
    https: {
      key: fs.readFileSync("local.k-pullup.com-key.pem"),
      cert: fs.readFileSync("local.k-pullup.com.pem"),
    },
    host: "local.k-pullup.com",
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
