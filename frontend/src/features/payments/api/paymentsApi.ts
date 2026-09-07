import { ApiError } from "../../../shared/api/error";
import { apiClient } from "../../../shared/api/client";
import type { components } from "../../../shared/api/openapi";

type Payment = components["schemas"]["Payment"];
type PaymentStatus = components["schemas"]["PaymentStatus"];

export async function listPayments(status?: PaymentStatus): Promise<Payment[]> {
  const { data, error, response } = await apiClient.GET("/dashboard/v1/payments", {
    params: { query: status === undefined ? {} : { status } }
  });

  if (data === undefined) {
    throw new ApiError(response.status, error);
  }

  return data.payments ?? [];
}
