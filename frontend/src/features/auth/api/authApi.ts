import { ApiError } from "../../../shared/api/error";
import { apiClient } from "../../../shared/api/client";
import type { components } from "../../../shared/api/openapi";

type LoginResponse = components["schemas"]["User"];

export async function login(email: string, password: string): Promise<LoginResponse> {
  const { data, error, response } = await apiClient.POST("/dashboard/v1/auth/login", {
    body: { email, password }
  });

  if (data === undefined) {
    throw new ApiError(response.status, error);
  }

  return data;
}
