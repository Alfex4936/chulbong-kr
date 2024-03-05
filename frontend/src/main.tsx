import { createRoot } from "react-dom/client";
import App from "./App";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import "./style.css";

declare global {
  interface Window {
    kakao: any;
  }
}

const container = document.getElementById("root");

const root = createRoot(container as HTMLElement);

const queryClient = new QueryClient();

root.render(
  // <React.StrictMode>
  <QueryClientProvider client={queryClient}>
    <App />
  </QueryClientProvider>
  // </React.StrictMode>
);
