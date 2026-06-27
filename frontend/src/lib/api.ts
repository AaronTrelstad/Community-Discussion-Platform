import axios, { InternalAxiosRequestConfig } from "axios";

interface RetryableRequestConfig extends InternalAxiosRequestConfig {
    _retry?: boolean;
}

const api = axios.create({
    baseURL: "/api", 
    headers: { "Content-Type": "application/json" },
    withCredentials: true,
});

let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value: unknown) => void;
    reject: (reason: unknown) => void;
}> = [];

function processQueue(error: unknown) {
    failedQueue.forEach((p) => (error ? p.reject(error) : p.resolve(null)));
    failedQueue = [];
}

api.interceptors.response.use(
    (response) => response,
    async (error: unknown) => {
        if (!axios.isAxiosError(error)) return Promise.reject(error);

        const original = error.config as RetryableRequestConfig;

        if (
            error.response?.status === 401 &&
            original &&
            !original._retry &&
            !original.url?.includes("/auth/refresh") &&
            !original.url?.includes("/auth/login")
        ) {
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                }).then(() => api(original));
            }

            original._retry = true;
            isRefreshing = true;

            try {
                await api.post("/auth/refresh");
                processQueue(null);
                return api(original);
            } catch (err) {
                processQueue(err);
                if (!window.location.pathname.includes("/login")) {
                    window.location.href = "/login";
                }
                return Promise.reject(err);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

export default api;
