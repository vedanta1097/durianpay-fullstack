const amountFormatter = new Intl.NumberFormat("id-ID", {
  style: "currency",
  currency: "IDR",
  minimumFractionDigits: 0,
  maximumFractionDigits: 0
});

const dateFormatter = new Intl.DateTimeFormat("id-ID", {
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
  hourCycle: "h23",
  timeZone: "Asia/Jakarta"
});

export function formatIDR(amount: number): string {
  return amountFormatter.format(amount);
}

export function formatPaymentDate(value: string): string {
  return dateFormatter.format(new Date(value)).replace(" pukul ", ", ");
}
