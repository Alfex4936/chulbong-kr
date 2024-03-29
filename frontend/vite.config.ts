import react from "@vitejs/plugin-react";
import { ConfigEnv, defineConfig, loadEnv } from "vite";
import VitePluginHtmlEnv from "vite-plugin-html-env";
import fs from "fs";

// import mkcert from "vite-plugin-mkcert";

export default ({ mode }: ConfigEnv) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  const isDevelop = process.env.VITE_DEVELOP === "true";

  return defineConfig({
    plugins: [
      react(),
      // mkcert(),
      VitePluginHtmlEnv(),
      VitePluginHtmlEnv({ compiler: true }),
    ],
    server: {
      // https: true,
      https: isDevelop
        ? {
            key: fs.readFileSync("local.k-pullup.com-key.pem"),
            cert: fs.readFileSync("local.k-pullup.com.pem"),
          }
        : false,
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
};
