import { describe, expect, it } from "vitest";
import { formatIDR, formatPaymentDate } from "./format";

describe("payment formatters", () => {
  it("formats integer amounts as IDR without decimals", () => {
    expect(formatIDR(125000)).toContain("125.000");
    expect(formatIDR(125000)).not.toContain(",00");
  });

  it("formats timestamps in Jakarta time", () => {
    const formatted = formatPaymentDate("2026-03-10T08:15:00Z");
    expect(formatted).toContain("10 Mar 2026");
    expect(formatted).toContain("15.15");
  });
});
