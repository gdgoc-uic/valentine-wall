export function isReadOnly() {
    return import.meta.env.VITE_READ_ONLY === "true";
}