import axios, { AxiosInstance } from "axios";

const baseInterceptors = (instance: AxiosInstance) => {
  let isRateLimitAlertShown = false;
  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (
        error.response &&
        error.response.status === 429 &&
        !isRateLimitAlertShown
      ) {
        isRateLimitAlertShown = true;
        alert("요청이 너무 많습니다. 나중에 다시 시도해주세요.");
        setTimeout(() => (isRateLimitAlertShown = false), 3000);
      }
      return Promise.reject(error);
    }
  );
};

const base = () => {
  const instance = axios.create({
    // timeout: 1000,
  });

  baseInterceptors(instance);

  return instance;
};

const instance = base();

export default instance;
