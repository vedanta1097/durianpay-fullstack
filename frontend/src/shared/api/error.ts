import type { components } from "./openapi";

type ErrorResponse = components["schemas"]["Error"];

export class ApiError extends Error {
  status: number;

  constructor(status: number, response?: ErrorResponse) {
    super(response?.message ?? "Request failed");
    this.name = "ApiError";
    this.status = status;
  }
}
