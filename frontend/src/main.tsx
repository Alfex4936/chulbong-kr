import { createRoot } from "react-dom/client";
import App from "./App";
import "./style.css";

declare global {
  interface Window {
    kakao: any;
  }
}

const container = document.getElementById("root");

const root = createRoot(container as HTMLElement);

root.render(
  // <React.StrictMode>
  <App />
  // </React.StrictMode>
);
