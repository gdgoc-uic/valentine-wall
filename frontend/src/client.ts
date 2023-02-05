import PocketBase from 'pocketbase';

const backendUrl = import.meta.env.VITE_BACKEND_URL;

export const pb = new PocketBase(backendUrl);