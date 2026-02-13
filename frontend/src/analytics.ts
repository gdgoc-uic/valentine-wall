import { Analytics, logEvent as firebaseLogEvent, setCurrentScreen as firebaseSetCurrentScreen } from 'firebase/analytics';
import { app } from './firebase';
import { getAnalytics, setAnalyticsCollectionEnabled } from 'firebase/analytics';

let analytics: Analytics | undefined;

// Initialize analytics only in production, client-side, and when Firebase is configured
if (!import.meta.env.SSR && import.meta.env.PROD && app) {
  analytics = getAnalytics(app);
}

export function logEvent(eventName: string, eventParams?: Record<string, any>) {
  if (analytics) {
    firebaseLogEvent(analytics, eventName, eventParams);
  }
}

export function setCurrentScreen(screenName: string) {
  if (analytics) {
    firebaseSetCurrentScreen(analytics, screenName);
  }
}

// Disable analytics in dev mode
if (!import.meta.env.PROD && analytics) {
  setAnalyticsCollectionEnabled(analytics, false);
}
