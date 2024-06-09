import type { HandleFetch } from "@sveltejs/kit";

export const handleFetch: HandleFetch = async ({event, request, fetch}) => {
    console.log("SERVER FETCH HANDLER RUN")

    if (request.url.startsWith('http://api:8080/')) {
        const cookie = event.request.headers.get('cookie') as string;
        console.log("Including cookie:", cookie)
        request.headers.set('cookie', cookie);
    }
    return fetch(request);
}