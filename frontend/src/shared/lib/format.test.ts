import { describe, expect, it } from "vitest";
import { formatIDR, formatPaymentDate } from "./format";

describe("payment formatters", () => {
  it("formats integer amounts as IDR without decimals", () => {
    expect(formatIDR(125000)).toContain("125.000");
    expect(formatIDR(125000)).not.toContain(",00");
  });

  it("formats timestamps in the browser timezone with a colon-separated time", () => {
    const value = "2026-03-10T08:15:00Z";
    const expected = new Intl.DateTimeFormat("en-GB", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hourCycle: "h23"
    }).format(new Date(value));

    expect(formatPaymentDate(value)).toBe(expected);
  });
});
