import { createRoot } from "react-dom/client";
import App from "./App";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import CheckoutPage from "./components/CheckoutPage/CheckoutPage";
import FailPaymentPage from "./components/FailPaymentPage/FailPaymentPage";
import SuccessPaymentPage from "./components/SuccessPaymentPage/SuccessPaymentPage";
import { HelmetProvider } from "react-helmet-async";
import * as Sentry from "@sentry/react";
import "./style.css";

Sentry.init({
  dsn:
    window.location.hostname === "localhost"
      ? undefined
      : import.meta.env.VITE_DSN,
  integrations: [
    Sentry.browserTracingIntegration(),
    Sentry.replayIntegration({
      maskAllText: false,
      blockAllMedia: false,
    }),
  ],
});


declare global {
  interface Window {
    kakao: any;
  }
}

const container = document.getElementById("root");

const root = createRoot(container as HTMLElement);

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/reset-password",
    element: <App />,
  },
  {
    path: "/marker",
    element: <App />,
  },
  {
    path: "/payment",
    element: <CheckoutPage />,
  },
  {
    path: "/payment/success",
    element: <SuccessPaymentPage />,
  },
  {
    path: "/payment/fail",
    element: <FailPaymentPage />,
  },
]);

const queryClient = new QueryClient();

root.render(
  // <React.StrictMode>
  <QueryClientProvider client={queryClient}>
    <HelmetProvider>
      <RouterProvider router={router} />
    </HelmetProvider>
  </QueryClientProvider>
  // </React.StrictMode>
);
