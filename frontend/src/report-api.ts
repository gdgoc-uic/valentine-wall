import { expandEndpoint } from "./client";

export const expandReportApiEndpoint = (endpoint: string): string => {
  return expandEndpoint(import.meta.env.VITE_REPORT_API_URL, endpoint);
}

export const headers = {
  'Content-Type': 'application/json',
  'xc-token': import.meta.env.VITE_REPORT_API_KEY
};