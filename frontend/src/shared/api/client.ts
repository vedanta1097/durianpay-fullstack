import createClient from "openapi-fetch";
import { useAuthStore } from "../../features/auth/store/authStore";
import type { paths } from "./openapi";

const baseUrl = import.meta.env.VITE_API_BASE_URL || window.location.origin;

export const apiClient = createClient<paths>({ baseUrl });

apiClient.use({
  onRequest({ request }) {
    const token = useAuthStore.getState().token;
    if (token !== null) {
      request.headers.set("Authorization", `Bearer ${token}`);
    }
    return request;
  }
});
