import type { HandleFetch } from "@sveltejs/kit";

export const handleFetch: HandleFetch = async ({event, request, fetch}) => {
    console.log("FETCH HANDLER RUN")

    if (request.url.startsWith('http://api:8080/')) {
    }
    return fetch(request);
}